package memchan

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/connopt"
)

const serverClosed = 0
const serverActive = 1

// Transport implements the Transport
type Transport struct {
	ConnectionOptions connopt.ConnectionOptions

	onNewConnection wwr.OnNewConnection
	isShuttingdown  wwr.IsShuttingDown

	bufferSize      uint32
	readTimeout     time.Duration
	connections     map[*Socket]*Socket
	connectionsLock *sync.Mutex
	status          uint32
	shutdown        chan struct{}
}

// Initialize implements the Transport interface
func (srv *Transport) Initialize(
	options wwr.ServerOptions,
	isShuttingdown wwr.IsShuttingDown,
	onNewConnection wwr.OnNewConnection,
) error {
	srv.readTimeout = options.ReadTimeout
	srv.bufferSize = options.MessageBufferSize
	srv.isShuttingdown = isShuttingdown
	srv.onNewConnection = onNewConnection
	srv.connections = make(map[*Socket]*Socket)
	srv.connectionsLock = &sync.Mutex{}
	srv.shutdown = make(chan struct{})
	srv.status = serverActive
	return nil
}

// Serve implements the Transport interface
func (srv *Transport) Serve() error {
	if atomic.LoadUint32(&srv.status) != serverActive {
		return errors.New("server is closed")
	}
	<-srv.shutdown
	return nil
}

// Shutdown implements the Transport interface
func (srv *Transport) Shutdown() error {
	if atomic.CompareAndSwapUint32(&srv.status, serverActive, serverClosed) {
		close(srv.shutdown)
		srv.connectionsLock.Lock()
		conns := make([]*Socket, len(srv.connections))
		index := 0
		for sock := range srv.connections {
			conns[index] = sock
			index++
		}
		srv.connectionsLock.Unlock()

		// Close all connections
		for _, sock := range conns {
			if err := sock.Close(); err != nil {
				srv.connectionsLock.Unlock()
				return fmt.Errorf("couldn't close socket %p: %s", sock, err)
			}
		}
	}
	return nil
}

// Address implements the Transport interface
func (srv *Transport) Address() url.URL {
	return url.URL{
		Scheme: "memchan",
	}
}

// onConnect is called in Socket.Dial by a client-type socket on connection
func (srv *Transport) onConnect(clientSocket *Socket) error {
	// Reject incoming connections during server shutdown
	if srv.isShuttingdown() {
		return errors.New("server is shutting down")
	}

	if atomic.LoadUint32(&srv.status) != serverActive {
		return errors.New("server is closed")
	}

	if clientSocket.remote == nil || clientSocket.status == nil {
		return errors.New("uninitialized socket")
	}

	srv.connectionsLock.Lock()
	srv.connections[clientSocket.remote] = clientSocket
	srv.connectionsLock.Unlock()

	go srv.onNewConnection(
		srv.ConnectionOptions,
		[]byte(fmt.Sprintf("webwire memchan client (%p)", clientSocket)),
		clientSocket.remote,
	)

	return nil
}

// onDisconnect is called in Socket.Close by a server-type socket on closure
func (srv *Transport) onDisconnect(serverSocket *Socket) {
	srv.connectionsLock.Lock()
	delete(srv.connections, serverSocket)
	srv.connectionsLock.Unlock()
}
