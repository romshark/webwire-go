package fasthttp

import (
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/qbeon/webwire-go/transport"
)

// ClientTransport implements the webwire client transport layer with fasthttp
type ClientTransport struct {
	// Upgrader specifies the websocket connection upgrader
	Dialer websocket.Dialer
}

// NewSocket implements the ClientTransport interface
func (cltTrans *ClientTransport) NewSocket(
	dialTimeout time.Duration,
) (transport.Socket, error) {
	// Reset handshake timeout to client-enforced dial timeout
	cltTrans.Dialer.HandshakeTimeout = dialTimeout

	return &Socket{
		connected: false,
		lock:      sync.RWMutex{},
		readLock:  sync.Mutex{},
		writeLock: sync.Mutex{},
		dialer:    cltTrans.Dialer,
	}, nil
}
