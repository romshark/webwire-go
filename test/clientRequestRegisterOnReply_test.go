package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRequestRegisterOnReply verifies the request register of the client
// is correctly updated when the request is successfully fulfilled
func TestClientRequestRegisterOnReply(t *testing.T) {
	var connection wwrclt.Client

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Verify pending requests
				assert.Equal(t, 1, connection.PendingRequests())

				// Wait until the request times out
				time.Sleep(300 * time.Millisecond)
				return nil, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)
	connection = client.connection

	// Connect the client to the server
	require.NoError(t, client.connection.Connect())

	// Verify pending requests
	require.Equal(t, 0, client.connection.PendingRequests())

	// Send request and await reply
	_, err := client.connection.Request(
		context.Background(),
		"",
		wwr.NewPayload(wwr.EncodingBinary, []byte("t")),
	)
	require.NoError(t, err)

	// Verify pending requests
	require.Equal(t, 0, client.connection.PendingRequests())
}
