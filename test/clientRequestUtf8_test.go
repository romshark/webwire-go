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

// TestClientRequestUtf8 tests requests with UTF8 encoded payloads
func TestClientRequestUtf8(t *testing.T) {
	expectedRequestPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("webwire_test_REQUEST_payload"),
	}
	expectedReplyPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("webwire_test_RESPONSE_message"),
	}

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
				assert.Equal(
					t,
					expectedRequestPayload.Encoding,
					msg.PayloadEncoding(),
				)
				assert.Equal(
					t,
					expectedRequestPayload.Data,
					msg.Payload(),
				)
				return expectedReplyPayload, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		expectedRequestPayload,
	)
	require.NoError(t, err)

	// Verify reply
	require.Equal(
		t,
		expectedReplyPayload.Encoding,
		reply.PayloadEncoding(),
	)
	require.Equal(
		t,
		expectedReplyPayload.Data,
		reply.Payload(),
	)
	reply.Close()
}
