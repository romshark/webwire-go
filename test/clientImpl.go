package test

import (
	"net/url"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

type callbackPoweredClientHooks struct {
	OnSessionCreated func(*wwr.Session)
	OnSessionClosed  func()
	OnDisconnected   func()
	OnSignal         func(wwr.Message)
}

// callbackPoweredClient implements the wwrclt.Implementation interface
type callbackPoweredClient struct {
	connection wwrclt.Client
	hooks      callbackPoweredClientHooks
}

// newCallbackPoweredClient constructs and returns a new echo client instance
func newCallbackPoweredClient(
	serverAddr url.URL,
	options wwrclt.Options,
	hooks callbackPoweredClientHooks,
) *callbackPoweredClient {
	newClt := &callbackPoweredClient{
		nil,
		hooks,
	}

	// Initialize connection
	conn, err := wwrclt.NewClient(serverAddr, newClt, options)
	if err != nil {
		panic(err)
	}

	newClt.connection = conn

	return newClt
}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnSessionCreated(newSession *wwr.Session) {
	if clt.hooks.OnSessionCreated != nil {
		clt.hooks.OnSessionCreated(newSession)
	}
}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnSessionClosed() {
	if clt.hooks.OnSessionClosed != nil {
		clt.hooks.OnSessionClosed()
	}
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnDisconnected() {
	if clt.hooks.OnDisconnected != nil {
		clt.hooks.OnDisconnected()
	}
}

// OnSignal implements the wwrclt.Implementation interface
func (clt *callbackPoweredClient) OnSignal(msg wwr.Message) {
	if clt.hooks.OnSignal != nil {
		clt.hooks.OnSignal(msg)
	}
}
