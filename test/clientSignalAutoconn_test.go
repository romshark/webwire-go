package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientSignalAutoconn tests sending signals on disconnected clients
// expecting it to connect automatically
func TestClientSignalAutoconn(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{},
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

	// Skip manual connection establishment and rely on autoconnect instead
	require.NoError(t, client.connection.Signal(
		"",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	))
}
