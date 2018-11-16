package test

import (
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

type testClientHooks struct {
	OnSessionCreated func(*wwr.Session)
	OnSessionClosed  func()
	OnDisconnected   func()
	OnSignal         func(wwr.Message)
}

// testClient implements the wwrclt.Implementation interface
type testClient struct {
	connection wwrclt.Client
	hooks      testClientHooks
}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *testClient) OnSessionCreated(newSession *wwr.Session) {
	if clt.hooks.OnSessionCreated != nil {
		clt.hooks.OnSessionCreated(newSession)
	}
}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *testClient) OnSessionClosed() {
	if clt.hooks.OnSessionClosed != nil {
		clt.hooks.OnSessionClosed()
	}
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *testClient) OnDisconnected() {
	if clt.hooks.OnDisconnected != nil {
		clt.hooks.OnDisconnected()
	}
}

// OnSignal implements the wwrclt.Implementation interface
func (clt *testClient) OnSignal(msg wwr.Message) {
	if clt.hooks.OnSignal != nil {
		clt.hooks.OnSignal(msg)
	}
}
