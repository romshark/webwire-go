package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientSignalDisconnectedErr tests client.Signal
// expecting it to return a DisconnectedErr when autoconn is disabled
// and the client is disconnected
func TestClientSignalDisconnectedErr(t *testing.T) {
	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			// Disable autoconnect to prevent automatic reconnection
			Autoconnect: wwr.Disabled,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	// Try to send a signal and expect a DisconnectedErr error
	err := client.connection.Signal(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("test")},
	)
	require.Error(t, err)
	require.IsType(t, wwr.DisconnectedErr{}, err)
}
