package test

import (
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientDisconnectedHook verifies the server is calling the
// onClientConnected and onClientDisconnected hooks properly
func TestClientDisconnectedHook(t *testing.T) {
	connectedHookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	disconnectedHookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server given only the request
	setup := setupTestServer(
		t,
		&serverImpl{
			onClientConnected: func(conn wwr.Connection) {
				connectedHookCalled.Progress(1)
			},
			onClientDisconnected: func(conn wwr.Connection, _ error) {
				disconnectedHookCalled.Progress(1)
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

	// Connect to the server
	require.NoError(t, client.connection.Connect())

	// Await the onClientDisconnected hook to be called on the server
	require.NoError(t,
		connectedHookCalled.Wait(),
		"server.OnClientConnected hook not called",
	)

	// Disconnect the client
	client.connection.Close()

	// Await the onClientDisconnected hook to be called on the server
	require.NoError(t,
		disconnectedHookCalled.Wait(),
		"server.OnClientDisconnected hook not called",
	)
}
