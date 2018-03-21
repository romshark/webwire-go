package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestCustomSessKeyGenInvalid tests custom session key generators returning invalid keys
func TestCustomSessKeyGenInvalid(t *testing.T) {
	// Initialize webwire server
	_, addr := setupServer(
		t,
		wwr.ServerOptions{
			SessionsEnabled: true,
			Hooks: wwr.Hooks{
				OnSessionKeyGeneration: func() string {
					// Return invalid session key
					return ""
				},
				OnRequest: func(ctx context.Context) (wwr.Payload, error) {
					defer func() {
						if err := recover(); err == nil {
							t.Errorf("Expected server to panic on invalid session key")
						}
					}()

					// Extract request message and requesting client from the context
					msg := ctx.Value(wwr.Msg).(wwr.Message)

					// Try to create a new session
					if err := msg.Client.CreateSession(nil); err != nil {
						return wwr.Payload{}, err
					}

					// Return the key of the newly created session (use default binary encoding)
					return wwr.Payload{}, nil
				},
				// Define dummy hooks to enable sessions on this server
				OnSessionCreated: func(_ *wwr.Client) error { return nil },
				OnSessionLookup:  func(_ string) (*wwr.Session, error) { return nil, nil },
				OnSessionClosed:  func(_ *wwr.Client) error { return nil },
			},
		},
	)

	// Initialize client
	client := wwrclt.NewClient(
		addr,
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)
	defer client.Close()

	// Send authentication request and await reply
	if _, err := client.Request("login", wwr.Payload{Data: []byte("testdata")}); err != nil {
		t.Fatalf("Request failed: %s", err)
	}
}
