package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionClosed tests the OnSessionClosed hook of the client
func TestClientOnSessionClosed(t *testing.T) {
	authenticated := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	hookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

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
				if err != nil {
					return nil, err
				}

				go func() {
					// Wait until the authentication request is finished
					assert.NoError(t,
						authenticated.Wait(),
						"Authentication timed out",
					)

					// Close the session
					assert.NoError(t, conn.CloseSession())
				}()

				return nil, nil
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
			OnSessionClosed: func() {
				hookCalled.Progress(1)
			},
		},
	)

	require.NoError(t, client.connection.Connect())

	// Send authentication request and await reply
	_, err := client.connection.Request(
		context.Background(),
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("credentials")),
	)
	require.NoError(t, err)
	authenticated.Progress(1)

	// Verify client session
	require.NoError(t, hookCalled.Wait(), "Hook not called")
}
