package fasthttp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// SockReadErr implements the SockReadErr interface
type SockReadErr struct {
	cause error
}

// Error implements the Go error interface
func (err SockReadErr) Error() string {
	return fmt.Sprintf("Reading socket failed: %s", err.cause)
}

// IsAbnormalCloseErr implements the SockReadErr interface
func (err SockReadErr) IsAbnormalCloseErr() bool {
	return websocket.IsUnexpectedCloseError(
		err.cause,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}

// SockReadWrongMsgTypeErr implements the SockReadErr interface
type SockReadWrongMsgTypeErr struct {
	messageType int
}

// Error implements the Go error interface
func (err SockReadWrongMsgTypeErr) Error() string {
	return fmt.Sprintf("invalid websocket message type: %d", err.messageType)
}

// IsAbnormalCloseErr implements the SockReadErr interface
func (err SockReadWrongMsgTypeErr) IsAbnormalCloseErr() bool {
	return false
}

// Socket implements the webwire.Socket interface using
// the fasthttp/websocket library
type Socket struct {
	connected bool
	lock      sync.RWMutex
	readLock  sync.Mutex
	writeLock sync.Mutex
	conn      *websocket.Conn
	dialer    websocket.Dialer
}

// NewConnectedSocket creates a new fasthttp/websocket based socket
// instance
func NewConnectedSocket(conn *websocket.Conn) *Socket {
	connected := false
	if conn != nil {
		connected = true
	}
	return &Socket{
		connected: connected,
		lock:      sync.RWMutex{},
		readLock:  sync.Mutex{},
		writeLock: sync.Mutex{},
		conn:      conn,
	}
}

// NewSocket creates a new disconnected fasthttp/websocket based socket
// instance
func NewSocket(
	tlsConfig *tls.Config,
	dialTimeout time.Duration,
) *Socket {
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	return &Socket{
		connected: false,
		lock:      sync.RWMutex{},
		readLock:  sync.Mutex{},
		writeLock: sync.Mutex{},
		dialer: websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: dialTimeout,
			TLSClientConfig:  tlsConfig.Clone(),
		},
	}
}

// Dial implements the webwire.Socket interface
func (sock *Socket) Dial(serverAddr url.URL) (err error) {
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
		return wwrerr.DisconnectedErr{
			Cause: fmt.Errorf("dial failure: %s", err),
		}
	}
	sock.conn = connection
	sock.connected = true
	sock.lock.Unlock()
	return nil
}

// fasthttpWriteCloser implements the io.WriteCloser interface
type fasthttpWriteCloser struct {
	writer  io.WriteCloser
	onClose func()
}

// Write implements the io.Writer interface
func (wc *fasthttpWriteCloser) Write(p []byte) (n int, err error) {
	return wc.writer.Write(p)
}

// Close implements the io.Closer interface
func (wc *fasthttpWriteCloser) Close() error {
	err := wc.writer.Close()
	wc.onClose()
	return err
}

// GetWriter implements the webwire.Socket interface
func (sock *Socket) GetWriter() (io.WriteCloser, error) {
	sock.writeLock.Lock()
	if !sock.connected {
		sock.writeLock.Unlock()
		return nil, wwrerr.DisconnectedErr{
			Cause: fmt.Errorf("can't write to a closed socket"),
		}
	}
	writer, err := sock.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		sock.writeLock.Unlock()
		return nil, err
	}
	return &fasthttpWriteCloser{
		writer: writer,
		onClose: func() {
			// Unlock the writer lock on writer closure
			sock.writeLock.Unlock()
		},
	}, nil
}

// Read implements the webwire.Socket interface
func (sock *Socket) Read(msg *message.Message) wwrerr.SockReadErr {
	sock.readLock.Lock()
	messageType, reader, err := sock.conn.NextReader()
	sock.readLock.Unlock()
	if err != nil {
		return SockReadErr{cause: err}
	}

	// Discard message in case of unexpected message types
	if messageType != websocket.BinaryMessage {
		io.Copy(ioutil.Discard, reader)
		return SockReadWrongMsgTypeErr{messageType: messageType}
	}

	// Try to read the socket into the buffer
	typeParsed, err := msg.Read(reader)
	if err != nil {
		return SockReadErr{cause: err}
	}
	if !typeParsed {
		return SockReadErr{cause: errors.New("no message type")}
	}

	return nil
}

// IsConnected implements the webwire.Socket interface
func (sock *Socket) IsConnected() bool {
	sock.lock.RLock()
	connected := sock.connected
	sock.lock.RUnlock()
	return connected
}

// RemoteAddr implements the webwire.Socket interface
func (sock *Socket) RemoteAddr() net.Addr {
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
func (sock *Socket) Close() error {
	sock.lock.Lock()
	sock.connected = false
	err := sock.conn.Close()
	sock.lock.Unlock()
	return err
}

// SetReadDeadline implements the webwire.Socket interface
func (sock *Socket) SetReadDeadline(deadline time.Time) error {
	sock.lock.Lock()
	err := sock.conn.SetReadDeadline(deadline)
	sock.lock.Unlock()
	return err
}

// OnPong implements the webwire.Socket interface
func (sock *Socket) OnPong(handler func(string) error) {
	sock.lock.Lock()
	sock.conn.SetPongHandler(handler)
	sock.lock.Unlock()
}

// OnPing implements the webwire.Socket interface
func (sock *Socket) OnPing(handler func(string) error) {
	sock.lock.Lock()
	sock.conn.SetPingHandler(handler)
	sock.lock.Unlock()
}

// WritePing implements the webwire.Socket interface
func (sock *Socket) WritePing(data []byte, deadline time.Time) error {
	sock.lock.Lock()
	err := sock.conn.WriteControl(websocket.PingMessage, data, deadline)
	sock.lock.Unlock()
	return err
}
