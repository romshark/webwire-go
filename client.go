package webwire

import (
	"time"
	"sync"
	"bytes"
	"github.com/gorilla/websocket"
)

// Client represents a client connected to the server
type Client struct {
	connectionTime time.Time
	conn *websocket.Conn
	lock *sync.Mutex
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
	msg.WriteRune(SIGNAL)
	msg.Write(payload)
	if err := clt.write(websocket.TextMessage, msg.Bytes()); err != nil {
		return err
	}
	return nil
}
