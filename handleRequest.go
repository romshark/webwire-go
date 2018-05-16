package webwire

import "context"

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleRequest(clt *Client, msg *Message) {
	replyPayload, returnedErr := srv.impl.OnRequest(
		context.Background(),
		clt,
		msg,
	)
	switch returnedErr.(type) {
	case nil:
		srv.fulfillMsg(clt, msg, replyPayload)
	case ReqErr:
		srv.failMsg(clt, msg, returnedErr)
	case *ReqErr:
		srv.failMsg(clt, msg, returnedErr)
	default:
		srv.errorLog.Printf("Internal error during request handling: %s", returnedErr)
		srv.failMsg(clt, msg, returnedErr)
	}
}
