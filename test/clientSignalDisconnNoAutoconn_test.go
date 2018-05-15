package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientSignalDisconnectedErr tests client.Signal
// expecting it to return a DisconnectedErr when autoconn is disabled
// and the client is disconnected
func TestClientSignalDisconnectedErr(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			// Disable autoconnect to prevent automatic reconnection
			Autoconnect: wwrclt.OptDisabled,
		},
		callbackPoweredClientHooks{},
	)

	err := client.connection.Signal("", wwr.Payload{Data: []byte("test")})
	if err == nil {
		t.Fatalf("Expected DisconnectedErr, got nil")
	}
	if err, converted := err.(wwr.DisconnectedErr); !converted {
		t.Fatalf("Expected DisconnectedErr, got: %s", err)
	}
}
