package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestRequestNoNameNoPayload tests sending requests without both a name and a
// payload expecting the server to reject the message
func TestRequestNoNameNoPayload(t *testing.T) {
	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
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
		TestClientHooks{},
	)

	// TODO: improve test by avoiding the use of the client but performing the
	// request over a raw socket to ensure the client doesn't filter the request
	// out preemtively so it never even reaches the server

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
