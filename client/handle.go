package client

import (
	"encoding/json"
	"fmt"

	webwire "github.com/qbeon/webwire-go"
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
)

func (clt *client) handleSessionCreated(msgPayload pld.Payload) {
	var encoded webwire.JSONEncodedSession
	if err := json.Unmarshal(msgPayload.Data, &encoded); err != nil {
		clt.errorLog.Printf("Failed unmarshalling session object: %s", err)
		return
	}

	// parse attached session info
	var parsedSessInfo webwire.SessionInfo
	if encoded.Info != nil && clt.options.SessionInfoParser != nil {
		parsedSessInfo = clt.options.SessionInfoParser(encoded.Info)
	}

	clt.sessionLock.Lock()
	clt.session = &webwire.Session{
		Key:      encoded.Key,
		Creation: encoded.Creation,
		Info:     parsedSessInfo,
	}
	clt.sessionLock.Unlock()
	clt.impl.OnSessionCreated(clt.session)
}

func (clt *client) handleSessionClosed() {
	// Destroy local session
	clt.sessionLock.Lock()
	clt.session = nil
	clt.sessionLock.Unlock()

	clt.impl.OnSessionClosed()
}

func (clt *client) handleFailure(
	reqIdent [8]byte,
	errCode,
	errMessage string,
) {
	// Fail request
	clt.requestManager.Fail(reqIdent, webwire.ReqErr{
		Code:    errCode,
		Message: errMessage,
	})
}

func (clt *client) handleInternalError(reqIdent [8]byte) {
	// Fail request
	clt.requestManager.Fail(reqIdent, webwire.ReqInternalErr{})
}

func (clt *client) handleReplyShutdown(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.ReqSrvShutdownErr{})
}

func (clt *client) handleSessionNotFound(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.SessNotFoundErr{})
}

func (clt *client) handleMaxSessConnsReached(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.MaxSessConnsReachedErr{})
}

func (clt *client) handleSessionsDisabled(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.SessionsDisabledErr{})
}

func (clt *client) handleReply(reqIdent [8]byte, payload pld.Payload) {
	clt.requestManager.Fulfill(reqIdent, payload)
}

func (clt *client) handleMessage(message []byte) error {
	if len(message) < 1 {
		return nil
	}

	var parsedMsg msg.Message
	typeDetermined, err := parsedMsg.Parse(message)
	if !typeDetermined {
		return fmt.Errorf("Couldn't determine message type")
	} else if err != nil {
		return err
	}

	switch parsedMsg.Type {
	case msg.MsgReplyBinary:
		clt.handleReply(parsedMsg.Identifier, parsedMsg.Payload)
	case msg.MsgReplyUtf8:
		clt.handleReply(parsedMsg.Identifier, parsedMsg.Payload)
	case msg.MsgReplyUtf16:
		clt.handleReply(parsedMsg.Identifier, parsedMsg.Payload)
	case msg.MsgReplyShutdown:
		clt.handleReplyShutdown(parsedMsg.Identifier)
	case msg.MsgSessionNotFound:
		clt.handleSessionNotFound(parsedMsg.Identifier)
	case msg.MsgMaxSessConnsReached:
		clt.handleMaxSessConnsReached(parsedMsg.Identifier)
	case msg.MsgSessionsDisabled:
		clt.handleSessionsDisabled(parsedMsg.Identifier)
	case msg.MsgErrorReply:
		// The message name contains the error code in case of
		// error reply messages, while the UTF8 encoded error message is
		// contained in the message payload
		clt.handleFailure(
			parsedMsg.Identifier,
			parsedMsg.Name,
			string(parsedMsg.Payload.Data),
		)
	case msg.MsgInternalError:
		clt.handleInternalError(parsedMsg.Identifier)

	case msg.MsgSignalBinary:
		fallthrough
	case msg.MsgSignalUtf8:
		fallthrough
	case msg.MsgSignalUtf16:
		clt.impl.OnSignal(webwire.NewMessageWrapper(&parsedMsg))

	case msg.MsgSessionCreated:
		clt.handleSessionCreated(parsedMsg.Payload)
	case msg.MsgSessionClosed:
		clt.handleSessionClosed()
	default:
		clt.warningLog.Printf(
			"Strange message type received: '%d'\n",
			parsedMsg.Type,
		)
	}
	return nil
}
