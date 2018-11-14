package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientRequestRegisterOnReply verifies the request register of the client
// is correctly updated when the request is successfully fulfilled
func TestClientRequestRegisterOnReply(t *testing.T) {
	var connection wwrclt.Client

	// Initialize webwire server given only the request
	setup := setupTestServer(
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
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		testClientHooks{},
	)
	connection = client.connection

	// Connect the client to the server
	require.NoError(t, client.connection.Connect())

	// Verify pending requests
	require.Equal(t, 0, client.connection.PendingRequests())

	// Send request and await reply
	_, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("t")},
	)
	require.NoError(t, err)

	// Verify pending requests
	require.Equal(t, 0, client.connection.PendingRequests())
}
