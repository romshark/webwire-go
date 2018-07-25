package webwire

import (
	msg "github.com/qbeon/webwire-go/message"
)

// handleSessionClosure handles session destruction requests

// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSessionClosure(clt *Client, message *msg.Message) {
	if !srv.sessionsEnabled {
		srv.failMsg(clt, message, SessionsDisabledErr{})
		return
	}

	if !clt.HasSession() {
		// Send confirmation even though no session was closed
		srv.fulfillMsg(clt, message, 0, nil)
		return
	}

	// Deregister session from active sessions registry
	srv.sessionRegistry.deregister(clt)

	// Synchronize session destruction to the client
	if err := clt.notifySessionClosed(); err != nil {
		srv.failMsg(clt, message, nil)
		srv.errorLog.Printf("CRITICAL: Internal server error, "+
			"couldn't notify client about the session destruction: %s",
			err,
		)
		return
	}

	// Reset the session on the client agent
	clt.setSession(nil)

	// Send confirmation
	srv.fulfillMsg(clt, message, 0, nil)
}
