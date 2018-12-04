package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSessionRestoration tests manual session restoration by key
func TestSessionRestoration(t *testing.T) {
	sessionStorage := make(map[string]*wwr.Session)

	currentStep := 1
	var createdSession *wwr.Session

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				if currentStep == 2 {
					// Expect the session to be automatically restored
					CompareSessions(t, createdSession, conn.Session())
					return wwr.Payload{}, nil
				}

				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				return wwr.Payload{}, err
			},
		},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				// Saves the session
				SessionCreated: func(conn wwr.Connection) error {
					sess := conn.Session()
					sessionStorage[sess.Key] = sess
					return nil
				},
				// Finds session by key
				SessionLookup: func(key string) (
					wwr.SessionLookupResult,
					error,
				) {
					// Expect the key of the created session to be looked up
					assert.Equal(t, createdSession.Key, key)
					assert.Contains(t, sessionStorage, key)
					session := sessionStorage[key]
					return wwr.NewSessionLookupResult(
						session.Creation,                      // Creation
						session.LastLookup,                    // LastLookup
						wwr.SessionInfoToVarMap(session.Info), // Info
					), nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	initialClient := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	/*****************************************************************\
		Step 1 - Create session and disconnect
	\*****************************************************************/

	// Create a new session
	_, err := initialClient.Connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{Data: []byte("auth")},
	)
	require.NoError(t, err)

	createdSession = initialClient.Connection.Session()

	// Disconnect client without closing the session
	initialClient.Connection.Close()

	/*****************************************************************\
		Step 2 - Create new client, restore session from key
	\*****************************************************************/
	currentStep = 2

	// Initialize client
	secondClient := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	// Ensure there's no active session on the second client
	require.Nil(t, secondClient.Connection.Session())

	// Try to manually restore the session
	// using the initial clients session key
	require.NoError(t, secondClient.Connection.RestoreSession(
		context.Background(),
		[]byte(createdSession.Key),
	))

	// Verify session
	sessionAfter := secondClient.Connection.Session()
	require.NotEqual(t, "", sessionAfter.Key)
	CompareSessions(t, createdSession, sessionAfter)
}
