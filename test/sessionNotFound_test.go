package test

import (
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/require"
)

// TestSessionNotFound tests restoration requests for inexistent sessions
// and expect them to fail returning the according error
func TestSessionNotFound(t *testing.T) {
	lookupTriggered := sync.WaitGroup{}
	lookupTriggered.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				SessionLookup: func(
					sessionKey string,
				) (wwr.SessionLookupResult, error) {
					lookupTriggered.Done()
					return nil, nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	// Skip manual connection establishment and rely on autoconnect instead
	reply := requestRestoreSession(t, sock, []byte("inexistentkey"))
	require.Equal(t, message.MsgReplySessionNotFound, reply.MsgType)

	lookupTriggered.Wait()
}
