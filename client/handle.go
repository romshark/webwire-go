package client

import (
	"encoding/json"
	"fmt"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	"github.com/qbeon/webwire-go/wwrerr"
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
	clt.requestManager.Fail(reqIdent, wwrerr.InternalErr{})
}

func (clt *client) handleReplyShutdown(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, wwrerr.ServerShutdownErr{})
}

func (clt *client) handleSessionNotFound(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, wwrerr.SessionNotFoundErr{})
}

func (clt *client) handleMaxSessConnsReached(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, wwrerr.MaxSessConnsReachedErr{})
}

func (clt *client) handleSessionsDisabled(reqIdent [8]byte) {
	clt.requestManager.Fail(reqIdent, wwrerr.SessionsDisabledErr{})
}

// handleMessage handles incoming messages
func (clt *client) handleMessage(msg *message.Message) (err error) {
	// Recover user-space panics to avoid leaking memory through unreleased
	// message buffer
	defer func() {
		if recvErr := recover(); recvErr != nil {
			r, ok := recvErr.(error)
			if !ok {
				err = fmt.Errorf("unexpected panic: %v", recvErr)
			} else {
				err = r
			}
		}
	}()

	switch msg.MsgType {
	case message.MsgReplyBinary:
		clt.requestManager.Fulfill(msg)
		// Don't release the buffer, make the user responsible for releasing it
	case message.MsgReplyUtf8:
		clt.requestManager.Fulfill(msg)
		// Don't release the buffer, make the user responsible for releasing it
	case message.MsgReplyUtf16:
		clt.requestManager.Fulfill(msg)
		// Don't release the buffer, make the user responsible for releasing it

	case message.MsgReplyShutdown:
		clt.handleReplyShutdown(msg.MsgIdentifier)
		// Release the buffer
		msg.Close()
	case message.MsgSessionNotFound:
		clt.handleSessionNotFound(msg.MsgIdentifier)
		// Release the buffer
		msg.Close()
	case message.MsgMaxSessConnsReached:
		clt.handleMaxSessConnsReached(msg.MsgIdentifier)
		// Release the buffer
		msg.Close()
	case message.MsgSessionsDisabled:
		clt.handleSessionsDisabled(msg.MsgIdentifier)
		// Release the buffer
		msg.Close()
	case message.MsgErrorReply:
		// The message name contains the error code in case of
		// error reply messages, while the UTF8 encoded error message is
		// contained in the message payload
		clt.requestManager.Fail(msg.MsgIdentifier, wwrerr.RequestErr{
			Code:    string(msg.MsgName),
			Message: string(msg.MsgPayload.Data),
		})
		// Release the buffer
		msg.Close()
	case message.MsgInternalError:
		clt.handleInternalError(msg.MsgIdentifier)
		// Release the buffer
		msg.Close()

	case message.MsgSignalBinary:
		fallthrough
	case message.MsgSignalUtf8:
		fallthrough
	case message.MsgSignalUtf16:
		clt.impl.OnSignal(msg)
		// Realease the buffer after the OnSignal user-space hook is executed
		// because it's referenced there through the payload
		msg.Close()

	case message.MsgSessionCreated:
		clt.handleSessionCreated(msg.MsgPayload)
		// Release the buffer after the OnSessionCreated user-space hook is
		// executed because it's referenced there through the payload
		msg.Close()
	case message.MsgSessionClosed:
		// Release the buffer before calling the OnSessionClosed user-space hook
		msg.Close()
		clt.handleSessionClosed()
	default:
		// Release the buffer
		msg.Close()
		clt.warningLog.Printf(
			"Strange message type received: '%d'\n",
			msg.MsgType,
		)
	}
	return nil
}
