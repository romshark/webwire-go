package webwire

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// connUpgrader implements the webwire.ConnUpgrader interface using
// the gorilla/websocket library
type connUpgrader struct {
	gorillaWsUpgrader websocket.Upgrader
}

// newConnUpgrader constructs a new default HTTP connection upgrader
// based on gorilla/websocket
func newConnUpgrader() *connUpgrader {
	return &connUpgrader{
		gorillaWsUpgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}
}

// Upgrade implements the webwire.ConnUpgrader interface
func (upgrader *connUpgrader) Upgrade(
	resp http.ResponseWriter,
	req *http.Request,
) (Socket, error) {
	conn, err := upgrader.gorillaWsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return nil, err
	}
	return newConnectedSocket(conn), nil
}

// sockReadErr implements the webwire.SockReadErr interface using
// the gorilla/websocket library
type sockReadErr struct {
	cause error
}

// Error implements the Go error interface
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

// socket implements the webwire.Socket interface using
// the gorilla/websocket library
type socket struct {
	connected bool
	lock      sync.RWMutex
	conn      *websocket.Conn
	dialer    websocket.Dialer
}

// newConnectedSocket creates a new gorilla/websocket based socket instance
func newConnectedSocket(conn *websocket.Conn) Socket {
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

// NewSocket creates a new disconnected gorilla/websocket based socket instance
func NewSocket(tlsConfig *tls.Config) Socket {
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	return &socket{
		connected: false,
		lock:      sync.RWMutex{},
		dialer: websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			TLSClientConfig:  tlsConfig.Clone(),
		},
	}
}

// Dial implements the webwire.Socket interface
func (sock *socket) Dial(serverAddr url.URL) (err error) {
	sock.lock.Lock()
	defer sock.lock.Unlock()
	if sock.connected {
		sock.conn.Close()
		sock.conn = nil
	}

	sock.conn, _, err = sock.dialer.Dial(serverAddr.String(), nil)
	if err != nil {
		return NewDisconnectedErr(fmt.Errorf("Dial failure: %s", err))
	}
	sock.connected = true
	return nil
}

// Write implements the webwire.Socket interface
func (sock *socket) Write(data []byte) error {
	sock.lock.Lock()
	defer sock.lock.Unlock()
	if !sock.connected {
		return DisconnectedErr{
			Cause: fmt.Errorf("Can't write to a socket"),
		}
	}
	return sock.conn.WriteMessage(websocket.BinaryMessage, data)
}

// Read implements the webwire.Socket interface
func (sock *socket) Read() ([]byte, SockReadErr) {
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

// SetReadDeadline implements the webwire.Socket interface
func (sock *socket) SetReadDeadline(deadline time.Time) error {
	sock.lock.Lock()
	err := sock.conn.SetReadDeadline(deadline)
	sock.lock.Unlock()
	return err
}

// OnPong implements the webwire.Socket interface
func (sock *socket) OnPong(handler func(string) error) {
	sock.lock.Lock()
	sock.conn.SetPongHandler(handler)
	sock.lock.Unlock()
}

// OnPing implements the webwire.Socket interface
func (sock *socket) OnPing(handler func(string) error) {
	sock.lock.Lock()
	sock.conn.SetPingHandler(handler)
	sock.lock.Unlock()
}

// WritePing implements the webwire.Socket interface
func (sock *socket) WritePing(data []byte, deadline time.Time) error {
	sock.lock.Lock()
	err := sock.conn.WriteControl(websocket.PingMessage, data, deadline)
	sock.lock.Unlock()
	return err
}
