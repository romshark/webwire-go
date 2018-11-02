package client

import (
	"encoding/json"
	"errors"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/msgbuf"

	webwire "github.com/qbeon/webwire-go"
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

func (clt *client) handleReply(msg *message.Message) {
	clt.requestManager.Fulfill(
		msg.Identifier,
		&webwire.BufferedEncodedPayload{
			Buffer:  msg.Buffer,
			Payload: msg.Payload,
			Closed:  false,
		},
	)
}

// handleMessage handles incoming messages
func (clt *client) handleMessage(buf *msgbuf.MessageBuffer) error {
	msg := &message.Message{
		Buffer: buf,
	}

	// Parse the message from the reader into the buffer
	typeDetermined, err := msg.Parse(buf.Data())
	if !typeDetermined {
		// Release the buffer before returning the error
		buf.Close()
		return errors.New("couldn't determine message type")
	} else if err != nil {
		// Release the buffer before returning the error
		buf.Close()
		return err
	}

	switch msg.Type {
	case message.MsgReplyBinary:
		clt.handleReply(msg)
		// Don't release the buffer, make the user responsible for releasing it
	case message.MsgReplyUtf8:
		clt.handleReply(msg)
		// Don't release the buffer, make the user responsible for releasing it
	case message.MsgReplyUtf16:
		clt.handleReply(msg)
		// Don't release the buffer, make the user responsible for releasing it

	case message.MsgReplyShutdown:
		clt.handleReplyShutdown(msg.Identifier)
		// Release the buffer
		buf.Close()
	case message.MsgSessionNotFound:
		clt.handleSessionNotFound(msg.Identifier)
		// Release the buffer
		buf.Close()
	case message.MsgMaxSessConnsReached:
		clt.handleMaxSessConnsReached(msg.Identifier)
		// Release the buffer
		buf.Close()
	case message.MsgSessionsDisabled:
		clt.handleSessionsDisabled(msg.Identifier)
		// Release the buffer
		buf.Close()
	case message.MsgErrorReply:
		// The message name contains the error code in case of
		// error reply messages, while the UTF8 encoded error message is
		// contained in the message payload
		clt.requestManager.Fail(msg.Identifier, webwire.ReqErr{
			Code:    string(msg.Name),
			Message: string(msg.Payload.Data),
		})
		// Release the buffer
		buf.Close()
	case message.MsgInternalError:
		clt.handleInternalError(msg.Identifier)
		// Release the buffer
		buf.Close()

	case message.MsgSignalBinary:
		fallthrough
	case message.MsgSignalUtf8:
		fallthrough
	case message.MsgSignalUtf16:
		clt.impl.OnSignal(msg.Name, &webwire.BufferedEncodedPayload{
			Buffer:  msg.Buffer,
			Payload: msg.Payload,
			Closed:  false,
		})
		// Realease the buffer af the OnSignal user-space hook is executed
		// because it's referenced there through the payload
		buf.Close()

	case message.MsgSessionCreated:
		clt.handleSessionCreated(msg.Payload)
		// Release the buffer after the OnSessionCreated user-space hook is
		// executed because it's referenced there through the payload
		buf.Close()
	case message.MsgSessionClosed:
		// Release the buffer before calling the OnSessionClosed user-space hook
		buf.Close()
		clt.handleSessionClosed()
	default:
		// Release the buffer
		buf.Close()
		clt.warningLog.Printf(
			"Strange message type received: '%d'\n",
			msg.Type,
		)
	}
	return nil
}
