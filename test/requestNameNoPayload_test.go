package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
)

// TestRequestNameNoPayload tests named requests without a payload
func TestRequestNameNoPayload(t *testing.T) {
	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ webwire.Connection,
				msg webwire.Message,
			) (webwire.Payload, error) {
				// Expect a named request
				msgName := msg.Name()
				assert.Equal(t, "n", msgName)

				// Expect no payload to arrive
				assert.Len(t, msg.Payload().Data(), 0)

				return nil, nil
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Send a named binary request without a payload and await reply
	_, err := client.connection.Request(
		context.Background(),
		"n",
		webwire.NewPayload(webwire.EncodingBinary, nil),
	)
	require.NoError(t, err)

	// Send a UTF16 encoded named binary request without a payload
	_, err = client.connection.Request(
		context.Background(),
		"n",
		webwire.NewPayload(webwire.EncodingUtf16, nil),
	)
	require.NoError(t, err)
}
