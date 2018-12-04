package client_test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/require"
)

// TestRequestNoNameNoPayload tests sending requests without both a name and a
// payload expecting the client to not even send the request message
func TestRequestNoNameNoPayload(t *testing.T) {
	// Initialize server
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{
			Request: func(
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
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		wwrtst.TestClientHooks{},
	)

	// Send request without a name and without a payload.
	// Expect a protocol error in return not sending the invalid request off
	reply, err := client.Connection.Request(
		context.Background(),
		nil,
		wwr.Payload{},
	)
	require.Error(t, err)
	require.IsType(t, wwr.ProtocolErr{}, err)
	require.Nil(t, reply)
}
