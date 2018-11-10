package webwire

import (
	"context"
	"fmt"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// handleMessage parses and handles incoming messages
func (srv *server) handleMessage(
	con *connection,
	msg *message.Message,
) (err error) {
	// Don't register a task handler for heartbeat messages
	//
	// TODO: probably this check should include any message type that's not
	// handled by handleMessage to avoid registering a handler
	if msg.MsgType == message.MsgHeartbeat {
		return nil
	}

	// Deregister the handler only if a handler was registered
	deregister := false
	if srv.registerHandler(con, msg) {
		deregister = true
	}
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
		if deregister {
			srv.deregisterHandler(con)
		}
	}()

	switch msg.MsgType {
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

	return nil
}

// registerHandler increments the number of currently executed handlers
// for this particular client.
// It blocks if the current number of max concurrent handlers was reached
// and frees only when a handler slot is freed for this handler to be executed
func (srv *server) registerHandler(
	con *connection,
	msg *message.Message,
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

	if failMsg && msg.RequiresReply() {
		// Don't process the message, fail it
		srv.failMsgShutdown(con, msg)
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
	replyPayload Payload,
) {
	writer, err := con.sock.GetWriter()
	if err != nil {
		srv.errorLog.Printf(
			"couldn't get writer for connection %p: %s",
			con,
			err,
		)
		return
	}

	if err := message.WriteMsgReply(
		writer,
		msg.MsgIdentifier,
		replyPayload.Encoding,
		replyPayload.Data,
	); err != nil {
		srv.errorLog.Printf(
			"couldn't write reply message for connection %p: %s",
			con,
			err,
		)
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

	writer, err := con.sock.GetWriter()
	if err != nil {
		srv.errorLog.Printf(
			"couldn't get writer for connection %p: %s",
			con,
			err,
		)
		return
	}

	switch err := reqErr.(type) {
	case wwrerr.RequestErr:
		if err := message.WriteMsgErrorReply(
			writer,
			msg.MsgIdentifier,
			[]byte(err.Code),
			[]byte(err.Message),
			true,
		); err != nil {
			srv.errorLog.Println("couldn't write error reply message: ", err)
			return
		}
	case *wwrerr.RequestErr:
		if err := message.WriteMsgErrorReply(
			writer,
			msg.MsgIdentifier,
			[]byte(err.Code),
			[]byte(err.Message),
			true,
		); err != nil {
			srv.errorLog.Println("couldn't write error reply message: ", err)
			return
		}
	case wwrerr.MaxSessConnsReachedErr:
		if err := message.WriteMsgSpecialRequestReply(
			writer,
			message.MsgMaxSessConnsReached,
			msg.MsgIdentifier,
		); err != nil {
			srv.errorLog.Println(
				"couldn't write max sessions reached message: ",
				err,
			)
			return
		}
	case wwrerr.SessionNotFoundErr:
		if err := message.WriteMsgSpecialRequestReply(
			writer,
			message.MsgSessionNotFound,
			msg.MsgIdentifier,
		); err != nil {
			srv.errorLog.Println(
				"couldn't write session not found message: ",
				err,
			)
			return
		}
	case wwrerr.SessionsDisabledErr:
		if err := message.WriteMsgSpecialRequestReply(
			writer,
			message.MsgSessionsDisabled,
			msg.MsgIdentifier,
		); err != nil {
			srv.errorLog.Println(
				"couldn't write sessions disabled message: ",
				err,
			)
			return
		}
	default:
		if err := message.WriteMsgSpecialRequestReply(
			writer,
			message.MsgInternalError,
			msg.MsgIdentifier,
		); err != nil {
			srv.errorLog.Println(
				"couldn't write internal error message: ",
				err,
			)
			return
		}
	}
}

// failMsgShutdown sends request failure reply due to current server shutdown
func (srv *server) failMsgShutdown(con *connection, msg *message.Message) {
	writer, err := con.sock.GetWriter()
	if err != nil {
		srv.errorLog.Printf(
			"couldn't get writer for connection %p: %s",
			con,
			err,
		)
	}

	if err := message.WriteMsgSpecialRequestReply(
		writer,
		message.MsgReplyShutdown,
		msg.MsgIdentifier,
	); err != nil {
		srv.errorLog.Println("failed writing shutdown reply message: ", err)
		return
	}
}
