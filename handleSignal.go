package webwire

import (
	"context"

	msg "github.com/qbeon/webwire-go/message"
)

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSignal(con *connection, message *msg.Message) {
	// Ignore incoming signals during shutdown
	if srv.isStopping() {
		return
	}
	srv.incOps()

	srv.impl.OnSignal(
		context.Background(),
		con,
		&MessageWrapper{
			actual: message,
		},
	)

	// Mark signal as done and shutdown the server if scheduled and no ops are left
	srv.decOps()
	if srv.isStopping() && srv.getOps() < 1 {
		close(srv.shutdownRdy)
	}
}
