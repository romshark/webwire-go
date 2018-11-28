package memchan

import (
	"time"

	wwr "github.com/qbeon/webwire-go"
)

// ClientTransport implements the ClientTransport interface
type ClientTransport struct {
	Server *Transport
}

// NewSocket implements the ClientTransport interface
func (cltTrans *ClientTransport) NewSocket(
	dialTimeout time.Duration,
) (wwr.ClientSocket, error) {
	if cltTrans.Server == nil {
		// Create a disconnected socket instance
		return newDisconnectedSocket(), nil
	}

	// Create a new entangled socket pair
	_, clt := NewEntangledSockets(cltTrans.Server)

	return clt, nil
}
