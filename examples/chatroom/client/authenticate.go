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

	if _, reqErr := client.Request("", webwire.Payload{
		Data: encodedCreds,
	}); reqErr != nil {
		log.Printf(
			"Authentication failed: %s : %s",
			reqErr.Code,
			reqErr.Message,
		)
	}
}
