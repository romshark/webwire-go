package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestCustomSessKeyGen tests custom session key generators
func TestCustomSessKeyGen(t *testing.T) {
	expectedSessionKey := "customkey123"

	// Initialize webwire server
	server := setupServer(
		t,
		wwr.ServerOptions{
			SessionsEnabled: true,
			Hooks: wwr.Hooks{
				OnSessionKeyGeneration: func() string {
					return expectedSessionKey
				},
				OnRequest: func(ctx context.Context) (wwr.Payload, error) {
					// Extract request message and requesting client from the context
					msg := ctx.Value(wwr.Msg).(wwr.Message)

					// Try to create a new session
					if err := msg.Client.CreateSession(nil); err != nil {
						return wwr.Payload{}, err
					}

					key := msg.Client.SessionKey()
					if key != expectedSessionKey {
						t.Errorf("Unexpected session key: %s | %s", expectedSessionKey, key)
					}

					// Return the key of the newly created session (use default binary encoding)
					return wwr.Payload{
						Data: []byte(key),
					}, nil
				},
			},
		},
	)

	// Initialize client
	client := wwrclt.NewClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)
	defer client.Close()

	// Send authentication request and await reply
	if _, err := client.Request("login", wwr.Payload{Data: []byte("testdata")}); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	if client.Session().Key != expectedSessionKey {
		t.Errorf("Unexpected session key: %s | %s", expectedSessionKey, client.Session().Key)
	}
}
