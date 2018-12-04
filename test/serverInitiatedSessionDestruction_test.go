package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerInitiatedSessionDestruction verifies
// server-initiated session destruction
func TestServerInitiatedSessionDestruction(t *testing.T) {
	sessionCreationCallbackCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	sessionDestructionCallbackCalled := tmdwg.NewTimedWaitGroup(
		1,
		1*time.Second,
	)
	var createdSession *wwr.Session
	expectedCredentials := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("secret_credentials"),
	}
	placeholderMessage := wwr.Payload{
		Encoding: wwr.EncodingBinary,
		Data:     []byte("nothinginteresting"),
	}
	currentStep := 1

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// On step 2 - verify session creation and correctness
				if currentStep == 2 {
					session := conn.Session()
					CompareSessions(t, createdSession, session)
					assert.Equal(t, session.Key, string(msg.Payload()))
					return wwr.Payload{}, nil
				}

				// on step 3 - close session and verify its destruction
				if currentStep == 3 {
					/******************************************************\
						Server-side session destruction initiation
					\******************************************************/
					// Attempt to destroy this clients session
					// on the end of the first step
					assert.NoError(t, conn.CloseSession())

					// Verify destruction
					assert.Nil(t,
						conn.Session(),
						"Expected the session to be destroyed",
					)

					return wwr.Payload{}, nil
				}

				// On step 4 - verify session destruction
				if currentStep == 4 {
					assert.Nil(t,
						conn.Session(),
						"Expected the session to be destroyed, got: %v",
					)
					return wwr.Payload{}, nil
				}

				// On step 1 - authenticate and create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				if err != nil {
					return wwr.Payload{}, err
				}

				// Return the key of the newly created session
				return wwr.Payload{Data: []byte(conn.SessionKey())}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{
			OnSessionCreated: func(_ *wwr.Session) {
				// Mark the client-side session creation callback executed
				sessionCreationCallbackCalled.Progress(1)
			},
			OnSessionClosed: func() {
				// Ensure this callback is called during the
				assert.Equal(t,
					3, currentStep,
					"Client-side session destruction callback "+
						"called at wrong step",
				)
				sessionDestructionCallbackCalled.Progress(1)
			},
		},
	)

	/*****************************************************************\
		Step 1 - Session Creation
	\*****************************************************************/
	// Send authentication request
	authReqReply, err := client.Connection.Request(
		context.Background(),
		[]byte("login"),
		expectedCredentials,
	)
	require.NoError(t, err)

	createdSession = client.Connection.Session()

	// Verify reply
	require.Equal(t, wwr.EncodingBinary, authReqReply.PayloadEncoding())
	require.Equal(t, []byte(createdSession.Key), authReqReply.Payload())
	authReqReply.Close()

	// Wait for the client-side session creation callback to be executed
	require.NoError(t,
		sessionCreationCallbackCalled.Wait(),
		"Session creation callback not called",
	)

	// Ensure the session was locally created
	require.NotEqual(t,
		nil,
		client.Connection.Session(),
		"Expected session on client-side",
	)

	/*****************************************************************\
		Step 2 - Session Creation Verification
	\*****************************************************************/
	currentStep = 2

	// Send a test-request to verify the session creation on the server
	_, err = client.Connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte(client.Connection.Session().Key)},
	)
	require.NoError(t, err)

	/*****************************************************************\
		Step 3 - Server-Side Session Destruction
	\*****************************************************************/
	currentStep = 3

	// Request session destruction
	_, err = client.Connection.Request(
		context.Background(),
		nil,
		placeholderMessage,
	)
	require.NoError(t, err)

	// Wait for the client-side session destruction callback to be called
	require.NoError(t,
		sessionDestructionCallbackCalled.Wait(),
		"Session destruction callback not called",
	)

	/*****************************************************************\
		Step 4 - Destruction Verification
	\*****************************************************************/
	currentStep = 4

	// Ensure the session is destroyed locally as well
	require.Nil(t,
		client.Connection.Session(),
		"Expected session to be destroyed on the client as well",
	)

	// Send a test-request to verify the session was destroyed on the server
	_, err = client.Connection.Request(
		context.Background(),
		nil,
		placeholderMessage,
	)
	require.NoError(t, err)
}
