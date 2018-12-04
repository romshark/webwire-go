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

// TestRestoreInexistentSession tests the restoration of an inexistent session
func TestRestoreInexistentSession(t *testing.T) {
	// Initialize server
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{},
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

	// Try to restore the session and expect it to fail
	// due to the session being inexistent
	sessionRestorationError := client.Connection.RestoreSession(
		context.Background(),
		[]byte("lalala"),
	)
	require.Error(t, sessionRestorationError)
	require.IsType(t, wwr.SessionNotFoundErr{}, sessionRestorationError)
}
