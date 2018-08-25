package webwire

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	msg "github.com/qbeon/webwire-go/message"
)

// ClientInfo represents basic information about a client connection
type ClientInfo struct {
	ConnectionTime time.Time
	UserAgent      string
	RemoteAddr     net.Addr
}

// connection represents a connected client connected to the server
type connection struct {
	// stateLock protects both isActive and tasks from concurrent access
	stateLock sync.RWMutex
	isActive  bool

	// tasks represents the number of currently performed tasks
	tasks int32

	// srv references the connection origin server instance
	srv *server

	// sock references the connection's socket
	sock Socket

	// sessionLock protects the session field from concurrent access
	sessionLock sync.RWMutex

	// session references the currently assigned session, can be null
	session *Session

	// info represents overall connection information
	info ClientInfo
}

// newConnection creates and returns a new client connection instance
func newConnection(
	socket Socket,
	userAgent string,
	srv *server,
) *connection {
	// the connection is considered closed when no socket is referenced
	var remoteAddr net.Addr
	isActive := false

	if socket != nil {
		isActive = true
		remoteAddr = socket.RemoteAddr()
	}

	return &connection{
		stateLock:   sync.RWMutex{},
		isActive:    isActive,
		tasks:       0,
		srv:         srv,
		sock:        socket,
		sessionLock: sync.RWMutex{},
		session:     nil,
		info: ClientInfo{
			time.Now(),
			userAgent,
			remoteAddr,
		},
	}
}

// IsActive implements the Connection interface
func (con *connection) IsActive() bool {
	con.stateLock.RLock()
	isActive := con.isActive
	con.stateLock.RUnlock()
	return isActive

}

// registerTask increments the number of currently executed tasks
func (con *connection) registerTask() {
	con.stateLock.Lock()
	con.tasks++
	con.stateLock.Unlock()
}

// deregisterTask decrements the number of currently executed tasks
// and closes this client connection if its shutdown is requested
// and the number of tasks reached zero
func (con *connection) deregisterTask() {
	unlink := false

	con.stateLock.Lock()
	con.tasks--
	if !con.isActive && con.tasks < 1 {
		unlink = true
	}
	con.stateLock.Unlock()

	if unlink {
		con.unlink()
	}
}

// setSession sets a new session for this client
func (con *connection) setSession(newSess *Session) {
	con.sessionLock.Lock()
	con.session = newSess
	con.sessionLock.Unlock()
}

// unlink resets the connection and marks it as disconnected
// preparing it for garbage collection
func (con *connection) unlink() {
	// Deregister session from active sessions registry
	con.srv.sessionRegistry.deregister(con)

	con.sessionLock.Lock()
	con.session = nil
	con.sessionLock.Unlock()

	// Close connection
	con.sock.Close()
}

// Info implements the Connection interface
func (con *connection) Info() ClientInfo {
	return con.info
}

// Signal implements the Connection interface
func (con *connection) Signal(name string, payload Payload) error {
	return con.sock.Write(msg.NewSignalMessage(
		name,
		payload.Encoding(),
		payload.Data(),
	))
}

// CreateSession implements the Connection interface
func (con *connection) CreateSession(attachment SessionInfo) error {
	if !con.srv.sessionsEnabled {
		return SessionsDisabledErr{}
	}

	if !con.sock.IsConnected() {
		return DisconnectedErr{
			Cause: fmt.Errorf(
				"Can't create session on disconnected connection",
			),
		}
	}

	con.sessionLock.Lock()

	// Abort if there's already another active session
	if con.session != nil {
		con.sessionLock.Unlock()
		return fmt.Errorf(
			"Another session (%s) on this client is already active",
			con.session.Key,
		)
	}

	// Create a new session
	newSession := NewSession(attachment, con.srv.sessionKeyGen.Generate)

	// Try to notify about session creation
	if err := con.notifySessionCreated(&newSession); err != nil {
		con.sessionLock.Unlock()
		return fmt.Errorf(
			"Couldn't notify client about the session creation: %s",
			err,
		)
	}

	// Register the session
	con.session = &newSession

	con.srv.sessionRegistry.register(con)
	con.sessionLock.Unlock()

	// Call session creation hook
	if err := con.srv.sessionManager.OnSessionCreated(con); err != nil {
		con.srv.errorLog.Printf("OnSessionCreated hook failed: %s", err)
	}

	return nil
}

func (con *connection) notifySessionCreated(newSession *Session) error {
	// Serialize session info
	var sessionInfo map[string]interface{}
	if newSession.Info != nil {
		sessionInfo = make(map[string]interface{})
		for _, field := range newSession.Info.Fields() {
			sessionInfo[field] = newSession.Info.Value(field)
		}
	}

	encoded, err := json.Marshal(JSONEncodedSession{
		newSession.Key,
		newSession.Creation,
		newSession.LastLookup,
		sessionInfo,
	})
	if err != nil {
		return fmt.Errorf("Couldn't marshal session object: %s", err)
	}

	// Notify client about the session creation
	message := make([]byte, 1+len(encoded))
	message[0] = msg.MsgSessionCreated

	for i := 0; i < len(encoded); i++ {
		message[1+i] = encoded[i]
	}
	return con.sock.Write(message)
}

func (con *connection) notifySessionClosed() error {
	// Notify client about the session destruction
	if err := con.sock.Write([]byte{msg.MsgSessionClosed}); err != nil {
		return fmt.Errorf(
			"Couldn't notify client about the session destruction: %s",
			err,
		)
	}
	return nil
}

// CloseSession implements the Connection interface
func (con *connection) CloseSession() error {
	if !con.srv.sessionsEnabled {
		return SessionsDisabledErr{}
	}

	con.sessionLock.Lock()
	if con.session == nil {
		con.sessionLock.Unlock()
		return nil
	}
	// Deregister session from active sessions registry
	con.srv.sessionRegistry.deregister(con)
	con.session = nil
	con.sessionLock.Unlock()

	return con.notifySessionClosed()
}

// HasSession implements the Connection interface
func (con *connection) HasSession() bool {
	con.sessionLock.RLock()
	hasSession := con.session != nil
	con.sessionLock.RUnlock()
	return hasSession
}

// Session implements the Connection interface
func (con *connection) Session() *Session {
	con.sessionLock.RLock()
	clone := con.session.Clone()
	con.sessionLock.RUnlock()
	return clone
}

// SessionKey implements the Connection interface
func (con *connection) SessionKey() string {
	con.sessionLock.RLock()
	if con.session == nil {
		con.sessionLock.RUnlock()
		return ""
	}
	key := con.session.Key
	con.sessionLock.RUnlock()
	return key
}

// SessionCreation implements the Connection interface
func (con *connection) SessionCreation() time.Time {
	con.sessionLock.RLock()
	if con.session == nil {
		con.sessionLock.RUnlock()
		return time.Time{}
	}
	creation := con.session.Creation
	con.sessionLock.RUnlock()
	return creation
}

// SessionInfo implements the Connection interface
func (con *connection) SessionInfo(name string) interface{} {
	con.sessionLock.RLock()
	if con.session == nil || con.session.Info == nil {
		con.sessionLock.RUnlock()
		return nil
	}
	val := con.session.Info.Value(name)
	con.sessionLock.RUnlock()
	return val
}

// Close implements the Connection interface
func (con *connection) Close() {
	unlink := false

	con.stateLock.Lock()
	if !con.isActive {
		con.stateLock.Unlock()
		return
	}
	con.isActive = false
	if con.tasks < 1 {
		unlink = true
	}
	con.stateLock.Unlock()

	if unlink {
		con.unlink()
	}
}
