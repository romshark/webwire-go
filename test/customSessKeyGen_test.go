package test

import (
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/assert"
)

// TestCustomSessKeyGen tests custom session key generators
func TestCustomSessKeyGen(t *testing.T) {
	finished := sync.WaitGroup{}
	finished.Add(1)
	expectedSessionKey := "customkey123"

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				defer finished.Done()

				// Try to create a new session
				assert.NoError(t, c.CreateSession(nil))
				assert.Equal(t, expectedSessionKey, c.SessionKey())
			},
		},
		wwr.ServerOptions{
			SessionKeyGenerator: &SessionKeyGen{
				OnGenerate: func() string {
					return expectedSessionKey
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	readSessionCreated(t, sock)

	finished.Wait()
}
