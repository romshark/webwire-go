package test

import (
	"sync"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientDisconnectedHook verifies the server is calling the
// onClientDisconnected hook properly
func TestClientDisconnectedHook(t *testing.T) {
	disconnectedHookCalled := newPending(1, 1*time.Second, true)
	var connectedClient *webwire.Client
	connectedClientLock := sync.Mutex{}

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(clt *webwire.Client) {
				connectedClientLock.Lock()
				connectedClient = clt
				connectedClientLock.Unlock()
			},
			onClientDisconnected: func(clt *webwire.Client) {
				connectedClientLock.Lock()
				if clt != connectedClient {
					t.Errorf(
						"Connected and disconnecting clients don't match: "+
							"disconnecting: %p | connected: %p",
						clt,
						connectedClient,
					)
				}
				connectedClientLock.Unlock()
				disconnectedHookCalled.Done()
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Connect to the server
	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect the client to the server: %s", err)
	}

	// Disconnect the client
	client.connection.Close()

	// Await the onClientDisconnected hook to be called on the server
	if err := disconnectedHookCalled.Wait(); err != nil {
		t.Fatal("server.OnClientDisconnected hook not called")
	}
}
