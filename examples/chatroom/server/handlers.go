package main

import (
	"log"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/chatroom/server/state"
)

func onClientConnected(newClient *wwr.Client) {
	state.State.AddConnected(newClient)
	log.Printf("New client connected: %s | %s", newClient.RemoteAddr(), newClient.UserAgent())
}

func onClientDisconnected(client *wwr.Client) {
	state.State.RemoveConnected(client)
	log.Printf("Client %s disconnected", client.RemoteAddr())
}

func onSessionCreated(client *wwr.Client) error {
	state.State.SaveSession(client)
	return nil
}

func onSessionLookup(key string) (*wwr.Session, error) {
	return state.State.FindSession(key)
}

func onSessionClosed(client *wwr.Client) error {
	state.State.CloseSession(client)
	log.Printf("Client %s closed the session", client.RemoteAddr())
	return nil
}
