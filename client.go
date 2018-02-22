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
	return clt.connectionTime
}

// RemoteAddr returns the address of the client.
// Returns empty string if the client is not connected.
func (clt *Client) RemoteAddr() net.Addr {
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
// It automatically synchronizes the new session to the client.
// The synchronization happens asynchronously using a signal
// and doesn't block the calling goroutine.
// Returns an error if there's already another session active
func (clt *Client) CreateSession(session *Session) error {
	if clt.Session != nil {
		return fmt.Errorf(
			"Another session (%s) on this client is already active",
			clt.Session.Key,
		)
	}
	if err := clt.srv.registerSession(clt, session); err != nil {
		return fmt.Errorf("Couldn't create session: %s", err)
	}
	clt.Session = session

	// Encode session into JSON
	encoded, err := json.Marshal(*clt.Session)
	if err != nil {
		return fmt.Errorf("Couldn't marshal session object: %s", err)
	}

	// Notify client about the session creation
	var msg bytes.Buffer
	msg.WriteRune(MsgSessionCreated)
	msg.Write(encoded)
	if err := clt.write(websocket.TextMessage, msg.Bytes()); err != nil {
		return fmt.Errorf(
			"Couldn't notify client about the session creation: %s",
			err,
		)
	}

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
	if clt.Session == nil {
		return nil
	}

	if err := clt.srv.deregisterSession(clt); err != nil {
		return fmt.Errorf("Couldn't close session: %s", err)
	}
	clt.Session = nil

	return clt.notifySessionClosed()
}
