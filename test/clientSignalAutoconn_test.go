package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientSignalAutoconn tests sending signals on disconnected clients
// expecting it to connect automatically
func TestClientSignalAutoconn(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Skip manual connection establishment and rely on autoconnect instead
	require.NoError(t, client.connection.Signal(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	))
}
