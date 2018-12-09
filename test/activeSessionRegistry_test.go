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

// TestActiveSessionRegistry verifies that the session registry
// of currently active sessions is properly updated
func TestActiveSessionRegistry(t *testing.T) {
	sessionCreated := sync.WaitGroup{}
	sessionCreated.Add(1)
	sessionClosed := sync.WaitGroup{}
	sessionClosed.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				// Try to create a new session
				assert.NoError(t, c.CreateSession(nil))
				sessionCreated.Done()
			},
			Signal: func(
				_ context.Context,
				c wwr.Connection,
				msg wwr.Message,
			) {
				// Close session on logout
				assert.NoError(t, c.CloseSession())

				sessionClosed.Done()
			},
		},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				SessionCreated: func(c wwr.Connection) error {
					return nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	require.Equal(t, 0, setup.Server.ActiveSessionsNum())

	// Initialize client
	sock, _ := setup.NewClientSocket()

	readSessionCreated(t, sock)

	sessionCreated.Wait()

	require.Equal(t, 1, setup.Server.ActiveSessionsNum())

	// Close session
	signal(t, sock, []byte("s"), payload.Payload{})

	readSessionClosed(t, sock)

	sessionClosed.Wait()

	require.Equal(t, 0, setup.Server.ActiveSessionsNum())
}
