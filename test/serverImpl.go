package test

import (
	"context"

	wwr "github.com/qbeon/webwire-go"
)

// ServerImpl implements the webwire.ServerImplementation interface
type ServerImpl struct {
	ClientConnected func(
		connectionOptions wwr.ConnectionOptions,
		connection wwr.Connection,
	)
	ClientDisconnected func(connection wwr.Connection, reason error)
	Signal             func(
		ctx context.Context,
		connection wwr.Connection,
		message wwr.Message,
	)
	Request func(
		ctx context.Context,
		connection wwr.Connection,
		message wwr.Message,
	) (response wwr.Payload, err error)
}

// OnClientConnected implements the webwire.ServerImplementation interface
func (srv *ServerImpl) OnClientConnected(
	opts wwr.ConnectionOptions,
	conn wwr.Connection,
) {
	if srv.ClientConnected != nil {
		srv.ClientConnected(opts, conn)
	}
}

// OnClientDisconnected implements the webwire.ServerImplementation interface
func (srv *ServerImpl) OnClientDisconnected(conn wwr.Connection, reason error) {
	if srv.ClientDisconnected != nil {
		srv.ClientDisconnected(conn, reason)
	}
}

// OnSignal implements the webwire.ServerImplementation interface
func (srv *ServerImpl) OnSignal(
	ctx context.Context,
	clt wwr.Connection,
	msg wwr.Message,
) {
	if srv.Signal != nil {
		srv.Signal(ctx, clt, msg)
	}
}

// OnRequest implements the webwire.ServerImplementation interface
func (srv *ServerImpl) OnRequest(
	ctx context.Context,
	clt wwr.Connection,
	msg wwr.Message,
) (response wwr.Payload, err error) {
	if srv.Request != nil {
		return srv.Request(ctx, clt, msg)
	}
	return wwr.Payload{}, nil
}
