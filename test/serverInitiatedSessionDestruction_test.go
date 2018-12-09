package test

import (
	"context"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerInitiatedSessionDestruction tests client-initiated session
// destruction
func TestServerInitiatedSessionDestruction(t *testing.T) {
	sessionDestructionCallbackCalled := sync.WaitGroup{}
	sessionDestructionCallbackCalled.Add(1)
	signalReceived := sync.WaitGroup{}
	signalReceived.Add(1)

	sessionKey := "testsessionkey"
	sessionCreation := time.Now()

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				c wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Verify session destruction
				assert.Nil(t, c.Session())
				return wwr.Payload{}, nil
			},
			Signal: func(
				_ context.Context,
				c wwr.Connection,
				_ wwr.Message,
			) {
				c.CloseSession()
				signalReceived.Done()
			},
		},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				SessionLookup: func(
					key string,
				) (wwr.SessionLookupResult, error) {
					if key != sessionKey {
						return nil, nil
					}
					return wwr.NewSessionLookupResult(
						sessionCreation, // Creation
						time.Now(),      // LastLookup
						nil,             // Info
					), nil
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
	require.Equal(t, -1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Nil(t, setup.Server.SessionConnections(sessionKey))

	sock, _ := setup.NewClientSocket()

	requestRestoreSessionSuccess(t, sock, []byte(sessionKey))

	assert.NotEqual(t, "", sessionKey)
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Equal(t, 1, len(setup.Server.SessionConnections(sessionKey)))

	// Initiate session destruction
	signal(t, sock, []byte("close_session"), payload.Payload{})

	readSessionClosed(t, sock)

	signalReceived.Wait()

	// Wait for the server to finally destroy the session
	sessionDestructionCallbackCalled.Wait()

	require.Equal(t, 0, setup.Server.ActiveSessionsNum())
	require.Equal(t, -1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Equal(t, 0, len(setup.Server.SessionConnections(sessionKey)))

	// Verify session destruction
	requestSuccess(t, sock, 32, []byte("verify"), payload.Payload{})
}
