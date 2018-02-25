package webwire

import (
	"bytes"
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

	lock *sync.Mutex
	conn *websocket.Conn

	connectionTime time.Time
	Session        *Session
}

// write sends the given data to the other side of the socket,
// it also protects the connection from concurrent writes
func (clt *Client) write(wsMsgType int, data []byte) error {
	clt.lock.Lock()
	defer clt.lock.Unlock()
	return clt.conn.WriteMessage(wsMsgType, data)
}

// ConnectionTime returns the time when the connection was established
func (clt *Client) ConnectionTime() time.Time {
	clt.lock.Lock()
	defer clt.lock.Unlock()
	return clt.connectionTime
}

// RemoteAddr returns the address of the client.
// Returns empty string if the client is not connected.
func (clt *Client) RemoteAddr() net.Addr {
	clt.lock.Lock()
	defer clt.lock.Unlock()
	if clt.conn == nil {
		return nil
	}
	return clt.conn.RemoteAddr()
}

// Signal sends a signal to the client
func (clt *Client) Signal(payload []byte) error {
	var msg bytes.Buffer
	msg.WriteRune(MsgSignal)
	msg.Write(payload)
	return clt.write(websocket.TextMessage, msg.Bytes())
}

// CreateSession creates a new session for this client.
// It automatically synchronizes the new session to the remote client.
// The synchronization happens asynchronously using a signal
// and doesn't block the calling goroutine.
// Returns an error if there's already another session active
func (clt *Client) CreateSession(attachment interface{}) error {
	if !clt.srv.sessionsEnabled {
		return fmt.Errorf("Sessions disabled")
	}

	clt.lock.Lock()

	// Abort if there's already another active session
	if clt.Session != nil {
		return fmt.Errorf(
			"Another session (%s) on this client is already active",
			clt.Session.Key,
		)
	}

	clt.lock.Unlock()

	// Create a new session
	newSession := NewSession(attachment)

	// Try to serialize the session
	encoded, err := json.Marshal(&newSession)
	if err != nil {
		return fmt.Errorf("Couldn't marshal session object: %s", err)
	}

	// Try to notify the remote client about the session creation
	var msg bytes.Buffer
	msg.WriteRune(MsgSessionCreated)
	msg.Write(encoded)
	if err := clt.write(websocket.TextMessage, msg.Bytes()); err != nil {
		return fmt.Errorf("Couldn't notify client about the session creation: %s", err)
	}

	// Register the session
	clt.lock.Lock()
	clt.Session = &newSession
	clt.srv.registerSession(clt)
	clt.lock.Unlock()

	return nil
}

func (clt *Client) notifySessionClosed() error {
	// Notify client about the session destruction
	if err := clt.write(
		websocket.TextMessage,
		[]byte(string(MsgSessionClosed)),
	); err != nil {
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
		return fmt.Errorf("Sessions disabled")
	}
	if clt.Session == nil {
		return nil
	}

	clt.srv.deregisterSession(clt)
	clt.Session = nil

	return clt.notifySessionClosed()
}
