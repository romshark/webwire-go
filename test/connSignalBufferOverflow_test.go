package test

import (
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/assert"
)

// TestConnSignalBufferOverflow tests Connection.Signal with a name and payload
// that would overflow the buffer
func TestConnSignalBufferOverflow(t *testing.T) {
	finished := sync.WaitGroup{}
	finished.Add(1)

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				defer finished.Done()

				payload := make([]byte, 2048)
				err := c.Signal(
					[]byte(nil),                // No name
					wwr.Payload{Data: payload}, // Payload too big
				)

				assert.Error(t, err)
				assert.IsType(t, wwr.ErrBufferOverflow{}, err)
			},
		},
		wwr.ServerOptions{
			MessageBufferSize: 1024,
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.NewClientSocket()

	finished.Wait()
}
