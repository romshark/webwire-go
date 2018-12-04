package test

import (
	"context"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerSideSessionClosureInexistent tests server-side closing of
// inexistent sessions
func TestServerSideSessionClosureInexistent(t *testing.T) {
	simultaneousClients := 4
	require.True(t, simultaneousClients > 1)

	var sessionKey string
	var createdSession *wwr.Session

	onSessionClosedHooksExecuted := sync.WaitGroup{}
	onSessionClosedHooksExecuted.Add(simultaneousClients)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				return wwr.Payload{}, err
			},
		},
		wwr.ServerOptions{
			MaxSessionConnections: uint(simultaneousClients),
		},
		nil, // Use the default transport implementation
	)

	// Initialize clients
	clients := make([]*TestClient, simultaneousClients)
	for i := 0; i < simultaneousClients; i++ {
		client := setup.NewClient(
			wwrclt.Options{
				DefaultRequestTimeout: 2 * time.Second,
				Autoconnect:           wwr.Disabled,
			},
			nil, // Use the default transport implementation
			TestClientHooks{
				OnSessionClosed: func() {
					onSessionClosedHooksExecuted.Done()
				},
			},
		)
		defer client.Connection.Close()
		clients[i] = client
	}

	// Connect clients
	for _, client := range clients {
		require.NoError(t, client.Connection.Connect())
	}

	// Authenticate first client to get the session key
	firstClient := clients[0]
	reply, err := firstClient.Connection.Request(
		context.Background(),
		[]byte("auth"),
		wwr.Payload{},
	)
	require.NoError(t, err)
	reply.Close()

	// Extract session key
	createdSession = firstClient.Connection.Session()
	sessionKey = createdSession.Key
	require.NotNil(t, createdSession)

	// Apply the session to other remaining clients
	for i := 1; i < len(clients); i++ {
		clt := clients[i]
		require.NoError(t, clt.Connection.RestoreSession(
			context.Background(),
			[]byte(sessionKey),
		))
	}

	// Ensure all clients are logged into 1 session
	for _, client := range clients {
		session := client.Connection.Session()
		require.Equal(t, sessionKey, session.Key)
		CompareSessions(t, createdSession, session)
	}

	// Compose an inexistent session key
	inexistentSessionKey := make([]byte, len(sessionKey))
	for i, c := range sessionKey {
		inexistentSessionKey[i] = byte(c)
	}
	inexistentSessionKey[0] = '0'

	// Try to close an inexistent session
	affectedConnections, closeErrors, err := setup.Server.CloseSession(
		string(inexistentSessionKey),
	)
	require.NoError(t, err)
	require.Len(t, affectedConnections, 0)
	require.Len(t, closeErrors, 0)

	// Ensure the session is still intact on all connections
	for _, client := range clients {
		require.NotNil(t, client.Connection.Session())
	}
}
