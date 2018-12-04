package test

import (
	"sync"
	"testing"
	"time"

	"github.com/qbeon/webwire-go/transport/memchan"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
)

// TestConnectionInfo tests the connection.Info method
func TestConnectionInfo(t *testing.T) {
	handlerFinished := sync.WaitGroup{}
	handlerFinished.Add(1)

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(
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
				handlerFinished.Done()
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
	setup.NewClient(
		wwrclt.Options{},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	handlerFinished.Wait()
}
