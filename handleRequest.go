package webwire

import (
	"context"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleRequest(con *connection, msg *message.Message) {
	// Recover potential user-space hook panics to avoid panicking the server
	defer func() {
		if recvErr := recover(); recvErr != nil {
			srv.errorLog.Printf("request handler panic: %v", recvErr)
			srv.failMsg(con, msg, nil)
		}
		srv.deregisterHandler(con)

		// Release message buffer
		msg.Close()
	}()

	// Execute user-space hook
	replyPayload, returnedErr := srv.impl.OnRequest(
		context.Background(),
		con,
		msg,
	)

	// Handle returned error
	switch returnedErr.(type) {
	case nil:
		srv.fulfillMsg(con, msg, replyPayload)
	case wwrerr.RequestErr:
		srv.failMsg(con, msg, returnedErr)
	case *wwrerr.RequestErr:
		srv.failMsg(con, msg, returnedErr)
	default:
		srv.errorLog.Printf(
			"request handler internal error: %v",
			returnedErr,
		)
		srv.failMsg(con, msg, nil)
	}
}
