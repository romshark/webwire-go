package test

import (
	"context"

	wwr "github.com/qbeon/webwire-go"
	"github.com/valyala/fasthttp"
)

// serverImpl implements the webwire.ServerImplementation interface
type serverImpl struct {
	beforeUpgrade        func(ctx *fasthttp.RequestCtx) wwr.ConnectionOptions
	onClientConnected    func(connection wwr.Connection)
	onClientDisconnected func(connection wwr.Connection, reason error)
	onSignal             func(
		ctx context.Context,
		connection wwr.Connection,
		message wwr.Message,
	)
	onRequest func(
		ctx context.Context,
		connection wwr.Connection,
		message wwr.Message,
	) (response wwr.Payload, err error)
}

// OnOptions implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnOptions(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "WEBWIRE")
}

// BeforeUpgrade implements the webwire.ServerImplementation interface
func (srv *serverImpl) BeforeUpgrade(
	ctx *fasthttp.RequestCtx,
) wwr.ConnectionOptions {
	return srv.beforeUpgrade(ctx)
}

// OnClientConnected implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnClientConnected(conn wwr.Connection) {
	srv.onClientConnected(conn)
}

// OnClientDisconnected implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnClientDisconnected(conn wwr.Connection, reason error) {
	srv.onClientDisconnected(conn, reason)
}

// OnSignal implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnSignal(
	ctx context.Context,
	clt wwr.Connection,
	msg wwr.Message,
) {
	srv.onSignal(ctx, clt, msg)
}

// OnRequest implements the webwire.ServerImplementation interface
func (srv *serverImpl) OnRequest(
	ctx context.Context,
	clt wwr.Connection,
	msg wwr.Message,
) (response wwr.Payload, err error) {
	return srv.onRequest(ctx, clt, msg)
}
