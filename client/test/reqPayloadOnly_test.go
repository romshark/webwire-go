package client_test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReqPayloadOnly tests requests without a name but only a payload
func TestReqPayloadOnly(t *testing.T) {
	expectedRequestPayloadBin := wwr.Payload{
		Encoding: wwr.EncodingBinary,
		Data:     []byte("3"),
	}
	expectedRequestPayloadUtf8 := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("6"),
	}
	expectedRequestPayloadUtf16 := wwr.Payload{
		Encoding: wwr.EncodingUtf16,
		Data:     []byte("19"),
	}

	// Initialize server
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{
			Request: func(
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
						expectedRequestPayloadUtf16.Data,
						msg.Payload(),
					)
				} else if msg.PayloadEncoding() == wwr.EncodingUtf8 {
					require.Equal(
						t,
						expectedRequestPayloadUtf8.Data,
						msg.Payload(),
					)
				} else {
					require.Equal(
						t,
						expectedRequestPayloadBin.Data,
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
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		wwrtst.TestClientHooks{},
	)

	// Send an unnamed binary request with a payload and await reply
	_, err := client.Connection.Request(
		context.Background(),
		nil,
		expectedRequestPayloadBin,
	)
	require.NoError(t, err)

	// Send an unnamed UTF8 encoded request with a payload and await reply
	_, err = client.Connection.Request(
		context.Background(),
		nil,
		expectedRequestPayloadUtf8,
	)
	require.NoError(t, err)

	// Send an unnamed UTF16 encoded binary request with a payload
	_, err = client.Connection.Request(
		context.Background(),
		nil,
		expectedRequestPayloadUtf16,
	)
	require.NoError(t, err)
}
