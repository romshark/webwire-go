package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientIsConnected tests the client.Status method
func TestClientIsConnected(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(t, &serverImpl{}, wwr.ServerOptions{})

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			Autoconnect:           wwr.Disabled,
		},
		callbackPoweredClientHooks{},
	)

	require.NotEqual(t,
		wwrclt.Connected, client.connection.Status(),
		"Expected client to be disconnected "+
			"before the connection establishment",
	)

	// Connect to the server
	require.NoError(t, client.connection.Connect())

	require.Equal(t,
		wwrclt.Connected, client.connection.Status(),
		"Expected client to be connected after the connection establishment",
	)

	// Disconnect the client
	client.connection.Close()

	require.NotEqual(t,
		wwrclt.Connected, client.connection.Status(),
		"Expected client to be disconnected after closure",
	)
}
