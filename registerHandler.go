package webwire

import (
	"context"

	"github.com/qbeon/webwire-go/message"
)

// registerHandler increments the number of currently executed handlers for this
// particular client and returns true if a handler was registered, otherwise
// returns false. It blocks if the current number of max concurrent handlers was
// reached until a handler slot is available
func (srv *server) registerHandler(
	con *connection,
	msg *message.Message,
) bool {
	failMsg := false

	if !con.IsActive() {
		return false
	}

	// Acquire handler slot if the number of concurrent handlers is limited
	if con.options.ConcurrencyLimit > 1 {
		con.handlerSlots.Acquire(context.Background(), 1)
	}

	srv.opsLock.Lock()
	if srv.shutdown {
		// defer failure due to shutdown of either the server or the connection
		failMsg = true
	} else {
		srv.currentOps++
	}
	srv.opsLock.Unlock()

	if failMsg && msg.RequiresReply() {
		// Don't process the message, fail it
		srv.failMsgShutdown(con, msg)
		return false
	}

	con.registerTask()
	return true
}
