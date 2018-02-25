package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// onRequest handles incoming client requests
// interpreted as authentication attempts. It parses and verifies
// the provided credentials and either rejects the authentication
// or confirms it eventually creating a session and returning the session key
func onRequest(ctx context.Context) ([]byte, *webwire.Error) {
	msg := ctx.Value(webwire.MESSAGE).(webwire.Message)
	client := msg.Client

	log.Printf("Client attempts authentication: %s", client.RemoteAddr())

	// Try to parse credentials
	var credentials shared.AuthenticationCredentials
	if err := json.Unmarshal(msg.Payload, &credentials); err != nil {
		return nil, &webwire.Error{
			Code:    "INTERNAL_ERROR",
			Message: fmt.Sprintf("Failed parsing credentials: %s", err),
		}
	}

	// Verify username
	password, userExists := userAccounts[credentials.Name]
	if !userExists {
		return nil, &webwire.Error{
			Code:    "INEXISTENT_USER",
			Message: fmt.Sprintf("No such user: '%s'", credentials.Name),
		}
	}

	// Verify password
	if password != credentials.Password {
		return nil, &webwire.Error{
			Code:    "WRONG_PASSWORD",
			Message: "Provided password is wrong",
		}
	}

	// Check if client already has an ongoing session
	if hasSession(client) {
		return nil, &webwire.Error{
			Code:    "SESSION_ACTIVE",
			Message: "Already have an ongoing session for this client",
		}
	}

	// Finally create a new session
	if err := client.CreateSession(map[string]string{
		"username": credentials.Name,
	}); err != nil {
		return nil, &webwire.Error{
			Code:    "INTERNAL_ERROR",
			Message: fmt.Sprintf("Couldn't create session: %s", err),
		}
	}

	log.Printf(
		"Created session for user %s (%s)",
		client.RemoteAddr(),
		credentials.Name,
	)

	// Reply to the request
	return []byte(client.Session.Key), nil
}
