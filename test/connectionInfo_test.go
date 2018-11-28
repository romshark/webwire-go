package test

import (
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConnectionInfo tests the connection.Info method
func TestConnectionInfo(t *testing.T) {
	handlerFinished := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize server
	setup := setupTestServer(
		t,
		&serverImpl{
			onClientConnected: func(conn wwr.Connection) {
				info := conn.Info()
				assert.WithinDuration(
					t,
					time.Now(),
					info.ConnectionTime,
					1*time.Second,
				)
				handlerFinished.Progress(1)
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.newClient(
		wwrclt.Options{},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, handlerFinished.Wait())
}
