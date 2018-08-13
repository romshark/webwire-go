package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestRequestEmpty tests empty requests without a name and without a payload
func TestRequestEmpty(t *testing.T) {
	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect the following request to not even arrive
				t.Error("Not expected but reached")
				return nil, nil
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

	// Send request without a name and without a payload.
	// Expect a protocol error in return not sending the invalid request off
	_, err := client.connection.Request(context.Background(), "", nil)
	require.Error(t, err)
	require.IsType(t, wwr.ProtocolErr{}, err)
}
