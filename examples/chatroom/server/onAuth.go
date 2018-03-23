package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/chatroom/server/state"
	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// onAuth handles incoming authentication requests.
// It parses and verifies the provided credentials and either rejects the authentication
// or confirms it eventually creating a session and returning the session key
func onAuth(ctx context.Context) (wwr.Payload, error) {
	msg := ctx.Value(wwr.Msg).(wwr.Message)
	client := msg.Client

	credentialsText, err := msg.Payload.Utf8()
	if err != nil {
		return wwr.Payload{}, wwr.ReqErr{
			Code:    "DECODING_FAILURE",
			Message: fmt.Sprintf("Failed decoding message: %s", err),
		}
	}

	log.Printf("Client attempts authentication: %s", client.RemoteAddr())

	// Try to parse credentials
	var credentials shared.AuthenticationCredentials
	if err := json.Unmarshal([]byte(credentialsText), &credentials); err != nil {
		return wwr.Payload{}, fmt.Errorf("Failed parsing credentials: %s", err)
	}

	// Verify username
	password, userExists := userAccounts[credentials.Name]
	if !userExists {
		return wwr.Payload{}, wwr.ReqErr{
			Code:    "INEXISTENT_USER",
			Message: fmt.Sprintf("No such user: '%s'", credentials.Name),
		}
	}

	// Verify password
	if password != credentials.Password {
		return wwr.Payload{}, wwr.ReqErr{
			Code:    "WRONG_PASSWORD",
			Message: "Provided password is wrong",
		}
	}

	// Check if client already has an ongoing session
	if state.State.HasSession(client) {
		return wwr.Payload{}, wwr.ReqErr{
			Code:    "SESSION_ACTIVE",
			Message: "Already have an ongoing session for this client",
		}
	}

	// Finally create a new session
	if err := client.CreateSession(map[string]string{
		"username": credentials.Name,
	}); err != nil {
		return wwr.Payload{}, fmt.Errorf("Couldn't create session: %s", err)
	}

	log.Printf(
		"Created session for user %s (%s)",
		client.RemoteAddr(),
		credentials.Name,
	)

	// Reply to the request, use default binary encoding
	return wwr.Payload{
		Data: []byte(client.Session.Key),
	}, nil
}
