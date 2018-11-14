package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientOnSessionCreated tests the OnSessionCreated hook of the client
func TestClientOnSessionCreated(t *testing.T) {
	hookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	var createdSession *wwr.Session
	var sessionFromHook *wwr.Session

	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
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
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{
			OnSessionCreated: func(newSession *wwr.Session) {
				sessionFromHook = newSession
				hookCalled.Progress(1)
			},
		},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	// Send authentication request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{
			Encoding: wwr.EncodingBinary,
			Data:     []byte("credentials"),
		},
	)
	require.NoError(t, err)
	reply.Close()

	createdSession = client.connection.Session()

	// Verify client session
	require.NoError(t, hookCalled.Wait(), "Hook not called")

	// Compare the actual created session with the session received in the hook
	compareSessions(t, createdSession, sessionFromHook)
}
