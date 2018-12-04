package test

import (
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

type TestClientHooks struct {
	OnSessionCreated func(*wwr.Session)
	OnSessionClosed  func()
	OnDisconnected   func()
	OnSignal         func(wwr.Message)
}

// TestClient implements the wwrclt.Implementation interface
type TestClient struct {
	Connection wwrclt.Client
	Hooks      TestClientHooks
}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *TestClient) OnSessionCreated(newSession *wwr.Session) {
	if clt.Hooks.OnSessionCreated != nil {
		clt.Hooks.OnSessionCreated(newSession)
	}
}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *TestClient) OnSessionClosed() {
	if clt.Hooks.OnSessionClosed != nil {
		clt.Hooks.OnSessionClosed()
	}
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *TestClient) OnDisconnected() {
	if clt.Hooks.OnDisconnected != nil {
		clt.Hooks.OnDisconnected()
	}
}

// OnSignal implements the wwrclt.Implementation interface
func (clt *TestClient) OnSignal(msg wwr.Message) {
	if clt.Hooks.OnSignal != nil {
		clt.Hooks.OnSignal(msg)
	}
}
