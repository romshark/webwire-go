package test

import (
	"context"
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientInitiatedSessionDestruction tests client-initiated session
// destruction
func TestClientInitiatedSessionDestruction(t *testing.T) {
	sessionCreated := sync.WaitGroup{}
	sessionCreated.Add(1)
	sessionDestructionCallbackCalled := sync.WaitGroup{}
	sessionDestructionCallbackCalled.Add(1)

	var sessionKey string

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				// Create a new session
				assert.NoError(t, c.CreateSession(nil))
			},
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Verify session destruction
				assert.Nil(t,
					conn.Session(),
					"Expected the session to be destroyed",
				)
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				SessionCreated: func(conn wwr.Connection) error {
					defer sessionCreated.Done()

					sessionKey = conn.SessionKey()

					return nil
				},
				SessionClosed: func(closedSessionKey string) error {
					defer sessionDestructionCallbackCalled.Done()

					// Ensure that the correct session was closed
					assert.Equal(t, sessionKey, closedSessionKey)

					return nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	require.Equal(t, 0, setup.Server.ActiveSessionsNum())

	sock, _ := setup.NewClientSocket()

	// Expect session creation notification message
	readSessionCreated(t, sock)

	sessionCreated.Wait()

	assert.NotEqual(t, "", sessionKey)
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Equal(t, 1, len(setup.Server.SessionConnections(sessionKey)))

	// Initiate session destruction
	requestCloseSessionSuccess(t, sock)

	// Wait for the server to finally destroy the session
	sessionDestructionCallbackCalled.Wait()

	require.Equal(t, 0, setup.Server.ActiveSessionsNum())
	require.Equal(t, -1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Equal(t, 0, len(setup.Server.SessionConnections(sessionKey)))

	// Verify session destruction
	requestSuccess(t, sock, 32, []byte("verify"), payload.Payload{})
}
