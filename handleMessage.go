package webwire

// handleMessage handles incoming messages
func (srv *server) handleMessage(clt *Client, msg *Message) error {
	switch msg.msgType {
	case MsgSignalBinary:
		fallthrough
	case MsgSignalUtf8:
		fallthrough
	case MsgSignalUtf16:
		srv.handleSignal(clt, msg)

	case MsgRequestBinary:
		fallthrough
	case MsgRequestUtf8:
		fallthrough
	case MsgRequestUtf16:
		srv.handleRequest(clt, msg)

	case MsgRestoreSession:
		return srv.handleSessionRestore(clt, msg)
	case MsgCloseSession:
		return srv.handleSessionClosure(clt, msg)
	}
	return nil
}

// fulfillMsg filfills the message sending the reply
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
