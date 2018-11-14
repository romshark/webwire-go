package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestEmptyReply verifies empty binary reply acceptance
func TestEmptyReply(t *testing.T) {
	// Initialize webwire server given only the request
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Return empty reply
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("test")},
	)
	require.NoError(t, err)

	// Verify reply is empty
	require.Equal(t, wwr.EncodingBinary, reply.PayloadEncoding())
	require.Len(t, reply.Payload(), 0)
	reply.Close()
}
