package memchan

import (
	"sync"
	"time"
)

// newSocket creates a new socket instance that must be entangled with another
// socket before it can be used
func newSocket(server *Transport, bufferSize uint32) *Socket {
	// Setup a new inactive timer
	readTimer := time.NewTimer(0)
	<-readTimer.C

	socket := &Socket{
		sockType:   SocketUninitialized,
		server:     server,
		readLock:   &sync.Mutex{},
		writerLock: &sync.Mutex{},
		reader:     make(chan []byte, 1),
		readerLock: &sync.Mutex{},
		readerErr:  make(chan error),
		readTimer:  readTimer,
		remote:     nil,
		status:     nil,
	}

	// Allocate the outbound buffer
	socket.outboundBuffer = NewBuffer(
		make([]byte, bufferSize),
		// Connect the onFlush callback to the corresponding slot method
		socket.onBufferFlush,
	)

	return socket
}
