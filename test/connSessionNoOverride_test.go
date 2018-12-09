package test

import (
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/assert"
)

// TestSessionNoOverride tests overriding of a connection session
// expecting it to fail
func TestSessionNoOverride(t *testing.T) {
	finished := sync.WaitGroup{}
	finished.Add(1)

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				defer finished.Done()

				assert.NoError(t, c.CreateSession(nil))
				sessionKey := c.SessionKey()

				// Try to override the previous session
				assert.Error(t, c.CreateSession(nil))

				// Ensure the session didn't change
				assert.Equal(t, sessionKey, c.SessionKey())
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	readSessionCreated(t, sock)

	finished.Wait()
}
