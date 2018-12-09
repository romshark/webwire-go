package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequestUtf16 tests requests with UTF16 encoded payloads
func TestRequestUtf16(t *testing.T) {
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
				assert.Equal(t, wwr.EncodingUtf16, msg.PayloadEncoding())
				assert.Equal(t, []byte{11, 20, 31, 40, 51, 60}, msg.Payload())
				return wwr.Payload{
					Encoding: wwr.EncodingUtf16,
					Data:     []byte{80, 91, 100, 111, 120, 131},
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
		Encoding: wwr.EncodingUtf16,
		Data:     []byte{11, 20, 31, 40, 51, 60},
	})

	// Verify reply
	require.Equal(t, wwr.EncodingUtf16, reply.MsgPayload.Encoding)
	require.Equal(t, []byte{80, 91, 100, 111, 120, 131}, reply.Payload())
}
