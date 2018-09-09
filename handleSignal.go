package webwire

import (
	"context"

	msg "github.com/qbeon/webwire-go/message"
)

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSignal(con *connection, message *msg.Message) {
	srv.opsLock.Lock()
	// Ignore incoming signals during shutdown
	if srv.shutdown {
		srv.opsLock.Unlock()
		return
	}
	srv.currentOps++
	srv.opsLock.Unlock()

	srv.impl.OnSignal(
		context.Background(),
		con,
		NewMessageWrapper(message),
	)

	// Mark signal as done and shutdown the server
	// if scheduled and no ops are left
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()
}
