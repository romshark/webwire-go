package test

import (
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/assert"
)

// TestConnIsConnected tests the Connection.IsActive method as well as the
// OnClientConnected and OnClientDisconnected server hooks
func TestConnIsConnected(t *testing.T) {
	ready := sync.WaitGroup{}
	clientDisconnected := sync.WaitGroup{}
	finished := sync.WaitGroup{}
	ready.Add(1)
	clientDisconnected.Add(1)
	finished.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(
				_ wwr.ConnectionOptions,
				newConn wwr.Connection,
			) {
				assert.True(t, newConn.IsActive())

				go func() {
					ready.Done()
					clientDisconnected.Wait()

					assert.False(t, newConn.IsActive())

					finished.Done()
				}()
			},
			ClientDisconnected: func(conn wwr.Connection, _ error) {
				assert.False(t, conn.IsActive())
				clientDisconnected.Done()
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client, _ := setup.NewClientSocket()

	// Wait for the connection to be set by the OnClientConnected handler
	ready.Wait()

	// Close the client connection and continue in the tester goroutine
	// spawned in the OnClientConnected handler of the server
	client.Close()

	// Wait for the tester goroutine to finish
	finished.Wait()
}
