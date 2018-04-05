package client

import (
	"encoding/json"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) handleSessionCreated(msgPayload webwire.Payload) {
	// Set new session
	var session webwire.Session

	if err := json.Unmarshal(msgPayload.Data, &session); err != nil {
		clt.errorLog.Printf("Failed unmarshalling session object: %s", err)
		return
	}

	clt.sessionLock.Lock()
	clt.session = &session
	clt.sessionLock.Unlock()
	clt.hooks.OnSessionCreated(&session)
}

func (clt *Client) handleSessionClosed() {
	// Destroy local session
	clt.sessionLock.Lock()
	clt.session = nil
	clt.sessionLock.Unlock()

	clt.hooks.OnSessionClosed()
}

func (clt *Client) handleFailure(
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

func (clt *Client) handleInternalError(reqIdent [8]byte) {
	// Fail request
	clt.requestManager.Fail(reqIdent, webwire.ReqInternalErr{})
}

func (clt *Client) handleReplyShutdown(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.ReqSrvShutdownErr{})
}

func (clt *Client) handleSessionNotFound(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.SessNotFoundErr{})
}

func (clt *Client) handleMaxSessConnsReached(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.MaxSessConnsReachedErr{})
}

func (clt *Client) handleSessionsDisabled(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, webwire.SessionsDisabledErr{})
}

func (clt *Client) handleReply(reqIdent [8]byte, payload webwire.Payload) {
	clt.requestManager.Fulfill(reqIdent, payload)
}

func (clt *Client) handleMessage(message []byte) error {
	if len(message) < 1 {
		return nil
	}

	var msg webwire.Message
	if err := msg.Parse(message); err != nil {
		return err
	}

	switch msg.MessageType() {
	case webwire.MsgReplyBinary:
		clt.handleReply(msg.Identifier(), msg.Payload)
	case webwire.MsgReplyUtf8:
		clt.handleReply(msg.Identifier(), msg.Payload)
	case webwire.MsgReplyUtf16:
		clt.handleReply(msg.Identifier(), msg.Payload)
	case webwire.MsgReplyShutdown:
		clt.handleReplyShutdown(msg.Identifier())
	case webwire.MsgSessionNotFound:
		clt.handleSessionNotFound(msg.Identifier())
	case webwire.MsgMaxSessConnsReached:
		clt.handleMaxSessConnsReached(msg.Identifier())
	case webwire.MsgSessionsDisabled:
		clt.handleSessionsDisabled(msg.Identifier())
	case webwire.MsgErrorReply:
		// The message name contains the error code in case of
		// error reply messages, while the UTF8 encoded error message is
		// contained in the message payload
		clt.handleFailure(
			msg.Identifier(),
			msg.Name,
			string(msg.Payload.Data),
		)
	case webwire.MsgInternalError:
		clt.handleInternalError(msg.Identifier())
	case webwire.MsgSignalBinary:
		clt.hooks.OnServerSignal(msg.Payload)
	case webwire.MsgSignalUtf8:
		clt.hooks.OnServerSignal(msg.Payload)
	case webwire.MsgSignalUtf16:
		clt.hooks.OnServerSignal(msg.Payload)
	case webwire.MsgSessionCreated:
		clt.handleSessionCreated(msg.Payload)
	case webwire.MsgSessionClosed:
		clt.handleSessionClosed()
	default:
		clt.warningLog.Printf(
			"Strange message type received: '%d'\n",
			msg.MessageType(),
		)
	}
	return nil
}
