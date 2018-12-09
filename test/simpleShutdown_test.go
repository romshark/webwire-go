package test

import (
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/require"
)

// TestSimpleShutdown tests simple shutdown without any pending tasks
func TestSimpleShutdown(t *testing.T) {
	connectedClientsNum := 5

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	clients := make([]wwr.Socket, connectedClientsNum)
	for i := 0; i < connectedClientsNum; i++ {
		sock, _ := setup.NewClientSocket()
		clients[i] = sock
	}

	require.NoError(t, setup.Server.Shutdown())
}
