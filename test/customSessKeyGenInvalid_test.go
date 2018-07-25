package test

import (
	"context"
	"testing"
	"time"

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
				clt *wwr.Client,
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
				if err := clt.CreateSession(nil); err != nil {
					return nil, err
				}

				// Return the key of the newly created session
				// (use default binary encoding)
				return nil, nil
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
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	// Send authentication request and await reply
	if _, err := client.connection.Request(
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}
}
