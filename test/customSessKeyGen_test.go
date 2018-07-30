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
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				if err := conn.CreateSession(nil); err != nil {
					return nil, err
				}

				key := conn.SessionKey()
				if key != expectedSessionKey {
					t.Errorf("Unexpected session key: %s | %s", expectedSessionKey, key)
				}

				// Return the key of the newly created session (use default binary encoding)
				return wwr.NewPayload(
					wwr.EncodingBinary,
					[]byte(key),
				), nil
			},
		},
		wwr.ServerOptions{
			SessionKeyGenerator: &sessionKeyGen{
				generate: func() string {
					return expectedSessionKey
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

	if client.connection.Session().Key != expectedSessionKey {
		t.Errorf(
			"Unexpected session key: %s | %s",
			expectedSessionKey,
			client.connection.Session().Key,
		)
	}
}
