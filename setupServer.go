package webwire

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// Options represents the options for a headed server setup
type Options struct {
	Addr                  string
	Hooks                 Hooks
	MaxSessionConnections uint
	WarnLog               io.Writer
	ErrorLog              io.Writer
}

// SetDefaults sets default values to undefined options
func (opts *Options) SetDefaults() {
	opts.Hooks.SetDefaults()

	if opts.WarnLog == nil {
		opts.WarnLog = os.Stdout
	}

	if opts.ErrorLog == nil {
		opts.ErrorLog = os.Stderr
	}
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
func SetupServer(opts Options) (
	wwrSrv *Server,
	httpServer *http.Server,
	addr string,
	runFunc func() error,
	err error,
) {
	wwrSrv = NewServer(opts)

	// Initialize HTTP server
	httpServer = &http.Server{
		Addr:    opts.Addr,
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
