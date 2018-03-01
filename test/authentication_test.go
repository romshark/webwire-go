package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestAuthentication verifies the server is connectable,
// and is able to receives requests and signals, create sessions
// and identify clients during request- and signal handling
func TestAuthentication(t *testing.T) {
	clientSignalReceived := NewPending(1, 1*time.Second, true)
	var createdSession *webwire.Session
	expectedCredentials := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("secret_credentials"),
	}
	expectedConfirmation := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("session_is_correct"),
	}
	currentStep := 1

	// Initialize webwire server
	server := setupServer(
		t,
		webwire.Hooks{
			OnSignal: func(ctx context.Context) {
				defer clientSignalReceived.Done()
				// Extract request message and requesting client from the context
				msg := ctx.Value(webwire.Msg).(webwire.Message)
				compareSessions(t, createdSession, msg.Client.Session)
			},
			OnRequest: func(ctx context.Context) (webwire.Payload, *webwire.Error) {
				// Extract request message and requesting client from the context
				msg := ctx.Value(webwire.Msg).(webwire.Message)

				// If already authenticated then check session
				if currentStep > 1 {
					compareSessions(t, createdSession, msg.Client.Session)
					return expectedConfirmation, nil
				}

				// Try to create a new session
				if err := msg.Client.CreateSession(nil); err != nil {
					return webwire.Payload{}, &webwire.Error{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("Internal server error: %s", err),
					}
				}

				// Authentication step is passed
				currentStep = 2

				// Return the key of the newly created session (use default binary encoding)
				return webwire.Payload{
					Data: []byte(msg.Client.Session.Key),
				}, nil
			},
			// Define dummy hooks to enable sessions on this server
			OnSessionCreated: func(_ *webwire.Client) error { return nil },
			OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
			OnSessionClosed:  func(_ *webwire.Client) error { return nil },
		},
	)
	go server.Run()

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	authReqReply, err := client.Request("login", expectedCredentials)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	tmp := client.Session()
	createdSession = &tmp

	// Verify reply
	comparePayload(
		t,
		"authentication reply",
		webwire.Payload{
			Data: []byte(createdSession.Key),
		},
		authReqReply,
	)

	// Send a test-request to verify the session on the server
	// and await response
	testReqReply, err := client.Request("test", expectedCredentials)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "test reply", expectedConfirmation, testReqReply)

	// Send a test-signal to verify the session on the server
	if err := client.Signal("test", expectedCredentials); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	if err := clientSignalReceived.Wait(); err != nil {
		t.Fatal("Client signal not received")
	}
}
