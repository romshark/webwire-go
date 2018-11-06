package webwire

import (
	"encoding/json"
	"fmt"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// handleSessionRestore handles session restoration (by session key) requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *server) handleSessionRestore(
	con *connection,
	msg *message.Message,
) {
	if !srv.sessionsEnabled {
		srv.failMsg(con, msg, wwrerr.SessionsDisabledErr{})
		return
	}

	key := string(msg.MsgPayload.Data)

	sessConsNum := srv.sessionRegistry.sessionConnectionsNum(key)
	if sessConsNum >= 0 && srv.sessionRegistry.maxConns > 0 &&
		uint(sessConsNum+1) > srv.sessionRegistry.maxConns {
		srv.failMsg(con, msg, wwrerr.MaxSessConnsReachedErr{})
		return
	}

	// Call session manager lookup hook
	result, err := srv.sessionManager.OnSessionLookup(key)

	if err != nil {
		// Fail message with internal error and log it in case the handler fails
		srv.failMsg(con, msg, nil)
		srv.errorLog.Printf("CRITICAL: Session search handler failed: %s", err)
		return
	}

	if result == nil {
		// Fail message with special error if the session wasn't found
		srv.failMsg(con, msg, wwrerr.SessionNotFoundErr{})
		return
	}

	sessionCreation := result.Creation()
	sessionLastLookup := result.LastLookup()
	sessionInfo := result.Info()

	// JSON encode the session
	encodedSessionObj := JSONEncodedSession{
		Key:        key,
		Creation:   sessionCreation,
		LastLookup: sessionLastLookup,
		Info:       sessionInfo,
	}
	encodedSession, err := json.Marshal(&encodedSessionObj)
	if err != nil {
		srv.failMsg(con, msg, nil)
		srv.errorLog.Printf(
			"Couldn't encode session object (%v): %s",
			encodedSessionObj,
			err,
		)
		return
	}

	// Parse attached session info
	var parsedSessInfo SessionInfo
	if sessionInfo != nil && srv.sessionInfoParser != nil {
		parsedSessInfo = srv.sessionInfoParser(sessionInfo)
	}

	con.setSession(&Session{
		Key:        key,
		Creation:   sessionCreation,
		LastLookup: sessionLastLookup,
		Info:       parsedSessInfo,
	})
	if err := srv.sessionRegistry.register(con); err != nil {
		panic(fmt.Errorf("The number of concurrent session connections was " +
			"unexpectedly exceeded",
		))
	}

	srv.fulfillMsg(
		con,
		msg,
		Payload{
			Encoding: EncodingUtf8,
			Data:     encodedSession,
		},
	)
}
