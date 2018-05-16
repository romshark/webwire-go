package webwire

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type clientAgentStatus = int32

const (
	statActive clientAgentStatus = iota
	statInactive
)

// ClientInfo represents basic information about a client agent
type ClientInfo struct {
	ConnectionTime time.Time
	UserAgent      string
	RemoteAddr     net.Addr
}

// Client represents a client connected to the server
type Client struct {
	statLock    sync.RWMutex
	stat        clientAgentStatus
	tasks       int32
	srv         *server
	conn        Socket
	sessionLock sync.RWMutex
	session     *Session
	info        ClientInfo
}

// newClientAgent creates and returns a new client agent instance
func newClientAgent(socket Socket, userAgent string, srv *server) *Client {
	var remoteAddr net.Addr
	stat := statInactive

	if socket != nil {
		stat = statActive
		remoteAddr = socket.RemoteAddr()
	}

	return &Client{
		statLock:    sync.RWMutex{},
		stat:        stat,
		tasks:       0,
		srv:         srv,
		conn:        socket,
		sessionLock: sync.RWMutex{},
		session:     nil,
		info: ClientInfo{
			time.Now(),
			userAgent,
			remoteAddr,
		},
	}
}

// isActive must return true if the client agent is in active state
// ready to accept incoming messages
func (clt *Client) isActive() bool {
	clt.statLock.RLock()
	defer clt.statLock.RUnlock()
	return clt.stat == statActive
}

// registerTask increments the number of currently executed tasks
func (clt *Client) registerTask() {
	clt.statLock.Lock()
	clt.tasks++
	clt.statLock.Unlock()
}

// deregisterTask decrements the number of currently executed tasks
// and closes this client agent if its shutdown is requested
// and the number of tasks reached zero
func (clt *Client) deregisterTask() {
	unlink := false

	clt.statLock.Lock()
	clt.tasks--
	if clt.stat == statInactive && clt.tasks < 1 {
		unlink = true
	}
	clt.statLock.Unlock()

	if unlink {
		clt.unlink()
	}
}

// setSession sets a new session for this client
func (clt *Client) setSession(newSess *Session) {
	clt.sessionLock.Lock()
	clt.session = newSess
	clt.sessionLock.Unlock()
}

// unlink resets the client agent and marks it as disconnected
// preparing it for garbage collection
func (clt *Client) unlink() {
	// Deregister session from active sessions registry
	clt.srv.sessionRegistry.deregister(clt)

	clt.sessionLock.Lock()
	clt.session = nil
	clt.sessionLock.Unlock()

	clt.statLock.Lock()
	clt.stat = statInactive
	clt.statLock.Unlock()

	// Close connection
	clt.conn.Close()
}

// Info returns information about the client agent including the
// client agent string, the remote address and the time of creation
func (clt *Client) Info() ClientInfo {
	return clt.info
}

// IsConnected returns true if the client is currently connected to the server,
// thus able to receive signals, otherwise returns false.
// Disconnected client agents are no longer useful and will be garbage collected
func (clt *Client) IsConnected() bool {
	clt.statLock.RLock()
	defer clt.statLock.RUnlock()
	return clt.stat == statActive
}

// Signal sends a named signal containing the given payload to the client
func (clt *Client) Signal(name string, payload Payload) error {
	return clt.conn.Write(NewSignalMessage(name, payload))
}

