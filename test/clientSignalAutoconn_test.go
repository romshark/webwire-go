package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientSignalAutoconn tests sending signals on disconnected clients
// expecting it to connect automatically
func TestClientSignalAutoconn(t *testing.T) {
	// Initialize webwire server given only the request
	setup := setupTestServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		testClientHooks{},
	)

	// Skip manual connection establishment and rely on autoconnect instead
	require.NoError(t, client.connection.Signal(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	))
}
