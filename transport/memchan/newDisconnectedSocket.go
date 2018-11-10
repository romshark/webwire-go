package memchan

import (
	"sync"
	"time"
)

// NewDisconnectedSocket creates a new disconnected socket instance that needs
// to be connected to a server-side socket before it can be written to and read
// from
func NewDisconnectedSocket(
	server *Transport,
	bufferSize uint32,
) *Socket {
	// Setup a new inactive timer
	readTimer := time.NewTimer(0)
	<-readTimer.C

	sock := &Socket{
		server:     server,
		readLock:   &sync.Mutex{},
		writerLock: &sync.Mutex{},
		reader:     make(chan []byte),
		readerErr:  make(chan error),
		close:      make(chan struct{}),
		readTimer:  readTimer,

		// the following fields will be set during the connection establishment
		connectionStatus: nil,
	}

	// Allocate the outbound buffer
	sock.outboundBuffer = NewBuffer(
		make([]byte, bufferSize),
		// Connect the onFlush callback to the corresponding slot method
		sock.onBufferFlush,
	)

	return sock
}
