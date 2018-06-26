package webwire

import (
	"context"
)

// handleMessage handles incoming messages
func (srv *server) handleMessage(clt *Client, message []byte) {
	// Parse message
	var msg Message
	msgTypeParsed, parserErr := msg.Parse(message)
	if !msgTypeParsed {
		// Couldn't determine message type, drop message
		return
	} else if parserErr != nil {
		// Couldn't parse message, protocol error
		srv.warnLog.Println("Parser error:", parserErr)

		// Respond with an error but don't break the connection
		// because protocol errors are not critical errors
		srv.failMsg(clt, &msg, ProtocolErr{})
		return
	}

	// Deregister the handler only if a handler was registered
	if srv.registerHandler(clt, &msg) {
		defer srv.deregisterHandler(clt)
	}

	switch msg.msgType {
	case MsgSignalBinary:
		fallthrough
	case MsgSignalUtf8:
		fallthrough
	case MsgSignalUtf16:
		srv.handleSignal(clt, &msg)

	case MsgRequestBinary:
		fallthrough
	case MsgRequestUtf8:
		fallthrough
	case MsgRequestUtf16:
		srv.handleRequest(clt, &msg)

	case MsgRestoreSession:
		srv.handleSessionRestore(clt, &msg)
	case MsgCloseSession:
		srv.handleSessionClosure(clt, &msg)
	}
}

// registerHandler increments the number of currently executed handlers.
// It blocks if the current number of max concurrent handlers was reached
// and frees only when a handler slot is freed for this handler to be executed
func (srv *server) registerHandler(clt *Client, msg *Message) bool {
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

	if failMsg && msg.RequiresResponse() {
		// Don't process the message, fail it
		srv.failMsgShutdown(clt, msg)
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
func (srv *server) fulfillMsg(clt *Client, msg *Message, reply Payload) {
	// Send reply
	if err := clt.conn.Write(
		NewReplyMessage(msg.id, reply),
	); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsg fails the message returning an error reply
func (srv *server) failMsg(clt *Client, msg *Message, reqErr error) {
	// Don't send any failure reply if the type of the message
	// doesn't expect any response
	if !msg.RequiresReply() {
		return
	}

	var replyMsg []byte
	switch err := reqErr.(type) {
	case ReqErr:
		replyMsg = NewErrorReplyMessage(msg.id, err.Code, err.Message)
	case *ReqErr:
		replyMsg = NewErrorReplyMessage(msg.id, err.Code, err.Message)
	case MaxSessConnsReachedErr:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgMaxSessConnsReached,
			msg.id,
		)
	case SessNotFoundErr:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgSessionNotFound,
			msg.id,
		)
	case SessionsDisabledErr:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgSessionsDisabled,
			msg.id,
		)
	case ProtocolErr:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgReplyProtocolError,
			msg.id,
		)
	default:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgInternalError,
			msg.id,
		)
	}

	// Send request failure notification
	if err := clt.conn.Write(replyMsg); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsgShutdown sends request failure reply due to current server shutdown
func (srv *server) failMsgShutdown(clt *Client, msg *Message) {
	if err := clt.conn.Write(NewSpecialRequestReplyMessage(
		MsgReplyShutdown,
		msg.id,
	)); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}
