package webwire

import (
	"encoding/json"
	"fmt"
)

// handleSessionRestore handles session restoration (by session key) requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSessionRestore(clt *Client, msg *Message) error {
	if !srv.sessionsEnabled {
		srv.failMsg(clt, msg, SessionsDisabledErr{})
		return nil
	}

	key := string(msg.Payload.Data)

	sessConsNum := srv.sessionRegistry.sessionConnectionsNum(key)
	if sessConsNum >= 0 && srv.sessionRegistry.maxConns > 0 &&
		uint(sessConsNum+1) > srv.sessionRegistry.maxConns {
		srv.failMsg(clt, msg, MaxSessConnsReachedErr{})
		return nil
	}

	//session, err := srv.sessionManager.OnSessionLookup(key)
	exists, creation, info, err := srv.sessionManager.OnSessionLookup(key)
	if err != nil {
		srv.failMsg(clt, msg, nil)
		return fmt.Errorf("CRITICAL: Session search handler failed: %s", err)
	}
	if !exists {
		srv.failMsg(clt, msg, SessNotFoundErr{})
		return nil
	}

	encodedSessionObj := JSONEncodedSession{
		Key:      key,
		Creation: creation,
		Info:     info,
	}

	// JSON encode the session
	encodedSession, err := json.Marshal(&encodedSessionObj)
	if err != nil {
		// TODO: return internal server error and log it
		srv.failMsg(clt, msg, nil)
		return fmt.Errorf(
			"Couldn't encode session object (%v): %s",
			encodedSessionObj,
			err,
		)
	}

	// parse attached session info
	var parsedSessInfo SessionInfo
	if info != nil && srv.sessionInfoParser != nil {
		parsedSessInfo = srv.sessionInfoParser(info)
	}

	clt.setSession(&Session{
		Key:      key,
		Creation: creation,
		Info:     parsedSessInfo,
	})
	if err := srv.sessionRegistry.register(clt); err != nil {
		panic(fmt.Errorf("The number of concurrent session connections was " +
			"unexpectedly exceeded",
		))
	}

	srv.fulfillMsg(clt, msg, Payload{
		Encoding: EncodingUtf8,
		Data:     encodedSession,
	})

	return nil
}
