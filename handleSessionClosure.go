package webwire

import (
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// handleSessionClosure handles session destruction requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSessionClosure(
	con *connection,
	msg *message.Message,
) {
	// Recover potential user-space hook panics to avoid panicking the server
	defer func() {
		if recvErr := recover(); recvErr != nil {
			srv.errorLog.Printf(
				"session closure handler panic: %v",
				recvErr,
			)
			srv.failMsg(con, msg, nil)
		}
		srv.deregisterHandler(con)

		// Release message buffer
		msg.Close()
	}()

	if !srv.sessionsEnabled {
		srv.failMsg(con, msg, wwrerr.SessionsDisabledErr{})
		return
	}

	if !con.HasSession() {
		// Send confirmation even though no session was closed
		srv.fulfillMsg(con, msg, Payload{})
		return
	}

	// Deregister session from active sessions registry
	srv.sessionRegistry.deregister(con)

	// Synchronize session destruction to the client
	if err := con.notifySessionClosed(); err != nil {
		srv.failMsg(con, msg, nil)
		srv.errorLog.Printf("internal server error: "+
			"couldn't notify client about the session destruction: %s",
			err,
		)
		return
	}

	// Reset the session on the connection
	con.setSession(nil)

	// Send confirmation
	srv.fulfillMsg(con, msg, Payload{})
}
