package client_test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReqRegisterOnReply tests the request register of the client assuming it's
// correctly updated when the request is successfully fulfilled
func TestReqRegisterOnReply(t *testing.T) {
	var connection wwrclt.Client

	// Initialize webwire server given only the request
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{
			Request: func(
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
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		wwrtst.TestClientHooks{},
	)
	connection = client.Connection

	// Verify pending requests
	require.Equal(t, 0, client.Connection.PendingRequests())

	// Send request and await reply
	_, err := client.Connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("t")},
	)
	require.NoError(t, err)

	// Verify pending requests
	require.Equal(t, 0, client.Connection.PendingRequests())
}
