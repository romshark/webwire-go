package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestEmptyReply tests returning empty binary replies from the request handler
func TestEmptyReply(t *testing.T) {
	// Initialize webwire server given only the request
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
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
	client := setup.NewClient(
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	// Send request and await reply
	reply, err := client.Connection.Request(
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
