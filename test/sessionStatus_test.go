package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/require"
)

// TestSessionStatus tests session monitoring methods
func TestSessionStatus(t *testing.T) {
	sessionKey := "testsessionkey"

	sessionCreation := time.Now()

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				SessionLookup: func(key string) (
					wwr.SessionLookupResult,
					error,
				) {
					if key != string(sessionKey) {
						// Session not found
						return nil, nil
					}
					return wwr.NewSessionLookupResult(
						sessionCreation, // Creation
						time.Now(),      // LastLookup
						nil,             // Info
					), nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	require.Equal(t, 0, setup.Server.ActiveSessionsNum())

	// Initialize client A
	clientA, _ := setup.NewClientSocket()

	requestRestoreSessionSuccess(t, clientA, []byte(sessionKey))

	// Check status, expect 1 session with 1 connection
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Equal(t, 1, len(setup.Server.SessionConnections(sessionKey)))

	// Initialize client B
	clientB, _ := setup.NewClientSocket()

	requestRestoreSessionSuccess(t, clientB, []byte(sessionKey))

	// Check status, expect 1 session with 2 connections
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 2, setup.Server.SessionConnectionsNum(sessionKey))
	require.Equal(t, 2, len(setup.Server.SessionConnections(sessionKey)))

	// Close first connection
	require.NoError(t, clientA.Close())

	// Wait for the server to close client A
	time.Sleep(50 * time.Millisecond)

	// Check status, expect 1 session with 1 connection
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Equal(t, 1, len(setup.Server.SessionConnections(sessionKey)))

	// Close session
	requestCloseSessionSuccess(t, clientB)

	// Wait for the server to close client B
	time.Sleep(50 * time.Millisecond)

	// Check status, expect 0 sessions
	require.Equal(t, 0, setup.Server.ActiveSessionsNum())
	require.Equal(t, -1, setup.Server.SessionConnectionsNum(sessionKey))
	require.Nil(t, setup.Server.SessionConnections(sessionKey))
}
