package client

import (
	"fmt"
	"net"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	webwire "github.com/qbeon/webwire-go"
)

// sockReadErr implements the webwire.SockReadErr interface using the gorilla/websocket library
type sockReadErr struct {
	cause error
}

func (err sockReadErr) Error() string {
	return fmt.Sprintf("Reading socket failed: %s", err.cause)
}

// IsAbnormalCloseErr implements the webwire.SockReadErr interface
func (err sockReadErr) IsAbnormalCloseErr() bool {
	return websocket.IsUnexpectedCloseError(
		err,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}

// socket implements the webwire.Socket interface using the gorilla/websocket library
type socket struct {
	connected bool
	lock      sync.RWMutex
	conn      *websocket.Conn
}

// newSocket creates a new gorilla/websocket based socket instance
func newSocket(conn *websocket.Conn) *socket {
	connected := false
	if conn != nil {
		connected = true
	}
	return &socket{
		connected: connected,
		lock:      sync.RWMutex{},
		conn:      conn,
	}
}

func (sock *socket) Dial(serverAddr string) (err error) {
	connURL := url.URL{Scheme: "ws", Host: serverAddr, Path: "/"}
	sock.lock.Lock()
	defer sock.lock.Unlock()
	if sock.connected {
		sock.conn.Close()
		sock.conn = nil
	}
	sock.conn, _, err = websocket.DefaultDialer.Dial(connURL.String(), nil)
	if err != nil {
		return webwire.NewDisconnectedErr(fmt.Errorf("Dial failure: %s", err))
	}
	sock.connected = true
	return nil
}

// Write implements the webwire.Socket interface
func (sock *socket) Write(data []byte) error {
	sock.lock.Lock()
	defer sock.lock.Unlock()
	if !sock.connected {
		return webwire.DisconnectedErr{
			Cause: fmt.Errorf("Can't write to a socket"),
		}
	}
	return sock.conn.WriteMessage(websocket.BinaryMessage, data)
}

// Read implements the webwire.Socket interface
func (sock *socket) Read() ([]byte, webwire.SockReadErr) {
	_, message, err := sock.conn.ReadMessage()
	if err != nil {
		return nil, sockReadErr{cause: err}
	}
	return message, nil
}

// IsConnected implements the webwire.Socket interface
func (sock *socket) IsConnected() bool {
	sock.lock.RLock()
	defer sock.lock.RUnlock()
	return sock.connected
}

// RemoteAddr implements the webwire.Socket interface
func (sock *socket) RemoteAddr() net.Addr {
	sock.lock.RLock()
	defer sock.lock.RUnlock()
	if sock.conn == nil {
		return nil
	}
	return sock.conn.RemoteAddr()
}

// Close implements the webwire.Socket interface
func (sock *socket) Close() error {
	sock.lock.Lock()
	defer sock.lock.Unlock()
	sock.connected = false
	return sock.conn.Close()
}
