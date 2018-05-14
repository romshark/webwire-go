package webwire

import "context"

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleRequest(clt *Client, msg *Message) {
	srv.opsLock.Lock()
	// Reject incoming requests during shutdown, return special shutdown error
	if srv.shutdown {
		srv.opsLock.Unlock()
		srv.failMsgShutdown(clt, msg)
		return
	}
	srv.currentOps++
	srv.opsLock.Unlock()

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

	// Mark request as done and shutdown the server if scheduled and no ops are left
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()
}
