package test

import (
	"context"
	"net/http"

	wwr "github.com/qbeon/webwire-go"
)

// serverImpl implements the webwire.ServerImplementation interface
type serverImpl struct {
	beforeUpgrade        func(resp http.ResponseWriter, req *http.Request) bool
	onClientConnected    func(client *wwr.Client)
	onClientDisconnected func(client *wwr.Client)
	onSignal             func(ctx context.Context)
	onRequest            func(ctx context.Context) (response wwr.Payload, err error)
}

// OnOptions implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnOptions(resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "WEBWIRE")
}

// BeforeUpgrade implements the webwire.ServerImplementation interface
func (srv *serverImpl) BeforeUpgrade(resp http.ResponseWriter, req *http.Request) bool {
	return srv.beforeUpgrade(resp, req)
}

// OnClientConnected implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnClientConnected(client *wwr.Client) {
	srv.onClientConnected(client)
}

// OnClientDisconnected implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnClientDisconnected(client *wwr.Client) {
	srv.onClientDisconnected(client)
}

// OnSignal implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnSignal(ctx context.Context) {
	srv.onSignal(ctx)
}

// OnRequest implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnRequest(ctx context.Context) (response wwr.Payload, err error) {
	return srv.onRequest(ctx)
}
