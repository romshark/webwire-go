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
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(conn wwr.Connection) {
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
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())
}
