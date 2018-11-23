package webwire

import "github.com/qbeon/webwire-go/message"

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

	if !srv.registerHandler(con, msg) {
		return nil
	}

	switch msg.MsgType {
	case message.MsgSignalBinary,
		message.MsgSignalUtf8,
		message.MsgSignalUtf16:
		if con.options.ConcurrencyLimit < 0 ||
			con.options.ConcurrencyLimit > 1 {
			go srv.handleSignal(con, msg)
		} else {
			srv.handleSignal(con, msg)
		}

	case message.MsgRequestBinary,
		message.MsgRequestUtf8,
		message.MsgRequestUtf16:
		if con.options.ConcurrencyLimit < 0 ||
			con.options.ConcurrencyLimit > 1 {
			go srv.handleRequest(con, msg)
		} else {
			srv.handleRequest(con, msg)
		}

	case message.MsgRestoreSession:
		srv.handleSessionRestore(con, msg)
	case message.MsgCloseSession:
		srv.handleSessionClosure(con, msg)
	default:
		// Immediately deregister handlers for unexpected message types
		srv.deregisterHandler(con)
	}

	return nil
}
