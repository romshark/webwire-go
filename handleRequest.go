package webwire

import (
	"context"

	"github.com/qbeon/webwire-go/message"
)

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleRequest(conn *connection, msg *message.Message) {
	replyPayload, returnedErr := srv.impl.OnRequest(
		context.Background(),
		conn,
		newMessageWrapper(msg),
	)
	switch returnedErr.(type) {
	case nil:
		// Initialize payload encoding & data
		var encoding PayloadEncoding
		var data []byte
		if replyPayload != nil {
			encoding = replyPayload.Encoding()
			data = replyPayload.Data()
		}

		srv.fulfillMsg(
			conn,
			msg,
			encoding,
			data,
		)
	case ReqErr:
		srv.failMsg(conn, msg, returnedErr)
	case *ReqErr:
		srv.failMsg(conn, msg, returnedErr)
	default:
		srv.errorLog.Printf(
			"Internal error during request handling: %s",
			returnedErr,
		)
		srv.failMsg(conn, msg, returnedErr)
	}
}
