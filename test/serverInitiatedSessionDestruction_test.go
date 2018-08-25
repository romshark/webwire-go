package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
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
	expectedCredentials := wwr.NewPayload(
		wwr.EncodingUtf8,
		[]byte("secret_credentials"),
	)
	placeholderMessage := wwr.NewPayload(
		wwr.EncodingBinary,
		[]byte("nothinginteresting"),
	)
	currentStep := 1

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// On step 2 - verify session creation and correctness
				if currentStep == 2 {
					session := conn.Session()
					compareSessions(t, createdSession, session)
					assert.Equal(t, session.Key, string(msg.Payload().Data()))
					return nil, nil
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

					return nil, nil
				}

				// On step 4 - verify session destruction
				if currentStep == 4 {
					assert.Nil(t,
						conn.Session(),
						"Expected the session to be destroyed, got: %v",
					)
					return nil, nil
				}

				// On step 1 - authenticate and create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				if err != nil {
					return nil, err
				}

				// Return the key of the newly created session
				return wwr.NewPayload(
					wwr.EncodingBinary,
					[]byte(conn.SessionKey()),
				), nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{
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
	require.NoError(t, client.connection.Connect())

	// Send authentication request
	authReqReply, err := client.connection.Request(
		context.Background(),
		"login",
		expectedCredentials,
	)
	require.NoError(t, err)

	createdSession = client.connection.Session()

	// Verify reply
	require.Equal(t,
		wwr.NewPayload(
			wwr.EncodingBinary,
			[]byte(createdSession.Key),
		),
		authReqReply,
	)

	// Wait for the client-side session creation callback to be executed
	require.NoError(t,
		sessionCreationCallbackCalled.Wait(),
		"Session creation callback not called",
	)

	// Ensure the session was locally created
	require.NotEqual(t,
		"",
		client.connection.Session(),
		"Expected session on client-side",
	)

	/*****************************************************************\
		Step 2 - Session Creation Verification
	\*****************************************************************/
	currentStep = 2

	// Send a test-request to verify the session creation on the server
	_, err = client.connection.Request(
		context.Background(),
		"",
		wwr.NewPayload(
			wwr.EncodingBinary,
			[]byte(client.connection.Session().Key),
		),
	)
	require.NoError(t, err)

	/*****************************************************************\
		Step 3 - Server-Side Session Destruction
	\*****************************************************************/
	currentStep = 3

	// Request session destruction
	_, err = client.connection.Request(
		context.Background(),
		"",
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
		client.connection.Session(),
		"Expected session to be destroyed on the client as well",
	)

	// Send a test-request to verify the session was destroyed on the server
	_, err = client.connection.Request(
		context.Background(),
		"",
		placeholderMessage,
	)
	require.NoError(t, err)
}
