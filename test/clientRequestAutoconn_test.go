package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientRequestAutoconn tests sending requests on disconnected clients
// expecting it to connect automatically
func TestClientRequestAutoconn(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client and skip manual connection establishment
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	)
	require.NoError(t, err, "Expected client.Request to automatically connect")
	reply.Close()
}
