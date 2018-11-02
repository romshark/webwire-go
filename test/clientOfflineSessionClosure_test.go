package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientOfflineSessionClosure tests offline session closure
func TestClientOfflineSessionClosure(t *testing.T) {
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
					// Expect the session to be removed
					assert.False(t,
						conn.HasSession(),
						"Expected client to be anonymous",
					)
					return nil, nil
				}

				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				return nil, err
			},
		},
		wwr.ServerOptions{
			SessionManager: &callbackPoweredSessionManager{
				// Saves the session
				SessionCreated: func(conn wwr.Connection) error {
					session := conn.Session()
					sessionStorage[session.Key] = session
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
					// Session found
					return wwr.NewSessionLookupResult(
						session.Creation,                      // Creation
						session.LastLookup,                    // LastLookup
						wwr.SessionInfoToVarMap(session.Info), // Info
					), nil
				},
			},
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	/*****************************************************************\
		Step 1 - Create session and disconnect
	\*****************************************************************/

	// Create a new session
	_, err := client.connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.NewPayload(wwr.EncodingBinary, []byte("auth")),
	)
	require.NoError(t, err)

	createdSession = client.connection.Session()

	// Disconnect client without closing the session
	client.connection.Close()

	// Ensure the session isn't lost
	require.NotEqual(t,
		wwrclt.Connected, client.connection.Status(),
		"Client is expected to be disconnected",
	)
	require.NotEqual(t,
		"", client.connection.Session().Key,
		"Session lost after disconnection",
	)

	/*****************************************************************\
		Step 2 - Close session, reconnect and verify
	\*****************************************************************/
	currentStep = 2

	require.NoError(t,
		client.connection.CloseSession(),
		"Offline session closure failed",
	)

	// Ensure the session is removed locally
	require.Nil(t, client.connection.Session(), "Session not removed")

	// Reconnect
	require.NoError(t, client.connection.Connect())

	// Ensure the client is anonymous
	_, err = client.connection.Request(
		context.Background(),
		[]byte("verify-restored"),
		wwr.NewPayload(wwr.EncodingBinary, []byte("is_restored?")),
	)
	require.NoError(t, err)
}
