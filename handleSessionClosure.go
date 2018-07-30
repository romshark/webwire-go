package webwire

import (
	msg "github.com/qbeon/webwire-go/message"
)

// handleSessionClosure handles session destruction requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSessionClosure(
	conn *connection,
	message *msg.Message,
) {
	if !srv.sessionsEnabled {
		srv.failMsg(conn, message, SessionsDisabledErr{})
		return
	}

	if !conn.HasSession() {
		// Send confirmation even though no session was closed
		srv.fulfillMsg(conn, message, 0, nil)
		return
	}

	// Deregister session from active sessions registry
	srv.sessionRegistry.deregister(conn)

	// Synchronize session destruction to the client
	if err := conn.notifySessionClosed(); err != nil {
		srv.failMsg(conn, message, nil)
		srv.errorLog.Printf("CRITICAL: Internal server error, "+
			"couldn't notify client about the session destruction: %s",
			err,
		)
		return
	}

	// Reset the session on the connection
	conn.setSession(nil)

	// Send confirmation
	srv.fulfillMsg(conn, message, 0, nil)
}
