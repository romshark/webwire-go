package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/qbeon/webwire-go"

	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// onSignal is invoked when the client receives a signal from the server
func onSignal(message webwire.Payload) {
	var msg shared.ChatMessage
	if err := json.Unmarshal(message.Data, &msg); err != nil {
		panic(fmt.Errorf("Failed parsing chat message: %s", err))
	}

	log.Printf("%s: %s\n", msg.User, msg.Msg)
}
