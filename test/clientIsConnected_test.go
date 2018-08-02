package test

import (
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientIsConnected tests the client.Status method
func TestClientIsConnected(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(t, &serverImpl{}, webwire.ServerOptions{})

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	if client.connection.Status() == webwireClient.Connected {
		t.Fatal("Expected client to be disconnected before the connection establishment")
	}

	// Connect to the server
	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect the client to the server: %s", err)
	}

	if client.connection.Status() != webwireClient.Connected {
		t.Fatal("Expected client to be connected after the connection establishment")
	}

	// Disconnect the client
	client.connection.Close()

	if client.connection.Status() == webwireClient.Connected {
		t.Fatal("Expected client to be disconnected after closure")
	}
}
