package webwire

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	msg "github.com/qbeon/webwire-go/message"
)

type connectionStatus = int32

const (
	statActive connectionStatus = iota
	statInactive
)

// ClientInfo represents basic information about a client connection
type ClientInfo struct {
	ConnectionTime time.Time
	UserAgent      string
	RemoteAddr     net.Addr
}

// connection represents a connected client connected to the server
type connection struct {
	statLock    sync.RWMutex
	stat        connectionStatus
	tasks       int32
	srv         *server
	sock        Socket
	sessionLock sync.RWMutex
	session     *Session
	info        ClientInfo
}

// newConnection creates and returns a new client connection instance
func newConnection(
	socket Socket,
	userAgent string,
	srv *server,
) *connection {
	var remoteAddr net.Addr
	stat := statInactive

	if socket != nil {
		stat = statActive
		remoteAddr = socket.RemoteAddr()
	}

	return &connection{
		statLock:    sync.RWMutex{},
		stat:        stat,
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
	con.statLock.RLock()
	defer con.statLock.RUnlock()
	return con.stat == statActive
}

// registerTask increments the number of currently executed tasks
func (con *connection) registerTask() {
	con.statLock.Lock()
	con.tasks++
	con.statLock.Unlock()
}

// deregisterTask decrements the number of currently executed tasks
// and closes this client connection if its shutdown is requested
// and the number of tasks reached zero
func (con *connection) deregisterTask() {
	unlink := false

	con.statLock.Lock()
	con.tasks--
	if con.stat == statInactive && con.tasks < 1 {
		unlink = true
	}
	con.statLock.Unlock()

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

	con.statLock.Lock()
	con.stat = statInactive
	con.statLock.Unlock()

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
		return fmt.Errorf("Couldn't notify client about the session creation: %s", err)
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
	defer con.sessionLock.RUnlock()
	return con.session != nil
}

// Session implements the Connection interface
func (con *connection) Session() *Session {
	con.sessionLock.RLock()
	defer con.sessionLock.RUnlock()
	return con.session.Clone()
}

// SessionKey implements the Connection interface
func (con *connection) SessionKey() string {
	con.sessionLock.RLock()
	defer con.sessionLock.RUnlock()
	if con.session == nil {
		return ""
	}
	return con.session.Key
}

// SessionCreation implements the Connection interface
func (con *connection) SessionCreation() time.Time {
	con.sessionLock.RLock()
	defer con.sessionLock.RUnlock()
	if con.session == nil {
		return time.Time{}
	}
	return con.session.Creation
}

// SessionInfo implements the Connection interface
func (con *connection) SessionInfo(name string) interface{} {
	con.sessionLock.RLock()
	defer con.sessionLock.RUnlock()
	if con.session == nil || con.session.Info == nil {
		return nil
	}
	return con.session.Info.Value(name)
}

// Close implements the Connection interface
func (con *connection) Close() {
	unlink := false

	con.statLock.Lock()
	if con.stat != statActive {
		con.statLock.Unlock()
		return
	}
	con.stat = statInactive
	if con.tasks < 1 {
		unlink = true
	}
	con.statLock.Unlock()

	if unlink {
		con.unlink()
	}
}
