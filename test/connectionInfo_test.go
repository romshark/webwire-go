package test

import (
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/transport/memchan"
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
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				defer handlerFinished.Done()
				assert.Equal(t, "samplestring", c.Info(1).(string))
				assert.Equal(t, uint64(42), c.Info(2).(uint64))
				assert.Nil(t, c.Info(3))
				assert.WithinDuration(
					t,
					time.Now(),
					c.Creation(),
					1*time.Second,
				)
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
	setup.NewClientSocket()

	handlerFinished.Wait()
}
