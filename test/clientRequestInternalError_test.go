package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRequestInternalError tests returning of non-ReqErr errors
// from the request handler
func TestClientRequestInternalError(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Fail the request by returning a non-ReqErr error
				return nil, fmt.Errorf(
					"don't worry, this internal error is expected",
				)
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send request and await reply
	reply, reqErr := client.connection.Request(
		context.Background(),
		"",
		wwr.NewPayload(wwr.EncodingUtf8, []byte("dummydata")),
	)

	// Verify returned error
	require.Error(t, reqErr)
	require.IsType(t, wwr.ReqInternalErr{}, reqErr)
	require.Nil(t, reply)
}
