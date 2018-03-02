package test

import (
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientIsConnected verifies correct client.IsConnected reporting
func TestClientIsConnected(t *testing.T) {
	// Initialize webwire server given only the request
	_, addr := setupServer(
		t,
		webwire.Hooks{},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	if client.IsConnected() {
		t.Fatal("Expected client to be disconnected before the connection establishment")
	}

	// Connect to the server
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect the client to the server: %s", err)
	}

	if !client.IsConnected() {
		t.Fatal("Expected client to be connected after the connection establishment")
	}

	// Disconnect the client
	client.Close()

	if client.IsConnected() {
		t.Fatal("Expected client to be disconnected after closure")
	}
}
