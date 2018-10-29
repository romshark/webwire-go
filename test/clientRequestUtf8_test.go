package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRequestUtf8 tests requests with UTF8 encoded payloads
func TestClientRequestUtf8(t *testing.T) {
	expectedRequestPayload := wwr.NewPayload(
		wwr.EncodingUtf8,
		[]byte("webwire_test_REQUEST_payload"),
	)
	expectedReplyPayload := wwr.NewPayload(
		wwr.EncodingUtf8,
		[]byte("webwire_test_RESPONSE_message"),
	)

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Verify request payload
				comparePayload(t, expectedRequestPayload, msg.Payload())
				return expectedReplyPayload, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		"",
		expectedRequestPayload,
	)
	require.NoError(t, err)

	// Verify reply
	comparePayload(t, expectedReplyPayload, reply)
}
