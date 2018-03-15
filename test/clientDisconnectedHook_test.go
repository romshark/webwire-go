package test

import (
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientDisconnectedHook verifies the server is calling the
// onClientDisconnected hook properly
func TestClientDisconnectedHook(t *testing.T) {
	disconnectedHookCalled := NewPending(1, 1*time.Second, true)
	var connectedClient *webwire.Client

	// Initialize webwire server given only the request
	_, addr := setupServer(
		t,
		webwire.ServerOptions{
			Hooks: webwire.Hooks{
				OnClientConnected: func(clt *webwire.Client) {
					connectedClient = clt
				},
				OnClientDisconnected: func(clt *webwire.Client) {
					if clt != connectedClient {
						t.Errorf(
							"Connected and disconnecting clients don't match: "+
								"disconnecting: %p | connected: %p",
							clt,
							connectedClient,
						)
					}
					disconnectedHookCalled.Done()
				},
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	// Connect to the server
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect the client to the server: %s", err)
	}

	// Disconnect the client
	client.Close()

	// Await the onClientDisconnected hook to be called on the server
	if err := disconnectedHookCalled.Wait(); err != nil {
		t.Fatal("server.OnClientDisconnected hook not called")
	}
}
