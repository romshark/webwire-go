package test

import (
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"

	"github.com/qbeon/webwire-go"
	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/require"
)

// TestMaxConcSessConn tests 4 maximum concurrent connections of a session
func TestMaxConcSessConn(t *testing.T) {
	concurrentConns := uint(4)

	var sessionKey = "testsessionkey"
	sessionCreation := time.Now()

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{
			MaxSessionConnections: concurrentConns,
			SessionManager: &SessionManager{
				SessionLookup: func(key string) (
					webwire.SessionLookupResult,
					error,
				) {
					if key != sessionKey {
						// Session not found
						return nil, nil
					}
					return webwire.NewSessionLookupResult(
						sessionCreation, // Creation
						time.Now(),      // LastLookup
						nil,             // Info
					), nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize clients
	clients := make([]wwr.Socket, concurrentConns)
	for i := uint(0); i < concurrentConns; i++ {
		sock, _ := setup.NewClientSocket()
		clients[i] = sock

		requestRestoreSessionSuccess(t, sock, []byte(sessionKey))
	}

	// Ensure that the last superfluous client is rejected
	superfluousClient, _ := setup.NewClientSocket()

	// Try to restore the session one more time and expect this request to fail
	// due to reached limit
	reply := requestRestoreSession(t, superfluousClient, []byte(sessionKey))
	require.Equal(t, message.MsgMaxSessConnsReached, reply.MsgType)
}
