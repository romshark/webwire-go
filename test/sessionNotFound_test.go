package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestSessionNotFound tests restoration requests for inexistent sessions
// and expect them to fail returning the according error
func TestSessionNotFound(t *testing.T) {
	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 100 * time.Millisecond,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	// Skip manual connection establishment and rely on autoconnect instead
	err := client.Connection.RestoreSession(
		context.Background(),
		[]byte("inexistent_key"),
	)
	require.Error(t, err)
	require.IsType(t, wwr.SessionNotFoundErr{}, err)
}
