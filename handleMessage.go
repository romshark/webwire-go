package webwire

import (
	"context"

	msg "github.com/qbeon/webwire-go/message"
)

// handleMessage handles incoming messages
func (srv *server) handleMessage(clt *Client, message []byte) {
	// Parse message
	var parsedMessage msg.Message
	msgTypeParsed, parserErr := parsedMessage.Parse(message)
	if !msgTypeParsed {
		// Couldn't determine message type, drop message
		return
	} else if parserErr != nil {
		// Couldn't parse message, protocol error
		srv.warnLog.Println("Parser error:", parserErr)

		// Respond with an error but don't break the connection
		// because protocol errors are not critical errors
		srv.failMsg(clt, &parsedMessage, ProtocolErr{})
		return
	}

	// Deregister the handler only if a handler was registered
	if srv.registerHandler(clt, &parsedMessage) {
		defer srv.deregisterHandler(clt)
	}

	switch parsedMessage.Type {
	case msg.MsgSignalBinary:
		fallthrough
	case msg.MsgSignalUtf8:
		fallthrough
	case msg.MsgSignalUtf16:
		srv.handleSignal(clt, &parsedMessage)

	case msg.MsgRequestBinary:
		fallthrough
	case msg.MsgRequestUtf8:
		fallthrough
	case msg.MsgRequestUtf16:
		srv.handleRequest(clt, &parsedMessage)

	case msg.MsgRestoreSession:
		srv.handleSessionRestore(clt, &parsedMessage)
	case msg.MsgCloseSession:
		srv.handleSessionClosure(clt, &parsedMessage)
	}
}

// registerHandler increments the number of currently executed handlers.
// It blocks if the current number of max concurrent handlers was reached
// and frees only when a handler slot is freed for this handler to be executed
func (srv *server) registerHandler(clt *Client, message *msg.Message) bool {
	failMsg := false

	// Wait for free handler slots
	// if the number of concurrent handler is limited
	if !clt.isActive() {
		return false
	}
	if srv.options.IsConcurrentHandlersLimited() {
		srv.handlerSlots.Acquire(context.Background(), 1)
	}

	srv.opsLock.Lock()
	if srv.shutdown || !clt.isActive() {
		// defer failure due to shutdown of either the server
		// or the client agent
		failMsg = true
	} else {
		srv.currentOps++
	}
	srv.opsLock.Unlock()

	if failMsg && message.RequiresResponse() {
		// Don't process the message, fail it
		srv.failMsgShutdown(clt, message)
		return false
	}

	clt.registerTask()
	return true
}

// deregisterHandler decrements the number of currently executed handlers
// and shuts down the server if scheduled and no more operations are left
func (srv *server) deregisterHandler(clt *Client) {
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()

	clt.deregisterTask()

	// Release a handler slot
	if srv.options.IsConcurrentHandlersLimited() {
		srv.handlerSlots.Release(1)
	}
}

// fulfillMsg fulfills the message sending the reply
func (srv *server) fulfillMsg(
	clt *Client,
	message *msg.Message,
	replyPayloadEncoding PayloadEncoding,
	replyPayloadData []byte,
) {
	// Send reply
	if err := clt.conn.Write(
		msg.NewReplyMessage(
			message.Identifier,
			replyPayloadEncoding,
			replyPayloadData,
		),
	); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsg fails the message returning an error reply
func (srv *server) failMsg(clt *Client, message *msg.Message, reqErr error) {
	// Don't send any failure reply if the type of the message
	// doesn't expect any response
	if !message.RequiresReply() {
		return
	}

	var replyMsg []byte
	switch err := reqErr.(type) {
	case ReqErr:
		replyMsg = msg.NewErrorReplyMessage(
			message.Identifier,
			err.Code,
			err.Message,
		)
	case *ReqErr:
		replyMsg = msg.NewErrorReplyMessage(
			message.Identifier,
			err.Code,
			err.Message,
		)
	case MaxSessConnsReachedErr:
		replyMsg = msg.NewSpecialRequestReplyMessage(
			msg.MsgMaxSessConnsReached,
			message.Identifier,
		)
	case SessNotFoundErr:
		replyMsg = msg.NewSpecialRequestReplyMessage(
			msg.MsgSessionNotFound,
			message.Identifier,
		)
	case SessionsDisabledErr:
		replyMsg = msg.NewSpecialRequestReplyMessage(
			msg.MsgSessionsDisabled,
			message.Identifier,
		)
	case ProtocolErr:
		replyMsg = msg.NewSpecialRequestReplyMessage(
			msg.MsgReplyProtocolError,
			message.Identifier,
		)
	default:
		replyMsg = msg.NewSpecialRequestReplyMessage(
			msg.MsgInternalError,
			message.Identifier,
		)
	}

	// Send request failure notification
	if err := clt.conn.Write(replyMsg); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsgShutdown sends request failure reply due to current server shutdown
func (srv *server) failMsgShutdown(clt *Client, message *msg.Message) {
	if err := clt.conn.Write(msg.NewSpecialRequestReplyMessage(
		msg.MsgReplyShutdown,
		message.Identifier,
	)); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}
