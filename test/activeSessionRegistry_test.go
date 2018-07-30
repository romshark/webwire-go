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
				conn webwire.Connection,
				msg webwire.Message,
			) (webwire.Payload, error) {
				// Close session on logout
				if msg.Name() == "logout" {
					if err := conn.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
					return nil, nil
				}

				// Try to create a new session
				if err := conn.CreateSession(nil); err != nil {
					return nil, err
				}

				// Return the key of the newly created session
				// (use default binary encoding)
				return webwire.NewPayload(
					webwire.EncodingBinary,
					[]byte(conn.SessionKey()),
				), nil
			},
		},
		webwire.ServerOptions{},
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
		webwire.NewPayload(webwire.EncodingUtf8, []byte("nothing")),
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberBefore := server.ActiveSessionsNum()
	if activeSessionNumberBefore != 1 {
		t.Fatalf(
			"Unexpected active session number after authentication: %d",
			activeSessionNumberBefore,
		)
	}

	// Send logout request
	if _, err := client.connection.Request(
		"logout",
		webwire.NewPayload(webwire.EncodingUtf8, []byte("nothing")),
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberAfter := server.ActiveSessionsNum()
	if activeSessionNumberAfter != 0 {
		t.Fatalf(
			"Unexpected active session number after logout: %d",
			activeSessionNumberAfter,
		)
	}
}
