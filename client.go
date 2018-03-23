package webwire

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a client connected to the server
type Client struct {
	srv *Server

	connected bool
	connLock  sync.RWMutex
	conn      *websocket.Conn

	connectionTime time.Time
	userAgent      string

	sessionLock sync.Mutex
	Session     *Session
}

// NewClientAgent creates and returns a new client agent instance
func NewClientAgent(conn *websocket.Conn, userAgent string, srv *Server) *Client {
	return &Client{
		srv,
		true,
		sync.RWMutex{},
		conn,
		time.Now(),
		userAgent,
		sync.Mutex{},
		nil,
	}
}

// UserAgent returns the user agent string associated with this client
func (clt *Client) UserAgent() string {
	return clt.userAgent
}

// write sends the given data to the other side of the socket,
// it also protects the connection from concurrent writes
func (clt *Client) write(data []byte) error {
	clt.connLock.Lock()
	defer clt.connLock.Unlock()
	if !clt.connected {
		return DisconnectedErr{
			cause: fmt.Errorf("Can't write to a disconnected client agent"),
		}
	}
	return clt.conn.WriteMessage(websocket.BinaryMessage, data)
}

// unlink resets the client agent and marks it as disconnected preparing it for garbage collection
func (clt *Client) unlink() {
	clt.connLock.Lock()
	clt.sessionLock.Lock()

	clt.connected = false
	clt.Session = nil
	clt.conn.Close()
	clt.conn = nil

	clt.sessionLock.Unlock()
	clt.connLock.Unlock()
}

// ConnectionTime returns the time when the connection was established
func (clt *Client) ConnectionTime() time.Time {
	clt.connLock.RLock()
	defer clt.connLock.RUnlock()
	return clt.connectionTime
}

// RemoteAddr returns the address of the client.
// Returns empty string if the client is not connected.
func (clt *Client) RemoteAddr() net.Addr {
	clt.connLock.RLock()
	defer clt.connLock.RUnlock()
	if clt.conn == nil {
		return nil
	}
	return clt.conn.RemoteAddr()
}

// IsConnected returns true if the client is currently connected to the server,
// thus able to receive signals, otherwise returns false.
// Disconnected client agents are no longer useful and will be garbage collected.
func (clt *Client) IsConnected() bool {
	clt.connLock.RLock()
	defer clt.connLock.RUnlock()
	return clt.connected
}

// Signal sends a named signal containing the given payload to the client
func (clt *Client) Signal(name string, payload Payload) error {
	return clt.write(NewSignalMessage(name, payload))
}

// CreateSession creates a new session for this client.
// It automatically synchronizes the new session to the remote client.
// The synchronization happens asynchronously using a signal
// and doesn't block the calling goroutine.
// Returns an error if there's already another session active
func (clt *Client) CreateSession(attachment interface{}) error {
	if !clt.srv.sessionsEnabled {
		return SessionsDisabledErr{}
	}

	clt.connLock.RLock()
	if !clt.connected {
		clt.connLock.RUnlock()
		return DisconnectedErr{
			cause: fmt.Errorf("Can't create session on disconnected client agent"),
		}
	}
	clt.connLock.RUnlock()

	clt.sessionLock.Lock()
	defer clt.sessionLock.Unlock()

	// Abort if there's already another active session
	if clt.Session != nil {
		return fmt.Errorf(
			"Another session (%s) on this client is already active",
			clt.Session.Key,
		)
	}

	// Create a new session
	newSession := NewSession(attachment, clt.srv.hooks.OnSessionKeyGeneration)

	// Try to notify about session creation
	if err := clt.notifySessionCreated(&newSession); err != nil {
		return fmt.Errorf("Couldn't notify client about the session creation: %s", err)
	}

	// Register the session
	clt.Session = &newSession
	clt.srv.registerSession(clt)

	return nil
}

func (clt *Client) notifySessionCreated(newSession *Session) error {
	encoded, err := json.Marshal(&newSession)
	if err != nil {
		return fmt.Errorf("Couldn't marshal session object: %s", err)
	}

	// Notify client about the session creation
	msg := make([]byte, 1+len(encoded))
	msg[0] = MsgSessionCreated

	for i := 0; i < len(encoded); i++ {
		msg[1+i] = encoded[i]
	}
	return clt.write(msg)
}

func (clt *Client) notifySessionClosed() error {
	// Notify client about the session destruction
	if err := clt.write([]byte{MsgSessionClosed}); err != nil {
		return fmt.Errorf(
			"Couldn't notify client about the session destruction: %s",
			err,
		)
	}
	return nil
}

// CloseSession destroys the currently active session for this client.
// It automatically synchronizes the session destruction to the client.
// The synchronization happens asynchronously using a signal
// and doesn't block the calling goroutine.
// Does nothing if there's no active session
func (clt *Client) CloseSession() error {
	if !clt.srv.sessionsEnabled {
		return SessionsDisabledErr{}
	}

	clt.sessionLock.Lock()
	defer clt.sessionLock.Unlock()

	if clt.Session == nil {
		return nil
	}

	clt.srv.deregisterSession(clt)
	clt.Session = nil

	return clt.notifySessionClosed()
}
