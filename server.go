package webwire

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

const protocolVersion = "1.4"

// server represents a headless WebWire server instance,
// where headless means there's no HTTP server that's hosting it
type server struct {
	impl              ServerImplementation
	httpServer        *http.Server
	listener          net.Listener
	sessionManager    SessionManager
	sessionKeyGen     SessionKeyGenerator
	sessionInfoParser SessionInfoParser

	// State
	addr            net.Addr
	options         ServerOptions
	stopping        uint32
	currentOps      int32
	shutdownRdy     chan bool
	connectionsLock *sync.Mutex
	handlerSlots    *semaphore.Weighted
	connections     []*connection
	sessionsEnabled bool
	sessionRegistry *sessionRegistry

	// Internals
	connUpgrader ConnUpgrader
	warnLog      *log.Logger
	errorLog     *log.Logger
}

func (srv *server) shutdownHTTPServer() error {
	if srv.httpServer == nil {
		return nil
	}
	if err := srv.httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("Couldn't properly shutdown HTTP server: %s", err)
	}
	return nil
}

// Run implements the Server interface
func (srv *server) Run() error {
	// Launch HTTP server
	if err := srv.httpServer.Serve(
		tcpKeepAliveListener{srv.listener.(*net.TCPListener)},
	); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP Server failure: %s", err)
	}

	return nil
}

// Addr implements the Server interface
func (srv *server) Addr() net.Addr {
	return srv.addr
}

// Shutdown implements the Server interface
func (srv *server) Shutdown() error {
	srv.setStopping()
	// Don't block if there's no currently processed operations
	if srv.getOps() < 1 {
		return srv.shutdownHTTPServer()
	}
	<-srv.shutdownRdy

	return srv.shutdownHTTPServer()
}

// ActiveSessionsNum implements the Server interface
func (srv *server) ActiveSessionsNum() int {
	return srv.sessionRegistry.activeSessionsNum()
}

// SessionConnectionsNum implements the Server interface
func (srv *server) SessionConnectionsNum(sessionKey string) int {
	return srv.sessionRegistry.sessionConnectionsNum(sessionKey)
}

// SessionConnections implements the Server interface
func (srv *server) SessionConnections(sessionKey string) []Connection {
	connections := srv.sessionRegistry.sessionConnections(sessionKey)
	if connections == nil {
		return nil
	}
	list := make([]Connection, len(connections))
	for i, connection := range connections {
		list[i] = connection
	}
	return list
}

// CloseSession implements the Server interface
func (srv *server) CloseSession(sessionKey string) int {
	connections := srv.sessionRegistry.sessionConnections(sessionKey)
	if connections == nil {
		return -1
	}
	for _, connection := range connections {
		connection.Close()
	}
	return len(connections)
}

func (srv *server) setStopping() {
	atomic.StoreUint32(&srv.stopping, 1)
}

func (srv *server) isStopping() bool {
	return atomic.LoadUint32(&srv.stopping) == 1
}

func (srv *server) incOps() {
	atomic.AddInt32(&srv.currentOps, 1)
}

func (srv *server) decOps() {
	atomic.AddInt32(&srv.currentOps, -1)
}

func (srv *server) getOps() int {
	return int(atomic.LoadInt32(&srv.currentOps))
}
