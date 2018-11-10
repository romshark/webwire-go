package memchan

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// Connect connects a client socket to a new server-side counterpart
func (srv *Transport) Connect(clt *Socket) error {
	// Reject incoming connections during server shutdown
	if srv.isShuttingdown() {
		return errors.New("server is shutting down")
	}

	if atomic.LoadUint32(&srv.status) != serverActive {
		return errors.New("server is closed")
	}

	if clt == nil {
		return errors.New("missing client socket")
	}

	// Setup a new server-side socket and connect it to the client socket
	serverSideSocket, err := NewServerSocket(clt, srv.BufferSize)
	if err != nil {
		return err
	}
	srv.clientsLock.Lock()
	srv.clients = append(srv.clients, serverSideSocket)
	srv.clientsLock.Unlock()

	go srv.onNewConnection(
		srv.ConnectionOptions,
		[]byte(fmt.Sprintf("webwire memchan client (%p)", clt)),
		serverSideSocket,
	)

	return nil
}
