package memchan

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

const statusDisconnected uint32 = 0
const statusConnected uint32 = 1

// SockReadErr implements the SockReadErr interface
type SockReadErr struct {
	// closed is true when the error was caused by a graceful socket closure
	closed bool

	err error
}

// Error implements the Go error interface
func (err SockReadErr) Error() string {
	return fmt.Sprintf("Reading socket failed: %s", err.err)
}

// IsAbnormalCloseErr implements the SockReadErr interface
func (err SockReadErr) IsAbnormalCloseErr() bool {
	return !err.closed
}

// Socket implements the transport.Socket interface using
// the fasthttp/websocket library
type Socket struct {
	// server references the remote server transport for sockets of client type.
	// this reference is nil for sockets of server type
	server *Transport

	remote *Socket

	// connectionStatus indicates the current connection status
	connectionStatus *uint32

	// readTimer enables reader timeout. It's started when the reader is applied
	// a deadline and cleared when it finished reading
	readTimer *time.Timer

	// readLock serializes access to the Read method
	readLock *sync.Mutex

	// writerLock serializes access to the writer returned from GetWriter
	writerLock *sync.Mutex

	// outboundBuffer is used as the primary writer to the remote socket. It's
	// onFlush callback must be connected to the sockets
	outboundBuffer Buffer

	// reader is the data receiving channel to the reader goroutine. It must be
	// triggered as soon as any data is received
	reader chan []byte

	// readerErr must be triggered as soon as the reader finished reading
	readerErr chan error

	// close is a signal causing the local reader to fail and close
	close chan struct{}
}

// onBufferFlush is a slot method that's called by the outbound buffer's onFlush
// callback. It notifies the remote socket about the buffer being available for
// reading and waits until it's finished reading eventually releasing the local
// writer buffer
func (sock *Socket) onBufferFlush(data []byte) error {
	if !sock.IsConnected() {
		sock.writerLock.Unlock()
		return errors.New("can't write to a closed socket")
	}

	// Notify the remote reader about the new message
	sock.remote.reader <- data

	// Wait for the remote reader to finish reading
	readerErr := <-sock.remote.readerErr
	sock.writerLock.Unlock()
	return readerErr
}

// Dial implements the transport.Socket interface
func (sock *Socket) Dial(serverAddr url.URL, deadline time.Time) (err error) {
	if sock.server == nil {
		return errors.New("cannot dial on a server socket")
	}
	return sock.server.Connect(sock)
}

// GetWriter implements the transport.Socket interface
func (sock *Socket) GetWriter() (io.WriteCloser, error) {
	sock.writerLock.Lock()

	// Check connection status
	if !sock.IsConnected() {
		sock.writerLock.Unlock()
		return nil, wwrerr.DisconnectedErr{
			Cause: fmt.Errorf("can't write to a closed socket"),
		}
	}

	// Don't immediately unlock the writer lock, let the writeCloser unlock it
	// as soon as the reader finished reading
	return &sock.outboundBuffer, nil
}

func (sock *Socket) readWithoutDeadline() (
	data []byte,
	err wwrerr.SockReadErr,
) {
	// Await either a message or remote/local socket closure
	select {
	case <-sock.close:
		// Socket closed
		err = SockReadErr{closed: true}
	case data = <-sock.reader:
		// Data received
	}
	return
}

func (sock *Socket) readWithDeadline(deadline time.Time) (
	data []byte,
	err wwrerr.SockReadErr,
) {
	sock.readTimer.Reset(deadline.Sub(time.Now()))

	// Await either a message, socket closure or deadline trigger
	select {
	case <-sock.readTimer.C:
		// Deadline exceeded
		err = SockReadErr{err: errors.New("read deadline exceeded")}
	case <-sock.close:
		// Socket closed
		err = SockReadErr{closed: true}
	case data = <-sock.reader:
		// Data received
	}

	if !sock.readTimer.Stop() {
		<-sock.readTimer.C
	}

	return
}

// Read implements the transport.Socket interface
func (sock *Socket) Read(
	msg *message.Message,
	deadline time.Time,
) (err wwrerr.SockReadErr) {
	// Set reader lock to ensure there's only one concurrent reader
	sock.readLock.Lock()

	// Check connection status
	if !sock.IsConnected() {
		sock.readLock.Unlock()
		return SockReadErr{closed: true}
	}

	var data []byte
	if deadline.IsZero() {
		// No deadline
		data, err = sock.readWithoutDeadline()
	} else {
		// Set deadline
		data, err = sock.readWithDeadline(deadline)
	}

	if err != nil {
		sock.Close()
		sock.readLock.Unlock()
		return err
	}

	// Try to read the socket into the message buffer and release the remote
	// writer by sending nil to the sock.reader channel
	typeParsed, parseErr := msg.ReadBytes(data)
	if parseErr != nil {
		sock.readerErr <- parseErr
		sock.readLock.Unlock()
		return SockReadErr{err: parseErr}
	}
	if !typeParsed {
		err := errors.New("no message type")
		sock.readerErr <- err
		sock.readLock.Unlock()
		return SockReadErr{err: err}
	}

	sock.readerErr <- nil
	sock.readLock.Unlock()
	return nil
}

// IsConnected implements the transport.Socket interface
func (sock *Socket) IsConnected() bool {
	return atomic.LoadUint32(sock.connectionStatus) == statusConnected
}

// RemoteAddr implements the transport.Socket interface
func (sock *Socket) RemoteAddr() net.Addr {
	// The in-memory socket implementation doesn't provide any address
	// information
	return nil
}

// Close implements the transport.Socket interface
func (sock *Socket) Close() error {
	if atomic.CompareAndSwapUint32(
		sock.connectionStatus,
		statusConnected,
		statusDisconnected,
	) {
		// Close the remote reader
		select {
		case sock.remote.close <- struct{}{}:
		default:
		}
		// Close the local reader
		select {
		case sock.close <- struct{}{}:
		default:
		}
	}
	return nil
}
