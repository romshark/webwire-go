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
				clt *wwr.Client,
				_ *wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				if err := clt.CreateSession(nil); err != nil {
					return wwr.Payload{}, err
				}

				key := clt.SessionKey()
				if key != expectedSessionKey {
					t.Errorf("Unexpected session key: %s | %s", expectedSessionKey, key)
				}

				// Return the key of the newly created session (use default binary encoding)
				return wwr.Payload{
					Data: []byte(key),
				}, nil
			},
		},
		wwr.ServerOptions{
			SessionsEnabled: true,
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
		nil, nil, nil, nil,
	)
	defer client.connection.Close()

	// Send authentication request and await reply
	if _, err := client.connection.Request(
		"login",
		wwr.Payload{Data: []byte("testdata")},
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
