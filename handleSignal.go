package webwire

import (
	"context"

	"github.com/qbeon/webwire-go/message"
)

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSignal(con *connection, msg *message.Message) {
	srv.impl.OnSignal(context.Background(), con, msg)

	srv.deregisterHandler(con)

	// Release message buffer
	msg.Close()
}
