package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientSessionRestoration tests manual session restoration by key
func TestClientSessionRestoration(t *testing.T) {
	sessionStorage := make(map[string]*wwr.Session)

	currentStep := 1
	var createdSession *wwr.Session

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				if currentStep == 2 {
					// Expect the session to be automatically restored
					compareSessions(t, createdSession, conn.Session())
					return nil, nil
				}

				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				if err != nil {
					return nil, err
				}

				// Return the key of the newly created session
				return nil, nil
			},
		},
		wwr.ServerOptions{
			SessionManager: &callbackPoweredSessionManager{
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
					return wwr.SessionLookupResult{
						Creation:   session.Creation,
						LastLookup: session.LastLookup,
						Info: wwr.SessionInfoToVarMap(
							session.Info,
						),
					}, nil
				},
			},
		},
	)

	// Initialize client
	initialClient := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, initialClient.connection.Connect())

	/*****************************************************************\
		Step 1 - Create session and disconnect
	\*****************************************************************/

	// Create a new session
	_, err := initialClient.connection.Request(
		context.Background(),
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("auth")),
	)
	require.NoError(t, err)

	createdSession = initialClient.connection.Session()

	// Disconnect client without closing the session
	initialClient.connection.Close()

	/*****************************************************************\
		Step 2 - Create new client, restore session from key
	\*****************************************************************/
	currentStep = 2

	// Initialize client
	secondClient := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, secondClient.connection.Connect())

	// Ensure there's no active session on the second client
	require.Nil(t, secondClient.connection.Session())

	// Try to manually restore the session
	// using the initial clients session key
	require.NoError(t, secondClient.connection.RestoreSession(
		[]byte(createdSession.Key),
	))

	// Verify session
	sessionAfter := secondClient.connection.Session()
	require.NotEqual(t, "", sessionAfter.Key)
	compareSessions(t, createdSession, sessionAfter)
}
