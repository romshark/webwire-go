package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestRestoreInexistentSession tests the restoration of an inexistent session
func TestRestoreInexistentSession(t *testing.T) {
	// Initialize server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client

	// Ensure that the last superfluous client is rejected
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
		nil, // No TLS configuration
	)

	require.NoError(t, client.connection.Connect())

	// Try to restore the session and expect it to fail
	// due to the session being inexistent
	sessionRestorationError := client.connection.RestoreSession(
		[]byte("lalala"),
	)
	require.Error(t, sessionRestorationError)
	require.IsType(t, wwr.SessNotFoundErr{}, sessionRestorationError)
}
