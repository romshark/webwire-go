package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
)

// TestConnSessionNoOverride tests overriding of a connection session
// expecting it to fail
func TestConnSessionNoOverride(t *testing.T) {
	// TODO: fix test, wait for the server to finish OnClientConnected before
	// returning from the test function

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(
				_ wwr.ConnectionOptions,
				conn wwr.Connection,
			) {
				assert.NoError(t, conn.CreateSession(nil))
				sessionKey := conn.SessionKey()

				// Try to override the previous session
				assert.Error(t, conn.CreateSession(nil))

				// Ensure the session didn't change
				assert.Equal(t, sessionKey, conn.SessionKey())
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)
}
