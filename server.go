package webwire

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/qbeon/webwire-go/message"
)

// server represents a headless WebWire server instance,
// where headless means there's no HTTP server that's hosting it
type server struct {
	impl              ServerImplementation
	sessionManager    SessionManager
	sessionKeyGen     SessionKeyGenerator
	sessionInfoParser SessionInfoParser
	addr              url.URL
	options           ServerOptions
	configMsg         []byte
	shutdown          bool
	shutdownRdy       chan bool
	currentOps        uint32
	opsLock           *sync.Mutex
	connectionsLock   *sync.Mutex
	connections       []*connection
	sessionsEnabled   bool
	sessionRegistry   *sessionRegistry
	messagePool       message.Pool

	// Internals
	warnLog  *log.Logger
	errorLog *log.Logger
}

// shutdownServer initiates the shutdown of the underlying transport layer
func (srv *server) shutdownServer() error {
	if err := srv.options.Transport.Shutdown(); err != nil {
		return fmt.Errorf("couldn't properly shutdown HTTP server: %s", err)
	}
	return nil
}

// Run implements the Server interface
func (srv *server) Run() error {
	return srv.options.Transport.Serve()
}

// Address implements the Server interface
func (srv *server) Address() url.URL {
	return srv.options.Transport.Address()
}

// Shutdown implements the Server interface
func (srv *server) Shutdown() error {
	srv.opsLock.Lock()
	srv.shutdown = true
	// Don't block if there's no currently processed operations
	if srv.currentOps < 1 {
		srv.opsLock.Unlock()
		return srv.shutdownServer()
	}
	srv.opsLock.Unlock()

	// Wait until the server is ready for shutdown
	<-srv.shutdownRdy

	return srv.shutdownServer()
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
	i := 0
	for connection := range connections {
		list[i] = connection
		i++
	}
	return list
}

// CloseSession implements the Server interface
func (srv *server) CloseSession(sessionKey string) (
	affectedConnections []Connection,
	errors []error,
	generalError error,
) {
	connections := srv.sessionRegistry.sessionConnections(sessionKey)

	errors = make([]error, len(connections))
	if connections == nil {
		return nil, nil, nil
	}
	affectedConnections = make([]Connection, len(connections))
	i := 0
	errNum := 0
	for connection := range connections {
		affectedConnections[i] = connection
		err := connection.CloseSession()
		if err != nil {
			errors[i] = err
			errNum++
		} else {
			errors[i] = nil
		}
		i++
	}

	if errNum > 0 {
		generalError = fmt.Errorf(
			"%d errors during the closure of a session",
			errNum,
		)
	}

	return affectedConnections, errors, generalError
}
