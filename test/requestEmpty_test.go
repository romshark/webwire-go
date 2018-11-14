package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestRequestEmpty tests empty requests without a name and without a payload
func TestRequestEmpty(t *testing.T) {
	// Initialize server
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect the following request to not even arrive
				t.Error("Not expected but reached")
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

	// Send request without a name and without a payload.
	// Expect a protocol error in return not sending the invalid request off
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{},
	)
	require.Error(t, err)
	require.IsType(t, wwr.ProtocolErr{}, err)
	require.Nil(t, reply)
}
