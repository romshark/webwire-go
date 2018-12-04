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

// TestRequestUtf8 tests requests with UTF8 encoded payloads
func TestRequestUtf8(t *testing.T) {
	expectedRequestPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("webwire_test_REQUEST_payload"),
	}
	expectedReplyPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("webwire_test_RESPONSE_message"),
	}

	// Initialize webwire server given only the request
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
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
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	// Send request and await reply
	reply, err := client.Connection.Request(
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
