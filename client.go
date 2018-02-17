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

// write protects the connection socket from concurrent writes
func (clt *Client) write(wsMsgType int, data []byte) error {
	clt.lock.Lock()
	defer clt.lock.Unlock()
	return clt.conn.WriteMessage(wsMsgType, data)
}

// ConnectionTime returns the time when the connection was established
func (clt *Client) ConnectionTime() time.Time {
	return clt.connectionTime
}

// Signal sends a signal to the client
func (clt *Client) Signal(payload []byte) error {
	var msg bytes.Buffer
	msg.WriteRune(MsgTyp_SIGNAL)
	msg.Write(payload)
	if err := clt.write(websocket.TextMessage, msg.Bytes()); err != nil {
		return err
	}
	return nil
}

func (clt *Client) CreateSession(session *Session) error {
	if clt.Session != nil {
		return fmt.Errorf("Another session (%s) on this client is already active", clt.Session.Key)
	}
	if err := clt.srv.registerSession(session); err != nil {
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

func (clt *Client) CloseSession() error {
	return fmt.Errorf("CloseSession is not yet implemented")
}
