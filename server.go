package webwire

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

const protocolVersion = "1.3"

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
	shutdown        bool
	shutdownRdy     chan bool
	currentOps      uint32
	opsLock         sync.Mutex
	clientsLock     *sync.Mutex
	clients         []*Client
	sessionsEnabled bool
	sessionRegistry *sessionRegistry

	// Internals
	connUpgrader ConnUpgrader
	warnLog      *log.Logger
	errorLog     *log.Logger
}

func (srv *server) deregisterAgent(clt *Client) {
	// Call the session destruction hook only when the client agent
	// was the last one remaining
	if srv.sessionRegistry.deregister(clt) == 0 {
		if err := srv.sessionManager.OnSessionClosed(clt.session.Key); err != nil {
			srv.errorLog.Printf("OnSessionClosed hook failed: %s", err)
		}
	}
}

func (srv *server) shutdownHttpServer() error {
	if srv.httpServer == nil {
		return nil
	}
	if err := srv.httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("Couldn't properly shutdown HTTP server: %s", err)
	}
	return nil
}

// Run implements the WebwireServer interface
func (srv *server) Run() error {
	// Launch HTTP server
	if err := srv.httpServer.Serve(
		tcpKeepAliveListener{srv.listener.(*net.TCPListener)},
	); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP Server failure: %s", err)
	}

	return nil
}

// Addr implements the WebwireServer interface
func (srv *server) Addr() net.Addr {
	return srv.addr
}

// Shutdown implements the WebwireServer interface
func (srv *server) Shutdown() error {
	srv.opsLock.Lock()
	srv.shutdown = true
	// Don't block if there's no currently processed operations
	if srv.currentOps < 1 {
		return srv.shutdownHttpServer()
	}
	srv.opsLock.Unlock()
	<-srv.shutdownRdy

	return srv.shutdownHttpServer()
}

// ActiveSessionsNum implements the WebwireServer interface
func (srv *server) ActiveSessionsNum() int {
	return srv.sessionRegistry.activeSessionsNum()
}

// SessionConnectionsNum implements the WebwireServer interface
func (srv *server) SessionConnectionsNum(sessionKey string) int {
	return srv.sessionRegistry.sessionConnectionsNum(sessionKey)
}

// SessionConnections implements the WebwireServer interface
func (srv *server) SessionConnections(sessionKey string) []ClientInfo {
	agents := srv.sessionRegistry.sessionConnections(sessionKey)
	if agents == nil {
		return nil
	}
	list := make([]ClientInfo, len(agents))
	for index, clt := range agents {
		list[index] = clt.info
	}
	return list
}

// CloseSession implements the WebwireServer interface
func (srv *server) CloseSession(sessionKey string) int {
	agents := srv.sessionRegistry.sessionConnections(sessionKey)
	if agents == nil {
		return -1
	}
	for _, agent := range agents {
		agent.Close()
	}
	return len(agents)
}
