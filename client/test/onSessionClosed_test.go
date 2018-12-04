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

// TestOnSessionClosed tests the OnSessionClosed hook
func TestOnSessionClosed(t *testing.T) {
	authenticated := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	hookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

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
				if err != nil {
					return wwr.Payload{}, err
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

				return wwr.Payload{}, nil
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
			OnSessionClosed: func() {
				hookCalled.Progress(1)
			},
		},
	)

	require.NoError(t, client.Connection.Connect())

	// Send authentication request and await reply
	_, err := client.Connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{
			Encoding: wwr.EncodingBinary,
			Data:     []byte("credentials"),
		},
	)
	require.NoError(t, err)
	authenticated.Progress(1)

	// Verify client session
	require.NoError(t, hookCalled.Wait(), "Hook not called")
}
