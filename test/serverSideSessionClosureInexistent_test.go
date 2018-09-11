package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestServerSideSessionClosureInexistent tests server-side closing of
// inexistent sessions
func TestServerSideSessionClosureInexistent(t *testing.T) {
	simultaneousClients := 4
	require.True(t, simultaneousClients > 1)

	var sessionKey string
	var createdSession *wwr.Session

	onSessionClosedHooksExecuted := tmdwg.NewTimedWaitGroup(
		simultaneousClients,
		10*time.Second,
	)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				return nil, err
			},
		},
		wwr.ServerOptions{
			MaxSessionConnections: uint(simultaneousClients),
		},
	)

	// Initialize clients
	clients := make([]*callbackPoweredClient, simultaneousClients)
	for i := 0; i < simultaneousClients; i++ {
		client := newCallbackPoweredClient(
			server.Addr().String(),
			wwrclt.Options{
				DefaultRequestTimeout: 2 * time.Second,
				Autoconnect:           wwr.Disabled,
			},
			callbackPoweredClientHooks{
				OnSessionClosed: func() {
					onSessionClosedHooksExecuted.Progress(1)
				},
			},
		)
		defer client.connection.Close()
		clients[i] = client
	}

	// Connect clients
	for _, client := range clients {
		require.NoError(t, client.connection.Connect())
		time.Sleep(2000 * time.Millisecond)
	}

	// Authenticate first client to get the session key
	firstClient := clients[0]
	_, err := firstClient.connection.Request(context.Background(), "auth", nil)
	require.NoError(t, err)

	// Extract session key
	createdSession = firstClient.connection.Session()
	sessionKey = createdSession.Key
	require.NotNil(t, createdSession)

	// Apply the session to other remaining clients
	for i := 1; i < len(clients); i++ {
		clt := clients[i]
		require.NoError(t, clt.connection.RestoreSession([]byte(sessionKey)))
	}

	// Ensure all clients are logged into 1 session
	for _, client := range clients {
		session := client.connection.Session()
		require.Equal(t, sessionKey, session.Key)
		compareSessions(t, createdSession, session)
	}

	// Compose an inexistent session key
	inexistentSessionKey := make([]byte, len(sessionKey))
	for i, c := range sessionKey {
		inexistentSessionKey[i] = byte(c)
	}
	inexistentSessionKey[0] = '0'

	// Try to close an inexistent session
	affectedConnections, closeErrors, err := server.CloseSession(
		string(inexistentSessionKey),
	)
	require.NoError(t, err)
	require.Len(t, affectedConnections, 0)
	require.Len(t, closeErrors, 0)

	// Ensure the session is still intact on all connections
	for _, client := range clients {
		require.NotNil(t, client.connection.Session())
	}
}
