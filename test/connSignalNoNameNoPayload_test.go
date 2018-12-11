package test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	wwr "github.com/qbeon/webwire-go"
)

// TestConnSignalNoNameNoPayload tests Connection.Signal providing both a nil
// name and a nil payload
func TestConnSignalNoNameNoPayload(t *testing.T) {
	finished := sync.WaitGroup{}
	finished.Add(1)

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				defer finished.Done()

				assert.Error(t, c.Signal(
					[]byte(nil),                    // No name
					wwr.Payload{Data: []byte(nil)}, // No payload
				))
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.NewClientSocket()

	finished.Wait()
}
