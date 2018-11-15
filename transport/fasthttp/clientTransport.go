package fasthttp

import (
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	wwr "github.com/qbeon/webwire-go"
)

// ClientTransport implements the webwire client transport layer with fasthttp
type ClientTransport struct {
	// Upgrader specifies the websocket connection upgrader
	Dialer websocket.Dialer
}

// NewSocket implements the ClientTransport interface
func (cltTrans *ClientTransport) NewSocket(
	dialTimeout time.Duration,
) (wwr.ClientSocket, error) {
	// Reset handshake timeout to client-enforced dial timeout
	cltTrans.Dialer.HandshakeTimeout = dialTimeout

	return &Socket{
		connected: false,
		lock:      &sync.Mutex{},
		readLock:  &sync.Mutex{},
		writeLock: &sync.Mutex{},
		dialer:    cltTrans.Dialer,
	}, nil
}
