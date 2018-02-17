package webwire

import (
	"fmt"
	"time"
	"sync"
	"bytes"
	"github.com/gorilla/websocket"
)

// Client represents a client connected to the server
type Client struct {
	srv *Server

	lock *sync.Mutex
	conn *websocket.Conn

	connectionTime time.Time
	Session *Session
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
	return clt.connectionTime
}

// Signal sends a signal to this client
func (clt *Client) Signal(payload []byte) error {
	var msg bytes.Buffer
	msg.WriteRune(MsgTyp_SIGNAL)
	msg.Write(payload)
	if err := clt.write(websocket.TextMessage, msg.Bytes()); err != nil {
		return err
	}
	return nil
}

// CreateSession creates a new session for this client.
// Returns an error if there's already another session active
func (clt *Client) CreateSession(session *Session) error {
	if clt.Session != nil {
		return fmt.Errorf("Another session (%s) on this client is already active", clt.Session.Key)
	}
	if err := clt.srv.registerSession(clt, session); err != nil {
		return fmt.Errorf("Couldn't create session: %s", err)
	}
	clt.Session = session

	// Notify client about the session creation
	var msg bytes.Buffer
	msg.WriteRune(MsgTyp_SESS_CREATED)
	msg.WriteString(session.Key)
	if err := clt.write(websocket.TextMessage, msg.Bytes()); err != nil {
		return fmt.Errorf("Couldn't notify client about the session creation: %s", err)
	}

	return nil
}

// Close session destroys the currently active session for this client.
// Does nothing if there's no active session
func (clt *Client) CloseSession() error {
	if clt.Session == nil {
		return nil
	}

	if err := clt.srv.deregisterSession(clt); err != nil {
		return fmt.Errorf("Couldn't close session: %s", err)
	}
	clt.Session = nil

	// Notify client about the session destruction
	var msg bytes.Buffer
	msg.WriteRune(MsgTyp_SESS_CLOSED)
	if err := clt.write(websocket.TextMessage, msg.Bytes()); err != nil {
		return fmt.Errorf("Couldn't notify client about the session destruction: %s", err)
	}

	return nil
}
