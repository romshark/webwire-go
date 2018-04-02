package webwire

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts incoming client connections
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// HeadedServer represents a composition of a webwire server instance
// with a dedicated HTTP server hosting it
type HeadedServer struct {
	addr     net.Addr
	listener net.Listener
	httpSrv  *http.Server
	wwrSrv   *Server
}

// NewHeadedServer sets up a new headed WebWire server
// with an HTTP server hosting it
func NewHeadedServer(
	implementation ServerImplementation,
	opts HeadedServerOptions,
) (*HeadedServer, error) {
	newWwrSrv := NewServer(implementation, opts.ServerOptions)

	// Initialize HTTP server
	httpServer := &http.Server{
		Addr:    opts.ServerAddress,
		Handler: newWwrSrv,
	}

	newHeadedSrv := &HeadedServer{
		httpSrv: httpServer,
		wwrSrv:  newWwrSrv,
	}

	// Determine final address
	addr := httpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	// Initialize TCP/IP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up TCP/IP listener: %s", err)
	}

	newHeadedSrv.listener = listener
	newHeadedSrv.addr = listener.Addr()

	return newHeadedSrv, nil
}

// Addr returns the address of the hosting HTTP server
func (srv *HeadedServer) Addr() net.Addr {
	return srv.addr
}

// SessionRegistry returns the session registry instance of the webwire server
func (srv *HeadedServer) SessionRegistry() SessionRegistry {
	return srv.wwrSrv.sessionRegistry
}

// Run launches both, the webwire and the hosting HTTP server
// blocking the calling goroutine until the server is either gracefully
// shut down or crashes returning an error
func (srv *HeadedServer) Run() error {
	// Launch server
	err := srv.httpSrv.Serve(
		tcpKeepAliveListener{srv.listener.(*net.TCPListener)},
	)
	if err != http.ErrServerClosed {
		return fmt.Errorf("HTTP Server failure: %s", err)
	}
	return nil
}

// Shutdown will block the calling goroutine until both the webwire and
// the hosting HTTP servers are gracefully shut down
func (srv *HeadedServer) Shutdown() error {
	srv.wwrSrv.Shutdown()
	if err := srv.httpSrv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("Couldn't properly shutdown HTTP server: %s", err)
	}
	return nil
}
