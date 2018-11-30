package test

import (
	"testing"
	"time"

	"github.com/qbeon/webwire-go/transport/memchan"

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
			onClientConnected: func(
				_ wwr.ConnectionOptions,
				conn wwr.Connection,
			) {
				assert.Equal(t, "samplestring", conn.Info(1).(string))
				assert.Equal(t, uint64(42), conn.Info(2).(uint64))
				assert.Nil(t, conn.Info(3))
				assert.WithinDuration(
					t,
					time.Now(),
					conn.Creation(),
					1*time.Second,
				)
				handlerFinished.Progress(1)
			},
		},
		wwr.ServerOptions{},
		&memchan.Transport{
			OnBeforeCreation: func() wwr.ConnectionOptions {
				return wwr.ConnectionOptions{
					Connection: wwr.Accept,
					Info: map[int]interface{}{
						1: "samplestring",
						2: uint64(42),
					},
				}
			},
		},
	)

	// Initialize client
	setup.newClient(
		wwrclt.Options{},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, handlerFinished.Wait())
}
