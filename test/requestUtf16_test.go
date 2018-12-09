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

// TestRequestUtf16 tests requests with UTF16 encoded payloads
func TestRequestUtf16(t *testing.T) {
	testPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf16,
		Data:     []byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
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
				assert.Equal(t, wwr.EncodingUtf16, msg.PayloadEncoding())
				assert.Equal(t, testPayload.Data, msg.Payload())

				return wwr.Payload{
					Encoding: wwr.EncodingUtf16,
					Data: []byte{
						00, 115, 00, 97, 00, 109,
						00, 112, 00, 108, 00, 101,
					},
				}, nil
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
		wwr.Payload{
			Encoding: wwr.EncodingUtf16,
			Data: []byte{
				00, 115, 00, 97, 00, 109,
				00, 112, 00, 108, 00, 101,
			},
		},
	)
	require.NoError(t, err)

	// Verify reply
	assert.Equal(t, wwr.EncodingUtf16, reply.PayloadEncoding())
	assert.Equal(t, testPayload.Data, reply.Payload())
	reply.Close()
}