// CreateSession creates a new session for this client.
// It automatically synchronizes the new session to the remote client.
// The synchronization happens asynchronously using a signal
// and doesn't block the calling goroutine.
// Returns an error if there's already another session active
func (clt *Client) CreateSession(attachment SessionInfo) error {
	if !clt.srv.sessionsEnabled {
		return SessionsDisabledErr{}
	}

	if !clt.conn.IsConnected() {
		return DisconnectedErr{
			Cause: fmt.Errorf("Can't create session on disconnected client agent"),
		}
	}

	clt.sessionLock.Lock()

	// Abort if there's already another active session
	if clt.session != nil {
		clt.sessionLock.Unlock()
		return fmt.Errorf(
			"Another session (%s) on this client is already active",
			clt.session.Key,
		)
	}

	// Create a new session
	newSession := NewSession(attachment, clt.srv.sessionKeyGen.Generate)

	// Try to notify about session creation
	if err := clt.notifySessionCreated(&newSession); err != nil {
		clt.sessionLock.Unlock()
		return fmt.Errorf("Couldn't notify client about the session creation: %s", err)
	}

	// Register the session
	clt.session = &newSession

	clt.srv.sessionRegistry.register(clt)
	clt.sessionLock.Unlock()

	// Call session creation hook
	if err := clt.srv.sessionManager.OnSessionCreated(clt); err != nil {
		clt.srv.errorLog.Printf("OnSessionCreated hook failed: %s", err)
	}

	return nil
}

func (clt *Client) notifySessionCreated(newSession *Session) error {
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
		sessionInfo,
	})
	if err != nil {
		return fmt.Errorf("Couldn't marshal session object: %s", err)
	}

	// Notify client about the session creation
	msg := make([]byte, 1+len(encoded))
	msg[0] = MsgSessionCreated

	for i := 0; i < len(encoded); i++ {
		msg[1+i] = encoded[i]
	}
	return clt.conn.Write(msg)
}

func (clt *Client) notifySessionClosed() error {
	// Notify client about the session destruction
	if err := clt.conn.Write([]byte{MsgSessionClosed}); err != nil {
		return fmt.Errorf(
			"Couldn't notify client about the session destruction: %s",
			err,
		)
	}
	return nil
}

// CloseSession disables the currently active session for this client agent
// and acknowledges the client.
// The session will be destroyed if this is it's last connection remaining.
// Does nothing if there's no active session
func (clt *Client) CloseSession() error {
	if !clt.srv.sessionsEnabled {
		return SessionsDisabledErr{}
	}

	clt.sessionLock.Lock()
	if clt.session == nil {
		clt.sessionLock.Unlock()
		return nil
	}
	// Deregister session from active sessions registry
	clt.srv.sessionRegistry.deregister(clt)
	clt.session = nil
	clt.sessionLock.Unlock()

	return clt.notifySessionClosed()
}

// HasSession returns true if the client referred by this client agent instance
// currently has a session assigned, otherwise returns false
func (clt *Client) HasSession() bool {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	return clt.session != nil
}

// Session returns an exact copy of the session object or nil if there's no
// session currently assigned to this client
func (clt *Client) Session() *Session {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil {
		return nil
	}
	clone := &Session{
		Key:      clt.session.Key,
		Creation: clt.session.Creation,
	}
	if clt.session.Info != nil {
		clone.Info = clt.session.Info.Copy()
	}
	return clone
}

// SessionKey returns the key of the currently assigned session of the client this client agent
// refers to. Returns an empty string if there's no session assigned for this client
func (clt *Client) SessionKey() string {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil {
		return ""
	}
	return clt.session.Key
}

// SessionCreation returns the time of creation of the currently assigned session.
// Warning: be sure to check whether there's a session beforehand as this function will return
// garbage if there's currently no session assigned to the client this user agent refers to
func (clt *Client) SessionCreation() time.Time {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil {
		return time.Time{}
	}
	return clt.session.Creation
}

// SessionInfo returns a copy of the session info field value
// in the form of an empty interface to be casted to either concrete type
func (clt *Client) SessionInfo(name string) interface{} {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil || clt.session.Info == nil {
		return nil
	}
	return clt.session.Info.Value(name)
}

// Close marks the client agent for shutdown. It defers closing the connection
// and removing it from the session registry (if necessary)
// until all work is finished
func (clt *Client) Close() {
	unlink := false

	clt.statLock.Lock()
	if clt.stat != statActive {
		clt.statLock.Unlock()
		return
	}
	clt.stat = statInactive
	if clt.tasks < 1 {
		unlink = true
	}
	clt.statLock.Unlock()

	if unlink {
		clt.unlink()
	}
}
