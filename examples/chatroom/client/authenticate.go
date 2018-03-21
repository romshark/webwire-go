package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/qbeon/webwire-go"

	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// authenticate tries to login using the password and name from the CLI
func authenticate() {
	encodedCreds, err := json.Marshal(shared.AuthenticationCredentials{
		Name:     *username,
		Password: *password,
	})
	if err != nil {
		panic(fmt.Errorf("Couldn't marshal credentials: %s", err))
	}

	_, reqErr := client.Request("auth", webwire.Payload{Data: encodedCreds})
	switch err := reqErr.(type) {
	case nil:
		break
	case webwire.ReqErr:
		log.Printf("Authentication failed: %s : %s", err.Code, err.Message)
	case webwire.ReqSrvShutdownErr:
		log.Print("Authentication failed, server is currently being shut down")
	default:
		log.Print("Authentication failed, unknown error")
	}
}
