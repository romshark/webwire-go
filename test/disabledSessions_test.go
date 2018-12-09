package test

import (
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDisabledSessions tests errors returned by CreateSession, CloseSession and
// client.RestoreSession when sessions are disabled
func TestDisabledSessions(t *testing.T) {
	finished := sync.WaitGroup{}
	finished.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				assert.Nil(t, c.Session())

				// Try to create a new session
				createErr := c.CreateSession(nil)
				assert.IsType(t, wwr.ErrSessionsDisabled{}, createErr)

				// Try to close a session
				closeErr := c.CloseSession()
				assert.IsType(t, wwr.ErrSessionsDisabled{}, closeErr)

				finished.Done()
			},
		},
		wwr.ServerOptions{
			Sessions: wwr.Disabled,
			SessionManager: &SessionManager{
				SessionCreated: func(c wwr.Connection) error {
					t.Fatal("unexpected hook call")
					return nil
				},
				SessionLookup: func(
					sessionKey string,
				) (wwr.SessionLookupResult, error) {
					t.Fatal("unexpected hook call")
					return nil, nil
				},
				SessionClosed: func(sessionKey string) error {
					t.Fatal("unexpected hook call")
					return nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	finished.Wait()

	// Try to restore a session
	reply := requestRestoreSession(t, sock, []byte("testsessionkey"))
	require.Equal(t, message.MsgReplySessionsDisabled, reply.MsgType)
}
