package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientRequestError tests server-side request errors properly
// failing the client-side requests
func TestClientRequestError(t *testing.T) {
	expectedReplyError := wwr.RequestErr{
		Code:    "SAMPLE_ERROR",
		Message: "Sample error message",
	}

	// Initialize webwire server given only the request
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Fail the request by returning an error
				return wwr.Payload{}, expectedReplyError
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		testClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{
			Encoding: wwr.EncodingUtf8,
			Data:     []byte("webwire_test_REQUEST_payload"),
		},
	)

	// Verify returned error
	require.Error(t, err)
	require.IsType(t, wwr.RequestErr{}, err)
	require.Equal(t, err.(wwr.RequestErr).Code, expectedReplyError.Code)
	require.Equal(t, err.(wwr.RequestErr).Message, expectedReplyError.Message)
	require.Nil(t, reply)
}
