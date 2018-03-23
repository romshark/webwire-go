package client

import webwire "github.com/qbeon/webwire-go"

// Hooks represents all callback hook functions
type Hooks struct {
	// OnDisconnected is an optional callback.
	// It's invoked when the client is disconnected from the server for any reason.
	OnDisconnected func()

	// OnServerSignal is an optional callback.
	// It's invoked when the webwire client receives a signal from the server
	OnServerSignal func(payload webwire.Payload)

	// OnSessionCreated is an optional callback.
	// It's invoked when the webwire client receives a new session
	OnSessionCreated func(*webwire.Session)

	// OnSessionClosed is an optional callback.
	// It's invoked when the clients session was closed
	// either by the server or by himself
	OnSessionClosed func()
}

// SetDefaults sets undefined required hooks
func (hooks *Hooks) SetDefaults() {
	if hooks.OnDisconnected == nil {
		hooks.OnDisconnected = func() {}
	}

	if hooks.OnServerSignal == nil {
		hooks.OnServerSignal = func(_ webwire.Payload) {}
	}

	if hooks.OnSessionCreated == nil {
		hooks.OnSessionCreated = func(_ *webwire.Session) {}
	}

	if hooks.OnSessionClosed == nil {
		hooks.OnSessionClosed = func() {}
	}
}
