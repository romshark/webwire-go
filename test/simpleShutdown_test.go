package test

import (
	"testing"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestSimpleShutdown tests simple shutdown without any pending tasks
func TestSimpleShutdown(t *testing.T) {
	connectedClientsNum := 5

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	clients := make([]*callbackPoweredClient, connectedClientsNum)
	for i := 0; i < connectedClientsNum; i++ {
		client := newCallbackPoweredClient(
			server.AddressURL(),
			wwrclt.Options{
				Autoconnect: wwr.Disabled,
			},
			callbackPoweredClientHooks{},
		)
		require.NoError(t, client.connection.Connect())
		defer client.connection.Close()
		clients[i] = client
	}

	require.NoError(t, server.Shutdown())
}
