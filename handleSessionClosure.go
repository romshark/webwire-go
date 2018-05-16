package webwire

// handleSessionClosure handles session destruction requests
import "fmt"

// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSessionClosure(clt *Client, msg *Message) error {
	if !srv.sessionsEnabled {
		srv.failMsg(clt, msg, SessionsDisabledErr{})
		return nil
	}

	if !clt.HasSession() {
		// Send confirmation even though no session was closed
		srv.fulfillMsg(clt, msg, Payload{})
		return nil
	}

	// Deregister session from active sessions registry
	srv.sessionRegistry.deregister(clt)

	// Synchronize session destruction to the client
	if err := clt.notifySessionClosed(); err != nil {
		srv.failMsg(clt, msg, nil)
		return fmt.Errorf("CRITICAL: Internal server error, "+
			"couldn't notify client about the session destruction: %s",
			err,
		)
	}

	// Reset the session on the client agent
	clt.setSession(nil)

	// Send confirmation
	srv.fulfillMsg(clt, msg, Payload{})

	return nil
}
