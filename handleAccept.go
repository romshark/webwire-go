package webwire

import (
	"bytes"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var methodNameOptions = []byte("OPTIONS")

func (srv *server) handleAccept(ctx *fasthttp.RequestCtx) {
	// Reject incoming connections during shutdown,
	// pretend the server is temporarily unavailable
	srv.opsLock.Lock()
	if srv.shutdown {
		srv.opsLock.Unlock()
		ctx.Response.Header.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return
	}
	srv.opsLock.Unlock()

	method := ctx.Method()
	if bytes.Equal(method, methodNameOptions) {
		return
	}
	srv.impl.OnOptions(ctx)

	connectionOptions := srv.impl.BeforeUpgrade(ctx)

	// Abort connection establishment if no options are provided
	if connectionOptions.Connection != Accept {
		return
	}

	// Copy the user agent string
	ua := ctx.UserAgent()
	userAgent := make([]byte, len(ua))
	copy(userAgent, ua)

	if err := srv.upgrader.Upgrade(
		ctx,
		func(conn *websocket.Conn) {
			srv.handleConnection(connectionOptions, userAgent, conn)
		},
	); err != nil {
		// Establish connection
		srv.errorLog.Print("upgrade failed:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}
