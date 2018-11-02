package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestEmptyReplyUtf16 verifies empty UTF16 encoded reply acceptance
func TestEmptyReplyUtf16(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Return empty reply
				return wwr.NewPayload(wwr.EncodingUtf16, nil), nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.NewPayload(wwr.EncodingBinary, []byte("test")),
	)
	require.NoError(t, err)

	// Verify reply is empty
	require.Equal(t, wwr.EncodingUtf16, reply.Encoding())
	require.Len(t, reply.Data(), 0)
}
