package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
)

// TestRequestNameOnly tests named requests without a payload
func TestRequestNameOnly(t *testing.T) {
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
				assert.Equal(t, []byte("name"), msg.Name())

				// Expect no payload to arrive
				assert.Equal(t, 0, len(msg.Payload()))

				switch msg.PayloadEncoding() {
				case wwr.EncodingUtf8:
					return wwr.Payload{Encoding: wwr.EncodingUtf8}, nil
				case wwr.EncodingUtf16:
					return wwr.Payload{Encoding: wwr.EncodingUtf16}, nil
				}
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	requestSuccess(t, sock, 32, []byte("name"), payload.Payload{})

	// Send a named UTF8 encoded request without a payload and await reply
	requestSuccess(t, sock, 32, []byte("name"), payload.Payload{
		Encoding: payload.Utf8,
	})

	// Send a UTF16 encoded named binary request without a payload
	requestSuccess(t, sock, 32, []byte("name"), payload.Payload{
		Encoding: payload.Utf16,
	})
}
