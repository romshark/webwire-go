package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/qbeon/webwire-go"

	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// Authenticate tries to login using the password and name from the CLI
func (clt *ChatroomClient) Authenticate(login, password string) {
	encodedCreds, err := json.Marshal(shared.AuthenticationCredentials{
		Name:     login,
		Password: password,
	})
	if err != nil {
		panic(fmt.Errorf("Couldn't marshal credentials: %s", err))
	}

	reply, reqErr := clt.connection.Request(
		context.Background(),
		[]byte("auth"),
		webwire.Payload{
			Encoding: webwire.EncodingBinary,
			Data:     encodedCreds,
		},
	)
	reply.Close()
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
