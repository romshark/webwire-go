package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSessionCreationOnClosedConn tests the creation of a session
// on a disconnected connection
func TestSessionCreationOnClosedConn(t *testing.T) {
	// Initialize server
	setup := setupTestServer(
		t,
		&serverImpl{
			onClientConnected: func(
				_ wwr.ConnectionOptions,
				conn wwr.Connection,
			) {
				conn.Close()
				err := conn.CreateSession(nil)
				assert.Error(t, err)
				assert.IsType(t, wwr.DisconnectedErr{}, err)
			},
			onClientDisconnected: func(conn wwr.Connection, _ error) {
				err := conn.CreateSession(nil)
				assert.Error(t, err)
				assert.IsType(t, wwr.DisconnectedErr{}, err)
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, client.connection.Connect())
}
