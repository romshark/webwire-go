package test

import (
	"errors"
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/assert"
)

// TestCustomSessKeyGenInvalid tests custom session key generators returning
// invalid keys
func TestCustomSessKeyGenInvalid(t *testing.T) {
	finished := sync.WaitGroup{}
	finished.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				defer func() {
					recoveredErr := recover()
					assert.NotNil(t, recoveredErr)
					assert.IsType(t, errors.New(""), recoveredErr)

					finished.Done()
				}()

				// Try to create a new session
				err := c.CreateSession(nil)
				assert.NoError(t, err)
			},
		},
		wwr.ServerOptions{
			SessionKeyGenerator: &SessionKeyGen{
				OnGenerate: func() string {
					// Return invalid session key
					return ""
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.NewClientSocket()

	finished.Wait()
}
