package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRequestUtf16 tests requests with UTF16 encoded payloads
func TestClientRequestUtf16(t *testing.T) {
	testPayload := wwr.NewPayload(
		wwr.EncodingUtf16,
		[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	)
	verifyPayload := func(payload wwr.Payload) {
		assert.Equal(t, wwr.EncodingUtf16, payload.Encoding())
		assert.Equal(t, testPayload.Data(), payload.Data())
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
				verifyPayload(msg.Payload())

				return wwr.NewPayload(
					wwr.EncodingUtf16,
					[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
				), nil
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

	require.NoError(t, client.connection.Connect())

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.NewPayload(
			wwr.EncodingUtf16,
			[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
		),
	)
	require.NoError(t, err)

	// Verify reply
	verifyPayload(reply)
}
