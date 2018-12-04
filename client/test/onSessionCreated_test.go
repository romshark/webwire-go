package client_test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOnSessionCreated tests the OnSessionCreated hook
func TestOnSessionCreated(t *testing.T) {
	hookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	var createdSession *wwr.Session
	var sessionFromHook *wwr.Session

	// Initialize webwire server
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				return wwr.Payload{}, err
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		wwrtst.TestClientHooks{
			OnSessionCreated: func(newSession *wwr.Session) {
				sessionFromHook = newSession
				hookCalled.Progress(1)
			},
		},
	)
	defer client.Connection.Close()

	require.NoError(t, client.Connection.Connect())

	// Send authentication request and await reply
	reply, err := client.Connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{
			Encoding: wwr.EncodingBinary,
			Data:     []byte("credentials"),
		},
	)
	require.NoError(t, err)
	reply.Close()

	createdSession = client.Connection.Session()

	// Verify client session
	require.NoError(t, hookCalled.Wait(), "Hook not called")

	// Compare the actual created session with the session received in the hook
	wwrtst.CompareSessions(t, createdSession, sessionFromHook)
}
