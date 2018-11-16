package main

import (
	"encoding/json"
	"fmt"
	"log"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// OnSessionCreated implements the webwireClient.Implementation interface
// it's invoked when a new session is assigned to the client
func (clt *ChatroomClient) OnSessionCreated(newSession *webwire.Session) {
	username := newSession.Info.Value("username").(string)
	log.Printf("Authenticated as %s", username)
}

// OnSignal implements the webwireClient.Implementation interface.
// it's invoked when the client receives a signal from the server
// containing a chatroom message
func (clt *ChatroomClient) OnSignal(msg webwire.Message) {
	var chatMsg shared.ChatMessage

	// Interpret the message as UTF8 encoded JSON
	jsonString, err := msg.PayloadUtf8()
	if err != nil {
		log.Printf("Couldn't decode incoming message: %s\n", err)
	}

	if err := json.Unmarshal([]byte(jsonString), &chatMsg); err != nil {
		panic(fmt.Errorf("Failed parsing chat message: %s", err))
	}

	log.Printf("%s: %s\n", chatMsg.User, chatMsg.Msg)
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *ChatroomClient) OnDisconnected() {}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *ChatroomClient) OnSessionClosed() {}
