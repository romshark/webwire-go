package test

import (
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/assert"
)

// TestSessionCreationOnClosedConn tests the creation of a session on a
// disconnected connection
func TestSessionCreationOnClosedConn(t *testing.T) {
	onConnectedFinished := sync.WaitGroup{}
	onConnectedFinished.Add(1)
	onDisconnectedFinished := sync.WaitGroup{}
	onDisconnectedFinished.Add(1)

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(
				_ wwr.ConnectionOptions,
				conn wwr.Connection,
			) {
				defer onConnectedFinished.Done()
				conn.Close()

				err := conn.CreateSession(nil)
				assert.Error(t, err)
				assert.IsType(t, wwr.DisconnectedErr{}, err)
			},
			ClientDisconnected: func(conn wwr.Connection, _ error) {
				defer onDisconnectedFinished.Done()
				err := conn.CreateSession(nil)
				assert.Error(t, err)
				assert.IsType(t, wwr.DisconnectedErr{}, err)
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.NewClientSocket()

	onConnectedFinished.Wait()
	onDisconnectedFinished.Wait()
}
