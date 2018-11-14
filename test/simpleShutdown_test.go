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
	setup := setupTestServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	clients := make([]*testClient, connectedClientsNum)
	for i := 0; i < connectedClientsNum; i++ {
		client := setup.newClient(
			wwrclt.Options{
				Autoconnect: wwr.Disabled,
			},
			testClientHooks{},
		)
		require.NoError(t, client.connection.Connect())
		defer client.connection.Close()
		clients[i] = client
	}

	require.NoError(t, setup.Server.Shutdown())
}
