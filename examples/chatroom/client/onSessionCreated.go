package main

import (
	"log"

	webwire "github.com/qbeon/webwire-go"
)

// onSessionCreated will be invoked when a session is created
func onSessionCreated(session *webwire.Session) {
	username := session.Info["username"].(string)
	log.Printf("Authenticated as %s", username)
}
