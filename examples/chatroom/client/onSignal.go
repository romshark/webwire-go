package main

import (
	"log"
	"fmt"
	"encoding/json"

	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// onSignal is invoked when the client receives a signal from the server
func onSignal(message []byte) {
	var msg shared.ChatMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		panic(fmt.Errorf("Failed parsing chat message: %s", err))
	}

	log.Printf("%s: %s\n", msg.User, msg.Msg)
}
