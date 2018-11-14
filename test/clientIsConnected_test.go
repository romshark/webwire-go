package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientIsConnected tests the client.Status method
func TestClientIsConnected(t *testing.T) {
	// Initialize webwire server given only the request
	setup := setupTestServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			Autoconnect:           wwr.Disabled,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NotEqual(t,
		wwrclt.StatusConnected, client.connection.Status(),
		"Expected client to be disconnected "+
			"before the connection establishment",
	)

	// Connect to the server
	require.NoError(t, client.connection.Connect())

	require.Equal(t,
		wwrclt.StatusConnected, client.connection.Status(),
		"Expected client to be connected after the connection establishment",
	)

	// Disconnect the client
	client.connection.Close()

	require.NotEqual(t,
		wwrclt.StatusConnected, client.connection.Status(),
		"Expected client to be disconnected after closure",
	)
}
