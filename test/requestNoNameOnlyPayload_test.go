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

// TestRequestNoNameOnlyPayload tests requests without a name but only a payload
func TestRequestNoNameOnlyPayload(t *testing.T) {
	expectedRequestPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("3"),
	}
	expectedRequestPayloadUtf16 := wwr.Payload{
		Encoding: wwr.EncodingUtf16,
		Data:     []byte("12"),
	}

	// Initialize server
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect a named request
				msgName := msg.Name()
				assert.Nil(t, msgName)

				if msg.PayloadEncoding() == wwr.EncodingUtf16 {
					require.Equal(
						t,
						expectedRequestPayloadUtf16.Encoding,
						msg.PayloadEncoding(),
					)
					require.Equal(
						t,
						expectedRequestPayloadUtf16.Data,
						msg.Payload(),
					)
				} else {
					require.Equal(
						t,
						expectedRequestPayload.Encoding,
						msg.PayloadEncoding(),
					)
					require.Equal(
						t,
						expectedRequestPayload.Data,
						msg.Payload(),
					)
				}

				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	// Send an unnamed binary request with a payload and await reply
	_, err := client.connection.Request(
		context.Background(),
		nil,
		expectedRequestPayload,
	)
	require.NoError(t, err)

	// Send an unnamed UTF16 encoded binary request with a payload
	_, err = client.connection.Request(
		context.Background(),
		nil,
		expectedRequestPayloadUtf16,
	)
	require.NoError(t, err)
}
