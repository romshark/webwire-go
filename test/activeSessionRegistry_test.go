package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestActiveSessionRegistry verifies that the session registry
// of currently active sessions is properly updated
func TestActiveSessionRegistry(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *webwire.Client,
				msg *webwire.Message,
			) (webwire.Payload, error) {
				// Close session on logout
				if msg.Name == "logout" {
					if err := clt.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
					return webwire.Payload{}, nil
				}

				// Try to create a new session
				if err := clt.CreateSession(nil); err != nil {
					return webwire.Payload{}, err
				}

				// Return the key of the newly created session
				// (use default binary encoding)
				return webwire.Payload{
					Data: []byte(clt.SessionKey()),
				}, nil
			},
		},
		webwire.ServerOptions{
			SessionsEnabled: true,
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: time.Second * 2,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request
	if _, err := client.connection.Request(
		"login",
		webwire.Payload{
			Encoding: webwire.EncodingUtf8,
			Data:     []byte("nothing"),
		},
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberBefore := server.SessionRegistry().ActiveSessions()
	if activeSessionNumberBefore != 1 {
		t.Fatalf(
			"Unexpected active session number after authentication: %d",
			activeSessionNumberBefore,
		)
	}

	// Send logout request
	if _, err := client.connection.Request(
		"logout",
		webwire.Payload{
			Encoding: webwire.EncodingUtf8,
			Data:     []byte("nothing"),
		},
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberAfter := server.SessionRegistry().ActiveSessions()
	if activeSessionNumberAfter != 0 {
		t.Fatalf(
			"Unexpected active session number after logout: %d",
			activeSessionNumberAfter,
		)
	}
}
