package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientRequestDisconnected tests sending requests on disconnected clients
func TestClientRequestDisconnected(t *testing.T) {
	// Initialize webwire server given only the request
	setup := setupTestServer(
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
		nil, // Use the default transport implementation
	)

	// Initialize client and skip manual connection establishment
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			Autoconnect:           wwrclt.StatusDisabled,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	)
	require.NoError(t, err)
	reply.Close()
}
