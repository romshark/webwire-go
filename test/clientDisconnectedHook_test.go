package test

import (
	"os"
	"sync"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
)

// TestClientDisconnectedHook verifies the server is calling the
// onClientDisconnected hook properly
func TestClientDisconnectedHook(t *testing.T) {
	var disconnectedHook sync.WaitGroup
	disconnectedHook.Add(1)
	var connectedClient *webwire.Client

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		webwire.Hooks{
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
				disconnectedHook.Done()
			},
		},
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil, nil, nil,
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	// Connect to the server
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect the client to the server: %s", err)
	}

	// Disconnect the client
	client.Close()

	// Await the onClientDisconnected hook to be called on the server
	disconnectedHook.Wait()
}
