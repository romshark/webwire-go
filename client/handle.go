package client

import (
	"encoding/json"
	"fmt"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) handleSessionCreated(msgPayload webwire.Payload) {
	var encoded webwire.JSONEncodedSession
	if err := json.Unmarshal(msgPayload.Data, &encoded); err != nil {
		clt.errorLog.Printf("Failed unmarshalling session object: %s", err)
		return
	}

	// parse attached session info
	var parsedSessInfo webwire.SessionInfo
	if encoded.Info != nil && clt.sessionInfoParser != nil {
		parsedSessInfo = clt.sessionInfoParser(encoded.Info)
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

func (clt *Client) handleSessionClosed() {
	// Destroy local session
	clt.sessionLock.Lock()
	clt.session = nil
	clt.sessionLock.Unlock()

	clt.impl.OnSessionClosed()
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
	typeDetermined, err := msg.Parse(message)
	if !typeDetermined {
		return fmt.Errorf("Couldn't determine message type")
	} else if err != nil {
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
		fallthrough
	case webwire.MsgSignalUtf8:
		fallthrough
	case webwire.MsgSignalUtf16:
		clt.impl.OnSignal(msg.Payload)

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
