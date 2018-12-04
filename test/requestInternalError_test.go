package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestRequestInternalError tests returning non-ReqErr errors from the request
// handler
func TestRequestInternalError(t *testing.T) {
	// Initialize webwire server given only the request
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Fail the request by returning a non-ReqErr error
				return wwr.Payload{}, fmt.Errorf(
					"don't worry, this internal error is expected",
				)
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
		TestClientHooks{},
	)

	// Send request and await reply
	reply, reqErr := client.Connection.Request(
		context.Background(),
		[]byte("test"),
		wwr.Payload{},
	)

	// Verify returned error
	require.Error(t, reqErr)
	require.IsType(t, wwr.InternalErr{}, reqErr)
	require.Nil(t, reply)
}
