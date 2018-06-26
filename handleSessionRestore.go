package webwire

import (
	"encoding/json"
	"fmt"
)

// handleSessionRestore handles session restoration (by session key) requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSessionRestore(clt *Client, msg *Message) {
	if !srv.sessionsEnabled {
		srv.failMsg(clt, msg, SessionsDisabledErr{})
		return
	}

	key := string(msg.Payload.Data)

	sessConsNum := srv.sessionRegistry.sessionConnectionsNum(key)
	if sessConsNum >= 0 && srv.sessionRegistry.maxConns > 0 &&
		uint(sessConsNum+1) > srv.sessionRegistry.maxConns {
		srv.failMsg(clt, msg, MaxSessConnsReachedErr{})
		return
	}

	// Call session manager lookup hook
	result, err := srv.sessionManager.OnSessionLookup(key)

	// Inspect error if any
	switch err := err.(type) {
	case SessNotFoundErr:
		srv.failMsg(clt, msg, SessNotFoundErr{})
		return
	default:
		srv.failMsg(clt, msg, nil)
		srv.errorLog.Printf("CRITICAL: Session search handler failed: %s", err)
		return
	case nil:
	}

	// JSON encode the session
	encodedSessionObj := JSONEncodedSession{
		Key:        key,
		Creation:   result.Creation,
		LastLookup: result.LastLookup,
		Info:       result.Info,
	}
	encodedSession, err := json.Marshal(&encodedSessionObj)
	if err != nil {
		srv.failMsg(clt, msg, nil)
		srv.errorLog.Printf(
			"Couldn't encode session object (%v): %s",
			encodedSessionObj,
			err,
		)
		return
	}

	// Parse attached session info
	var parsedSessInfo SessionInfo
	if result.Info != nil && srv.sessionInfoParser != nil {
		parsedSessInfo = srv.sessionInfoParser(result.Info)
	}

	clt.setSession(&Session{
		Key:        key,
		Creation:   result.Creation,
		LastLookup: result.LastLookup,
		Info:       parsedSessInfo,
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
}
