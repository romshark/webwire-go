package webwire

import (
	"context"
	"time"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/msgbuf"
)

// handleMessage parses and handles incoming messages
func (srv *server) handleMessage(con *connection, buf *msgbuf.MessageBuffer) {
	defer buf.Close()

	msg := &message.Message{}

	// Parse message
	msgTypeParsed, parserErr := msg.Parse(buf.Data())
	if !msgTypeParsed {
		// Couldn't determine message type, drop message
		return
	} else if parserErr != nil {
		// Couldn't parse message, protocol error
		srv.warnLog.Println("Parser error:", parserErr)

		// Respond with an error but don't break the connection
		// because protocol errors are not critical errors
		srv.failMsg(con, msg, ProtocolErr{})
		return
	}

	// Reset the read deadline on valid message
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
	if msg.Type == message.MsgHeartbeat {
		return
	}

	// Deregister the handler only if a handler was registered
	if srv.registerHandler(con, msg) {
		// This defer is necessary to ensure the handler is deregistered even in
		// such situations when one of the handlers panics
		defer srv.deregisterHandler(con)
	}

	switch msg.Type {
	case message.MsgSignalBinary:
		fallthrough
	case message.MsgSignalUtf8:
		fallthrough
	case message.MsgSignalUtf16:
		srv.handleSignal(con, msg)

	case message.MsgRequestBinary:
		fallthrough
	case message.MsgRequestUtf8:
		fallthrough
	case message.MsgRequestUtf16:
		srv.handleRequest(con, msg)

	case message.MsgRestoreSession:
		srv.handleSessionRestore(con, msg)
	case message.MsgCloseSession:
		srv.handleSessionClosure(con, msg)
	}
}

// registerHandler increments the number of currently executed handlers
// for this particular client.
// It blocks if the current number of max concurrent handlers was reached
// and frees only when a handler slot is freed for this handler to be executed
func (srv *server) registerHandler(
	con *connection,
	message *message.Message,
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
	msg *message.Message,
	replyPayloadEncoding PayloadEncoding,
	replyPayloadData []byte,
) {
	// Send reply
	if err := con.sock.Write(
		message.NewReplyMessage(
			msg.Identifier,
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
	msg *message.Message,
	reqErr error,
) {
	// Don't send any failure reply if the type of the message
	// doesn't expect any response
	if !msg.RequiresReply() {
		return
	}

	var replyMsg []byte
	switch err := reqErr.(type) {
	case ReqErr:
		replyMsg = message.NewErrorReplyMessage(
			msg.Identifier,
			err.Code,
			err.Message,
		)
	case *ReqErr:
		replyMsg = message.NewErrorReplyMessage(
			msg.Identifier,
			err.Code,
			err.Message,
		)
	case MaxSessConnsReachedErr:
		replyMsg = message.NewSpecialRequestReplyMessage(
			message.MsgMaxSessConnsReached,
			msg.Identifier,
		)
	case SessNotFoundErr:
		replyMsg = message.NewSpecialRequestReplyMessage(
			message.MsgSessionNotFound,
			msg.Identifier,
		)
	case SessionsDisabledErr:
		replyMsg = message.NewSpecialRequestReplyMessage(
			message.MsgSessionsDisabled,
			msg.Identifier,
		)
	case ProtocolErr:
		replyMsg = message.NewSpecialRequestReplyMessage(
			message.MsgReplyProtocolError,
			msg.Identifier,
		)
	default:
		replyMsg = message.NewSpecialRequestReplyMessage(
			message.MsgInternalError,
			msg.Identifier,
		)
	}

	// Send request failure notification
	if err := con.sock.Write(replyMsg); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsgShutdown sends request failure reply due to current server shutdown
func (srv *server) failMsgShutdown(con *connection, msg *message.Message) {
	if err := con.sock.Write(message.NewSpecialRequestReplyMessage(
		message.MsgReplyShutdown,
		msg.Identifier,
	)); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}
