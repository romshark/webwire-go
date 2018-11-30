package test

import (
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientConnIsConnected tests the IsActive method of a connection
func TestClientConnIsConnected(t *testing.T) {
	var clientConn wwr.Connection
	connectionReady := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	clientDisconnected := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	testerGoroutineFinished := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{
			onClientConnected: func(
				_ wwr.ConnectionOptions,
				newConn wwr.Connection,
			) {
				assert.True(t,
					newConn.IsActive(),
					"Expected connection to be active",
				)
				clientConn = newConn

				go func() {
					connectionReady.Progress(1)
					assert.NoError(t,
						clientDisconnected.Wait(),
						"Client didn't disconnect",
					)

					assert.False(t,
						clientConn.IsActive(),
						"Expected connection to be inactive",
					)

					testerGoroutineFinished.Progress(1)
				}()
			},
			onClientDisconnected: func(_ wwr.Connection, _ error) {
				assert.False(t,
					clientConn.IsActive(),
					"Expected connection to be inactive",
				)

				// Try to send a signal to a inactive client and expect an error
				sigErr := clientConn.Signal(
					nil,
					wwr.Payload{
						Encoding: wwr.EncodingBinary,
						Data:     []byte("testdata"),
					},
				)
				assert.IsType(t, wwr.DisconnectedErr{}, sigErr)

				clientDisconnected.Progress(1)
			},
		},
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

	require.NoError(t, client.connection.Connect())

	// Wait for the connection to be set by the OnClientConnected handler
	require.NoError(t,
		connectionReady.Wait(),
		"Connection not ready after 1 second",
	)

	require.True(t,
		clientConn.IsActive(),
		"Expected connection to be active",
	)

	// Close the client connection and continue in the tester goroutine
	// spawned in the OnClientConnected handler of the server
	client.connection.Close()

	// Wait for the tester goroutine to finish
	require.NoError(t,
		testerGoroutineFinished.Wait(),
		"Tester goroutine didn't finish within 1 second",
	)
}
