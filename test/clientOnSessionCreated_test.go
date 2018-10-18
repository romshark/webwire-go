package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionCreated tests the OnSessionCreated hook of the client
func TestClientOnSessionCreated(t *testing.T) {
	hookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	var createdSession *wwr.Session
	var sessionFromHook *wwr.Session

	// Initialize webwire server
	server := setupServer(
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
				return nil, err
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{
			OnSessionCreated: func(newSession *wwr.Session) {
				sessionFromHook = newSession
				hookCalled.Progress(1)
			},
		},
		nil, // No TLS configuration
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	// Send authentication request and await reply
	_, err := client.connection.Request(
		context.Background(),
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("credentials")),
	)
	require.NoError(t, err)

	createdSession = client.connection.Session()

	// Verify client session
	require.NoError(t, hookCalled.Wait(), "Hook not called")

	// Compare the actual created session with the session received in the hook
	compareSessions(t, createdSession, sessionFromHook)
}
