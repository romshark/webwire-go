package webwire

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

var methodNameOptions = []byte("OPTIONS")
var methodNameWebwire = []byte("WEBWIRE")

func (srv *server) handleHttpRequest(ctx *fasthttp.RequestCtx) {
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
		srv.impl.OnOptions(ctx)
		return
	} else if bytes.Equal(method, methodNameWebwire) {
		srv.handleMetadata(ctx)
		return
	}

	connectionOptions := srv.impl.BeforeUpgrade(ctx)

	// Abort connection establishment if no options are provided
	if connectionOptions.Connection != Accept {
		return
	}

	// Establish connection
	if err := srv.upgrader.Upgrade(ctx, srv.handleConnection); err != nil {
		srv.errorLog.Print("upgrade failed:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}
