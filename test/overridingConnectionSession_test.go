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
	setup := setupTestServer(
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
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		testClientHooks{},
	)

	require.NoError(t, client.connection.Connect())
}
