package memchan

import (
	"errors"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qbeon/webwire-go/connopt"
	"github.com/qbeon/webwire-go/transport"
)

const serverClosed = 0
const serverActive = 1

// Transport implements the Transport
type Transport struct {
	BufferSize        uint32
	ConnectionOptions connopt.ConnectionOptions

	onNewConnection transport.OnNewConnection
	isShuttingdown  transport.IsShuttingDown

	readTimeout time.Duration
	clients     []*Socket
	clientsLock *sync.Mutex
	status      uint32
	shutdown    chan struct{}
}

// Initialize implements the Transport interface
func (srv *Transport) Initialize(
	host string,
	readTimeout time.Duration,
	messageBufferSize uint32,
	isShuttingdown transport.IsShuttingDown,
	onNewConnection transport.OnNewConnection,
) error {
	srv.readTimeout = readTimeout
	srv.isShuttingdown = isShuttingdown
	srv.onNewConnection = onNewConnection
	srv.clients = make([]*Socket, 0, 64)
	srv.clientsLock = &sync.Mutex{}
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
	}
	return nil
}

// Address implements the Transport interface
func (srv *Transport) Address() url.URL {
	return url.URL{
		Scheme: "memchan",
	}
}
