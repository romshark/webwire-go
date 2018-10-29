package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestRequestNoNameOnlyPayload tests requests without a name but only a payload
func TestRequestNoNameOnlyPayload(t *testing.T) {
	expectedRequestPayload := wwr.NewPayload(
		wwr.EncodingUtf8,
		[]byte("3"),
	)
	expectedRequestPayloadUtf16 := wwr.NewPayload(
		wwr.EncodingUtf16,
		[]byte("12"),
	)

	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect a named request
				msgName := msg.Name()
				assert.Equal(t, "", msgName)

				msgPayload := msg.Payload()
				if msgPayload.Encoding() == wwr.EncodingUtf16 {
					comparePayload(t, expectedRequestPayloadUtf16, msgPayload)
				} else {
					comparePayload(t, expectedRequestPayload, msgPayload)
				}

				return nil, nil
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

	// Send an unnamed binary request with a payload and await reply
	_, err := client.connection.Request(
		context.Background(),
		"",
		expectedRequestPayload,
	)
	require.NoError(t, err)

	// Send an unnamed UTF16 encoded binary request with a payload
	_, err = client.connection.Request(
		context.Background(),
		"",
		expectedRequestPayloadUtf16,
	)
	require.NoError(t, err)
}
