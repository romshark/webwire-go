package webwire

import (
	"context"
	"time"

	msg "github.com/qbeon/webwire-go/message"
)

// handleMessage handles incoming messages
func (srv *server) handleMessage(con *connection, message []byte) {
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
		srv.failMsg(con, &parsedMessage, ProtocolErr{})
		return
	}

	// Reset read deadline on valid message
	if err := con.sock.SetReadDeadline(
		time.Now().Add(srv.options.ReadTimeout),
	); err != nil {
		srv.errorLog.Printf("couldn't set read deadline: %s", err)
		return
	}

	// Don't register a task handler for heartbeat messages
	//
	// TODO: probably this check should include any message type that's not
	// handled by handleMessage to avoid registering a handler
	if parsedMessage.Type == msg.MsgHeartbeat {
		return
	}

	// Deregister the handler only if a handler was registered
	if srv.registerHandler(con, &parsedMessage) {
		// This defer is necessary to ensure the handler is deregistered even in
		// such situations when one of the handlers panics
		defer srv.deregisterHandler(con)
	}

	switch parsedMessage.Type {
	case msg.MsgSignalBinary:
		fallthrough
	case msg.MsgSignalUtf8:
		fallthrough
	case msg.MsgSignalUtf16:
		srv.handleSignal(con, &parsedMessage)

	case msg.MsgRequestBinary:
		fallthrough
	case msg.MsgRequestUtf8:
		fallthrough
	case msg.MsgRequestUtf16:
		srv.handleRequest(con, &parsedMessage)

	case msg.MsgRestoreSession:
		srv.handleSessionRestore(con, &parsedMessage)
	case msg.MsgCloseSession:
		srv.handleSessionClosure(con, &parsedMessage)
	}
}

// registerHandler increments the number of currently executed handlers
// for this particular client.
// It blocks if the current number of max concurrent handlers was reached
// and frees only when a handler slot is freed for this handler to be executed
func (srv *server) registerHandler(
	con *connection,
	message *msg.Message,
) bool {
	failMsg := false

	if !con.IsActive() {
		return false
	}

	// Wait for free handler slots
	// if the number of concurrent handlers is limited
	if con.options.ConcurrencyLimit > 0 {
		con.handlerSlots.Acquire(context.Background(), 1)
	}

	srv.opsLock.Lock()
	if srv.shutdown || !con.IsActive() {
		// defer failure due to shutdown of either the server or the connection
		failMsg = true
	} else {
		srv.currentOps++
	}
	srv.opsLock.Unlock()

	if failMsg && message.RequiresReply() {
		// Don't process the message, fail it
		srv.failMsgShutdown(con, message)
		return false
	}

	con.registerTask()
	return true
}

// deregisterHandler decrements the number of currently executed handlers
// and shuts down the server if scheduled and no more operations are left
func (srv *server) deregisterHandler(con *connection) {
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()

	con.deregisterTask()

	// Release a handler slot
	if con.options.ConcurrencyLimit > 0 {
		con.handlerSlots.Release(1)
	}
}

// fulfillMsg fulfills the message sending the reply
func (srv *server) fulfillMsg(
	con *connection,
	message *msg.Message,
	replyPayloadEncoding PayloadEncoding,
	replyPayloadData []byte,
) {
	// Send reply
	if err := con.sock.Write(
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
func (srv *server) failMsg(
	con *connection,
	message *msg.Message,
	reqErr error,
) {
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
	if err := con.sock.Write(replyMsg); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsgShutdown sends request failure reply due to current server shutdown
func (srv *server) failMsgShutdown(con *connection, message *msg.Message) {
	if err := con.sock.Write(msg.NewSpecialRequestReplyMessage(
		msg.MsgReplyShutdown,
		message.Identifier,
	)); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}
