package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
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
		server.Address(),
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
	require.IsType(t, wwr.SessionNotFoundErr{}, err)
}
