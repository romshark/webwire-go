package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestConnectionInfo tests the connection.Info method
func TestConnectionInfo(t *testing.T) {
	// Initialize server
	server := setupServer(
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
				assert.Equal(t, "Go-http-client/1.1", info.UserAgent)
				assert.NotNil(t, info.RemoteAddr)
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())
}
