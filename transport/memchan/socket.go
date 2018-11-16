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

	"github.com/qbeon/webwire-go/connopt"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

const statusConnected uint32 = 1
const statusDisconnected uint32 = 2

// SocketType represents the type of a socket
type SocketType uint

const (
	// SocketUninitialized is the default type of an uninitialized socket
	SocketUninitialized SocketType = iota

	// SocketServer represents server-side sockets
	SocketServer

	// SocketClient represents client-side sockets
	SocketClient
)

// Socket implements the transport.Socket interface using
// the fasthttp/websocket library
type Socket struct {
	server *Transport

	sockType SocketType

	// remote references the remote socket
	remote *Socket

	// status represents the connection status
	status *uint32

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
	reader     chan []byte
	readerErr  chan error
	readerLock *sync.Mutex
}

// onBufferFlush is a slot method that's called by the outbound buffer's onFlush
// callback. It notifies the remote socket about the buffer being available for
// reading and waits until it's finished reading eventually releasing the local
// writer buffer
func (sock *Socket) onBufferFlush(data []byte) (err error) {
	if !sock.IsConnected() {
		sock.writerLock.Unlock()
		return errors.New("can't write to a closed socket")
	}

	// Notify the remote reader about the new message
	sock.remote.getReader() <- data

	// Wait for the remote reader to finish reading
	err = <-sock.remote.readerErr

	sock.writerLock.Unlock()
	return
}

// Dial implements the transport.Socket interface
func (sock *Socket) Dial(serverAddr url.URL, deadline time.Time) error {
	if sock.sockType != SocketClient {
		return errors.New("cannot dial on a non-client socket")
	}

	if sock.server.ConnectionOptions.Connection != connopt.Accept {
		return errors.New("connection refused")
	}

	if !atomic.CompareAndSwapUint32(
		sock.status,
		statusDisconnected,
		statusConnected,
	) {
		return errors.New("socket already connected")
	}

	sock.resetReader()
	sock.remote.resetReader()

	// Execute server callback
	if err := sock.server.onConnect(sock); err != nil {
		return err
	}

	return nil
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

func (sock *Socket) getReader() chan []byte {
	sock.readerLock.Lock()
	reader := sock.reader
	sock.readerLock.Unlock()
	return reader
}

// closeReader closes the current reader
func (sock *Socket) closeReader() {
	sock.readerLock.Lock()
	close(sock.reader)
	sock.readerLock.Unlock()
}

// resetReader resets the socket reader
func (sock *Socket) resetReader() {
	sock.readerLock.Lock()
	sock.reader = make(chan []byte, 1)
	sock.readerLock.Unlock()
}

func (sock *Socket) readWithoutDeadline() (
	data []byte,
	err wwrerr.SockReadErr,
) {
	// Await either a message or remote/local socket closure
	data = <-sock.getReader()

	if data == nil {
		// Socket closed
		return nil, SockReadErr{closed: true}
	}

	return data, nil
}

func (sock *Socket) readWithDeadline(deadline time.Time) (
	data []byte,
	err wwrerr.SockReadErr,
) {
	sock.readTimer.Reset(time.Until(deadline))

	// Await either a message, socket closure or deadline trigger
	select {
	case <-sock.readTimer.C:
		// Deadline exceeded
		err = SockReadErr{err: errors.New("read deadline exceeded")}

	case result := <-sock.getReader():
		if result == nil {
			// Socket closed
			err = SockReadErr{closed: true}
		} else {
			data = result
		}
	}

	sock.readTimer.Stop()

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
		sock.readerErr <- nil
		sock.readLock.Unlock()
		return SockReadErr{err: parseErr}
	}
	if !typeParsed {
		sock.readerErr <- nil
		sock.readLock.Unlock()
		return SockReadErr{err: errors.New("no message type")}
	}

	sock.readerErr <- nil
	sock.readLock.Unlock()
	return nil
}

// IsConnected implements the transport.Socket interface
func (sock *Socket) IsConnected() bool {
	return atomic.LoadUint32(sock.status) == statusConnected
}

// RemoteAddr implements the transport.Socket interface
func (sock *Socket) RemoteAddr() net.Addr {
	switch sock.sockType {
	case SocketServer, SocketClient:
		return RemoteAddress{sock.remote}
	}
	return nil
}

// Close implements the transport.Socket interface
func (sock *Socket) Close() error {
	if !atomic.CompareAndSwapUint32(
		sock.status,
		statusConnected,
		statusDisconnected,
	) {
		return nil
	}

	// Close remote reader and writer
	sock.remote.closeReader()

	// Send reader confirmation if any remote writer is currently awaiting
	// the local reader to finish
	select {
	case sock.readerErr <- errors.New("socket closed"):
	default:
	}

	// Close the local reader
	sock.closeReader()

	// Execute server callback
	if sock.sockType == SocketClient {
		sock.server.onDisconnect(sock.remote)
	} else {
		sock.server.onDisconnect(sock)
	}

	return nil
}
