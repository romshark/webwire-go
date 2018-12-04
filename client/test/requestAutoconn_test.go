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

// TestRequestAutoconn tests sending requests on disconnected clients expecting
// it to automatically establish a connection
func TestRequestAutoconn(t *testing.T) {
	// Initialize webwire server given only the request
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client and skip manual connection establishment
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			Autoconnect:           wwr.Enabled,
		},
		nil, // Use the default transport implementation
		wwrtst.TestClientHooks{},
	)

	// Send request and await reply
	reply, err := client.Connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	)
	require.NoError(t, err, "Expected client.Request to automatically connect")
	reply.Close()
}
