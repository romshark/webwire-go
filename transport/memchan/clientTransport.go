package memchan

import (
	"errors"
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
	// Verify server reference
	if cltTrans.Server == nil {
		return nil, errors.New(
			"missing a reference to the memchan server in the client transport",
		)
	}

	// Create a new entangled socket pair
	_, clt := NewEntangledSockets(cltTrans.Server)

	return clt, nil
}
