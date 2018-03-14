package webwire

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// SetupOptions represents the options used during the setup of
// a new headed WebWire server instance
type SetupOptions struct {
	ServerAddress string
	ServerOptions ServerOptions
}

// SetDefaults sets default values to undefined options
func (opts *SetupOptions) SetDefaults() {
	opts.ServerOptions.SetDefaults()
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts incomming client connections
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// SetupServer sets up a new headed WebWire server with an HTTP server hosting it
func SetupServer(opts SetupOptions) (
	wwrSrv *Server,
	httpServer *http.Server,
	addr string,
	runFunc func() error,
	err error,
) {
	wwrSrv = NewServer(opts.ServerOptions)

	// Initialize HTTP server
	httpServer = &http.Server{
		Addr:    opts.ServerAddress,
		Handler: wwrSrv,
	}

	// Determine final address
	addr = httpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	// Initialize TCP/IP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, "", nil, fmt.Errorf("Failed setting up TCP/IP listener: %s", err)
	}

	runFunc = func() (err error) {
		// Launch server
		err = httpServer.Serve(
			tcpKeepAliveListener{listener.(*net.TCPListener)},
		)
		if err != nil {
			return fmt.Errorf("HTTP Server failure: %s", err)
		}
		return nil
	}

	addr = listener.Addr().String()

	return wwrSrv, httpServer, addr, runFunc, nil
}
