package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRequestError tests server-side request errors properly
// failing the client-side requests
func TestClientRequestError(t *testing.T) {
	expectedReplyError := wwr.ReqErr{
		Code:    "SAMPLE_ERROR",
		Message: "Sample error message",
	}

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Fail the request by returning an error
				return nil, expectedReplyError
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
		"",
		wwr.NewPayload(
			wwr.EncodingUtf8,
			[]byte("webwire_test_REQUEST_payload"),
		),
	)

	// Verify returned error
	require.Error(t, err)
	require.IsType(t, wwr.ReqErr{}, err)
	require.Equal(t, err.(wwr.ReqErr).Code, expectedReplyError.Code)
	require.Equal(t, err.(wwr.ReqErr).Message, expectedReplyError.Message)
	require.Nil(t, reply)
}
