package test

import (
	"testing"
	"os"
	"fmt"
	"time"
	"sync"
	"context"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
)

// TestSessionCreation verifies the server is connectable,
// and is able to receives requests and create sessions for the requesting client
func TestSessionCreation(t *testing.T) {
	var finish sync.WaitGroup
	var createdSession *webwire.Session
	finish.Add(2)

	// Initialize webwire server
	server := setupServer(
		t,
		nil,
		nil,
		nil,
		// onRequest
		func(ctx context.Context) ([]byte, *webwire.Error) {
			defer finish.Done()

			// Extract request message and requesting client from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// Create a new session
			newSession := webwire.NewSession(
				webwire.OsUnknown,
				"user agent",
				nil,
			)
			createdSession = &newSession

			// Try to register the newly created session and bind it to the client
			if err := msg.Client.CreateSession(createdSession); err != nil {
				return nil, &webwire.Error {
					Code: "INTERNAL_ERROR",
					Message: fmt.Sprintf("Internal server error: %s", err),
				}
			}

			// Return the key of the newly created session
			return []byte(createdSession.Key), nil
		},
		// OnSessionCreated
		func(session *webwire.Session) error {
			// Verify the session
			compareSessions(t, createdSession, session)
			return nil
		},
		// OnSessionLookup
		func(_ string) (*webwire.Session, error) {
			return nil, nil
		},
		// OnSessionClosed
		func(_ *webwire.Client) error {
			return nil
		},
		nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
		// On session creation
		func (newSession *webwire.Session) {
			defer finish.Done()

			// Verify reply
			compareSessions(t, createdSession, newSession)
		},
		nil,
		5 * time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	// Send request and await reply
	reply, err := client.Request([]byte("credentials"))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "server reply", []byte(createdSession.Key), reply)

	// Verify client session
	finish.Wait()
}
