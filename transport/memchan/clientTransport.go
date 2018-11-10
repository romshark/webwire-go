package memchan

import (
	"errors"
	"time"

	"github.com/qbeon/webwire-go/transport"
)

// ClientTransport implements the ClientTransport interface
type ClientTransport struct {
	Server     *Transport
	BufferSize uint32
}

// NewSocket implements the ClientTransport interface
func (cltTrans *ClientTransport) NewSocket(
	dialTimeout time.Duration,
) (transport.ClientSocket, error) {
	// Verify server reference
	if cltTrans.Server == nil {
		return nil, errors.New(
			"missing a reference to the memchan server in the client transport",
		)
	}

	// Set default buffer size if none is specified
	bufferSize := uint32(8 * 1024)
	if cltTrans.BufferSize > 0 {
		bufferSize = cltTrans.BufferSize
	}

	// Create a disconnected socket instance
	return NewDisconnectedSocket(cltTrans.Server, bufferSize), nil
}
