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
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	clients := make([]*TestClient, connectedClientsNum)
	for i := 0; i < connectedClientsNum; i++ {
		client := setup.NewClient(
			wwrclt.Options{
				Autoconnect: wwr.Disabled,
			},
			nil, // Use the default transport implementation
			TestClientHooks{},
		)
		require.NoError(t, client.Connection.Connect())
		clients[i] = client
	}

	require.NoError(t, setup.Server.Shutdown())
}
