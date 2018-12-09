package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequestUtf8 tests requests with UTF8 encoded payloads
func TestRequestUtf8(t *testing.T) {
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
				assert.Equal(t, wwr.EncodingUtf8, msg.PayloadEncoding())
				assert.Equal(t, []byte("sample data"), msg.Payload())
				return wwr.Payload{
					Encoding: wwr.EncodingUtf8,
					Data:     []byte("sample reply"),
				}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	// Send request and await reply
	reply := requestSuccess(t, sock, 32, nil, payload.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("sample data"),
	})

	// Verify reply
	require.Equal(t, wwr.EncodingUtf8, reply.MsgPayload.Encoding)
	require.Equal(t, []byte("sample reply"), reply.Payload())
}
