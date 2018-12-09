package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequestPayloadOnly tests requests without a name but only a payload
func TestRequestPayloadOnly(t *testing.T) {
	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect a named request
				msgName := msg.Name()
				assert.Nil(t, msgName)

				switch msg.PayloadEncoding() {
				case wwr.EncodingBinary:
					require.Equal(t, []byte("d"), msg.Payload())
					return wwr.Payload{Encoding: wwr.EncodingBinary}, nil
				case wwr.EncodingUtf8:
					require.Equal(t, []byte("d"), msg.Payload())
					return wwr.Payload{Encoding: wwr.EncodingUtf8}, nil
				case wwr.EncodingUtf16:
					require.Equal(t, []byte{32, 32}, msg.Payload())
					return wwr.Payload{Encoding: wwr.EncodingUtf16}, nil
				default:
					panic("unexpected message payload encoding type")
				}
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	// Send an unnamed binary request with a payload and await reply
	requestSuccess(t, sock, 32, nil, payload.Payload{Data: []byte("d")})

	// Send an unnamed UTF8 encoded request with a payload
	requestSuccess(t, sock, 32, nil, payload.Payload{
		Encoding: payload.Utf8,
		Data:     []byte("d"),
	})

	// Send an unnamed UTF16 encoded request with a payload
	requestSuccess(t, sock, 32, nil, payload.Payload{
		Encoding: payload.Utf16,
		Data:     []byte{32, 32},
	})
}
