package test

import (
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// callbackPoweredClient implements the wwrclt.Implementation interface
type callbackPoweredClient struct {
	connection       *wwrclt.Client
	onSessionCreated func(_ *wwr.Session)
	onSessionClosed  func()
	onDisconnected   func()
	onSignal         func(_ wwr.Payload)
}

// newCallbackPoweredClient constructs and returns a new echo client instance
func newCallbackPoweredClient(
	serverAddr string,
	opts wwrclt.Options,
	onSessionCreated func(_ *wwr.Session),
	onSessionClosed func(),
	onDisconnected func(),
	onSignal func(_ wwr.Payload),
) *callbackPoweredClient {
	newClt := &callbackPoweredClient{
		nil,
		onSessionCreated,
		onSessionClosed,
		onDisconnected,
		onSignal,
	}

	// Initialize connection
	newClt.connection = wwrclt.NewClient(serverAddr, newClt, opts)

	return newClt
}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnSessionCreated(newSession *wwr.Session) {
	if clt.onSessionCreated != nil {
		clt.onSessionCreated(newSession)
	}
}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnSessionClosed() {
	if clt.onSessionClosed != nil {
		clt.onSessionClosed()
	}
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnDisconnected() {
	if clt.onDisconnected != nil {
		clt.onDisconnected()
	}
}

// OnSignal implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnSignal(message wwr.Payload) {
	if clt.onSignal != nil {
		clt.onSignal(message)
	}
}
