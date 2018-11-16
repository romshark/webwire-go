package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCustomSessKeyGenInvalid tests custom session key generators
// returning invalid keys
func TestCustomSessKeyGenInvalid(t *testing.T) {
	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				defer func() {
					if err := recover(); err == nil {
						t.Errorf("Expected server to panic " +
							"on invalid session key",
						)
					}
				}()

				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				return wwr.Payload{}, err
			},
		},
		wwr.ServerOptions{
			SessionKeyGenerator: &sessionKeyGen{
				generate: func() string {
					// Return invalid session key
					return ""
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)
	defer client.connection.Close()

	// Send authentication request and await reply
	_, err := client.connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{Data: []byte("testdata")},
	)
	require.NoError(t, err)
}
