package main

import (
	"encoding/json"
	"fmt"
	"log"

	wwr "github.com/qbeon/webwire-go"

	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// broadcastMessage sends a message on behalf of the given user
// to all connected clients
func broadcastMessage(name string, msg string) {
	// Marshal message
	encoded, err := json.Marshal(shared.ChatMessage{
		User: name,
		Msg:  msg,
	})
	if err != nil {
		panic(fmt.Errorf("Couldn't marshal chat message: %s", err))
	}

	// Send message as a signal to each connected client
	log.Printf("Broadcast message to %d clients", len(connectedClients))
	for client := range connectedClients {
		// Send message as signal
		if err := client.Signal("", wwr.Payload{
			Encoding: wwr.EncodingUtf8,
			Data:     encoded,
		}); err != nil {
			log.Printf(
				"WARNING: failed sending signal to client %s : %s",
				client.RemoteAddr(),
				err,
			)
		}
	}
}
