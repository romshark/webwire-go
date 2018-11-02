package webwire

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/qbeon/webwire-go/msgbuf"
)

// fasthttpSockReadErr represents an implementation of the SockReadErr
// interface using the fasthttp/websocket library
type fasthttpSockReadErr struct {
	cause error
}

// Error implements the Go error interface
func (err fasthttpSockReadErr) Error() string {
	return fmt.Sprintf("Reading socket failed: %s", err.cause)
}

// IsAbnormalCloseErr implements the SockReadErr interface
func (err fasthttpSockReadErr) IsAbnormalCloseErr() bool {
	return websocket.IsUnexpectedCloseError(
		err.cause,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}

// fasthttpSockReadWrongMsgTypeErr represents an implementation of the
// SockReadErr interface
type fasthttpSockReadWrongMsgTypeErr struct {
	messageType int
}

// Error implements the Go error interface
func (err fasthttpSockReadWrongMsgTypeErr) Error() string {
	return fmt.Sprintf("invalid websocket message type: %d", err.messageType)
}

// IsAbnormalCloseErr implements the SockReadErr interface
func (err fasthttpSockReadWrongMsgTypeErr) IsAbnormalCloseErr() bool {
	return false
}

// fasthttpSocket implements the webwire.Socket interface using
// the fasthttp/websocket library
type fasthttpSocket struct {
	connected bool
	lock      sync.RWMutex
	readLock  sync.Mutex
	conn      *websocket.Conn
	dialer    websocket.Dialer
}

// newFasthttpConnectedSocket creates a new fasthttp/websocket based socket
// instance
func newFasthttpConnectedSocket(conn *websocket.Conn) Socket {
	connected := false
	if conn != nil {
		connected = true
	}
	return &fasthttpSocket{
		connected: connected,
		lock:      sync.RWMutex{},
		readLock:  sync.Mutex{},
		conn:      conn,
	}
}

// NewFasthttpSocket creates a new disconnected fasthttp/websocket based socket
// instance
func NewFasthttpSocket(
	tlsConfig *tls.Config,
	dialTimeout time.Duration,
) Socket {
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	return &fasthttpSocket{
		connected: false,
		lock:      sync.RWMutex{},
		dialer: websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: dialTimeout,
			TLSClientConfig:  tlsConfig.Clone(),
		},
	}
}

// Dial implements the webwire.Socket interface
func (sock *fasthttpSocket) Dial(serverAddr url.URL) (err error) {
	sock.lock.Lock()
	if sock.connected {
		sock.conn.Close()
		sock.conn = nil
		sock.connected = false
	}

	connection, _, err := sock.dialer.Dial(serverAddr.String(), nil)
	if err != nil {
		sock.connected = false
		sock.lock.Unlock()
		return NewDisconnectedErr(fmt.Errorf("dial failure: %s", err))
	}
	sock.conn = connection
	sock.connected = true
	sock.lock.Unlock()
	return nil
}

// Write implements the webwire.Socket interface
func (sock *fasthttpSocket) Write(data []byte) error {
	sock.lock.Lock()
	if !sock.connected {
		sock.lock.Unlock()
		return DisconnectedErr{
			Cause: fmt.Errorf("can't write to a closed socket"),
		}
	}
	err := sock.conn.WriteMessage(websocket.BinaryMessage, data)
	sock.lock.Unlock()
	return err
}

// Read implements the webwire.Socket interface
func (sock *fasthttpSocket) Read(buf *msgbuf.MessageBuffer) SockReadErr {
	sock.readLock.Lock()
	messageType, reader, err := sock.conn.NextReader()
	sock.readLock.Unlock()
	if err != nil {
		return fasthttpSockReadErr{cause: err}
	}

	// Discard message in case of unexpected message types
	if messageType != websocket.BinaryMessage {
		io.Copy(ioutil.Discard, reader)
		return fasthttpSockReadWrongMsgTypeErr{messageType: messageType}
	}

	// Try to read the socket into the buffer
	if err := buf.Read(reader); err != nil {
		return fasthttpSockReadErr{cause: err}
	}

	return nil
}

// IsConnected implements the webwire.Socket interface
func (sock *fasthttpSocket) IsConnected() bool {
	sock.lock.RLock()
	connected := sock.connected
	sock.lock.RUnlock()
	return connected
}

// RemoteAddr implements the webwire.Socket interface
func (sock *fasthttpSocket) RemoteAddr() net.Addr {
	sock.lock.RLock()
	if sock.conn == nil {
		sock.lock.RUnlock()
		return nil
	}
	addr := sock.conn.RemoteAddr()
	sock.lock.RUnlock()
	return addr
}

// Close implements the webwire.Socket interface
func (sock *fasthttpSocket) Close() error {
	sock.lock.Lock()
	sock.connected = false
	err := sock.conn.Close()
	sock.lock.Unlock()
	return err
}

// SetReadDeadline implements the webwire.Socket interface
func (sock *fasthttpSocket) SetReadDeadline(deadline time.Time) error {
	sock.lock.Lock()
	err := sock.conn.SetReadDeadline(deadline)
	sock.lock.Unlock()
	return err
}

// OnPong implements the webwire.Socket interface
func (sock *fasthttpSocket) OnPong(handler func(string) error) {
	sock.lock.Lock()
	sock.conn.SetPongHandler(handler)
	sock.lock.Unlock()
}

// OnPing implements the webwire.Socket interface
func (sock *fasthttpSocket) OnPing(handler func(string) error) {
	sock.lock.Lock()
	sock.conn.SetPingHandler(handler)
	sock.lock.Unlock()
}

// WritePing implements the webwire.Socket interface
func (sock *fasthttpSocket) WritePing(data []byte, deadline time.Time) error {
	sock.lock.Lock()
	err := sock.conn.WriteControl(websocket.PingMessage, data, deadline)
	sock.lock.Unlock()
	return err
}
