package main

import (
	webwire "github.com/qbeon/webwire-go"
	"log"
)

// onSessionCreated will be invoked when a session is created
func onSessionCreated(session *webwire.Session) {
	info := session.Info.(map[string]interface{})
	username := info["username"].(string)
	log.Printf("Authenticated as %s", username)
}
