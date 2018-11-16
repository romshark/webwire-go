package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequestNameNoPayload tests named requests without a payload
func TestRequestNameNoPayload(t *testing.T) {
	// Initialize server
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ webwire.Connection,
				msg webwire.Message,
			) (webwire.Payload, error) {
				// Expect a named request
				msgName := msg.Name()
				assert.Equal(t, []byte("n"), msgName)

				// Expect no payload to arrive
				assert.Len(t, msg.Payload(), 0)

				return webwire.Payload{}, nil
			},
		},
		webwire.ServerOptions{},
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

	// Send a named binary request without a payload and await reply
	_, err := client.connection.Request(
		context.Background(),
		[]byte("n"),
		webwire.Payload{
			Encoding: webwire.EncodingBinary,
			Data:     nil,
		},
	)
	require.NoError(t, err)

	// Send a UTF16 encoded named binary request without a payload
	_, err = client.connection.Request(
		context.Background(),
		[]byte("n"),
		webwire.Payload{
			Encoding: webwire.EncodingUtf16,
			Data:     nil,
		},
	)
	require.NoError(t, err)
}
