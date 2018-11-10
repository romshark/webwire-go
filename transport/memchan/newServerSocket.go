package memchan

import (
	"sync"
	"time"
)

// NewServerSocket creates a new in-memory server socket instance and
// connects it to the given remote socket, which must be an initialized but
// disconnected socket instance
func NewServerSocket(
	remoteClientSocket *Socket,
	bufferSize uint32,
) (*Socket, error) {
	// Setup a new inactive timer
	readTimer := time.NewTimer(0)
	<-readTimer.C

	connectionStatus := statusConnected

	serverSocket := &Socket{
		remote:           remoteClientSocket,
		connectionStatus: &connectionStatus,
		readLock:         &sync.Mutex{},
		writerLock:       &sync.Mutex{},
		reader:           make(chan []byte),
		readerErr:        make(chan error),
		close:            make(chan struct{}),
		readTimer:        readTimer,
	}

	// Allocate the outbound buffer
	serverSocket.outboundBuffer = NewBuffer(
		make([]byte, bufferSize),
		// Connect the onFlush callback to the corresponding slot method
		serverSocket.onBufferFlush,
	)

	// Make the sockets share the atomic connection status
	remoteClientSocket.remote = serverSocket
	remoteClientSocket.connectionStatus = &connectionStatus

	return serverSocket, nil
}
