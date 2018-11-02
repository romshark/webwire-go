package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRestSessAutoconn tests manual session restoration
// on disconnected client expecting client.RestoreSession
// to automatically establish a connection
func TestClientRestSessAutoconn(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 100 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
	)

	// Skip manual connection establishment and rely on autoconnect instead
	err := client.connection.RestoreSession(
		context.Background(),
		[]byte("inexistent_key"),
	)
	require.Error(t, err)
	require.IsType(t, wwr.SessNotFoundErr{}, err)
}
