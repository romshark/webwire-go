package webwire

import (
	"context"

	"github.com/qbeon/webwire-go/message"
)

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSignal(con *connection, msg *message.Message) {
	// Recover potential user-space hook panics to avoid panicking the server
	defer func() {
		if recvErr := recover(); recvErr != nil {
			srv.errorLog.Printf("signal handler failed: %v", recvErr)
		}
		srv.deregisterHandler(con)
	}()

	srv.impl.OnSignal(context.Background(), con, msg)
}
