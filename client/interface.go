package client

import webwire "github.com/qbeon/webwire-go"

// Implementation defines a webwire client implementation interface
type Implementation interface {
	// OnDisconnected is invoked when the client is disconnected
	// from the server for any reason.
	OnDisconnected()

	// OnSignal is invoked when the client receives a signal
	// from the server
	OnSignal(payload webwire.Payload)

	// OnSessionCreated is invoked when the client was assigned a new session
	OnSessionCreated(*webwire.Session)

	// OnSessionClosed is invoked when the client's session was closed
	// either by the server or the client itself
	OnSessionClosed()
}
