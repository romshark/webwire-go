package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestCustomSessKeyGenInvalid tests custom session key generators
// returning invalid keys
func TestCustomSessKeyGenInvalid(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
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
				return nil, err
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
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
		nil, // No TLS configuration
	)
	defer client.connection.Close()

	// Send authentication request and await reply
	_, err := client.connection.Request(
		context.Background(),
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	)
	require.NoError(t, err)
}
