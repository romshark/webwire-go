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

// TestServerSideSessionClosure tests server-side closing of sessions
func TestServerSideSessionClosure(t *testing.T) {
	simultaneousClients := 4
	require.True(t, simultaneousClients > 1)

	var sessionKey string
	var createdSession *wwr.Session

	onSessionClosedHooksExecuted := tmdwg.NewTimedWaitGroup(
		simultaneousClients,
		10*time.Second,
	)

	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
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
	clients := make([]*testClient, simultaneousClients)
	for i := 0; i < simultaneousClients; i++ {
		client := setup.newClient(
			wwrclt.Options{
				DefaultRequestTimeout: 2 * time.Second,
				Autoconnect:           wwr.Disabled,
			},
			nil, // Use the default transport implementation
			testClientHooks{
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
	}

	// Authenticate first client to get the session key
	firstClient := clients[0]
	reply, err := firstClient.connection.Request(
		context.Background(),
		[]byte("auth"),
		wwr.Payload{},
	)
	require.NoError(t, err)
	reply.Close()

	// Extract session key
	createdSession = firstClient.connection.Session()
	sessionKey = createdSession.Key
	require.NotNil(t, createdSession)

	// Apply the session to other remaining clients
	for i := 1; i < len(clients); i++ {
		clt := clients[i]
		require.NoError(t, clt.connection.RestoreSession(
			context.Background(),
			[]byte(sessionKey),
		))
	}

	// Ensure all clients are logged into 1 session
	for _, client := range clients {
		session := client.connection.Session()
		require.Equal(t, sessionKey, session.Key)
		compareSessions(t, createdSession, session)
	}

	// Close the session
	affectedConnections, closeErrors, err := setup.Server.CloseSession(
		sessionKey,
	)
	require.NoError(t, err)
	require.Len(t, affectedConnections, simultaneousClients)
	require.Len(t, closeErrors, simultaneousClients)
	for _, err := range closeErrors {
		require.NoError(t, err)
	}

	// Expect the session creation hook to be executed in the client
	require.NoError(t,
		onSessionClosedHooksExecuted.Wait(),
		"client.OnSessionClosed hook wasn't executed",
	)

	// Ensure the session was properly closed for all affected connections
	for _, client := range clients {
		require.Nil(t, client.connection.Session())
	}
}
