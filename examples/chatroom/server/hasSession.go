package main

import (
	webwire "github.com/qbeon/webwire-go"
)

// hasSession returns true if the given user already has an ongoing session,
// otherwise returns false
func hasSession(client *webwire.Client) bool {
	for _, clt := range sessions {
		if clt == client {
			return true
		}
	}
	return false
}
