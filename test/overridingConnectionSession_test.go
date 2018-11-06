package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOverridingConnectionSession tests overriding of a connection session
func TestOverridingConnectionSession(t *testing.T) {
	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(conn wwr.Connection) {
				assert.NoError(t, conn.CreateSession(nil))
				sessionKey := conn.SessionKey()

				// Try to override the previous session
				assert.Error(t, conn.CreateSession(nil))

				// Ensure the session didn't change
				assert.Equal(t, sessionKey, conn.SessionKey())
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())
}
