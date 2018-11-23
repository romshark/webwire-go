package webwire

import (
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

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
			msg.MsgIdentifierBytes,
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
			msg.MsgIdentifierBytes,
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
			msg.MsgIdentifierBytes,
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
			msg.MsgIdentifierBytes,
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
			msg.MsgIdentifierBytes,
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
			msg.MsgIdentifierBytes,
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
		msg.MsgIdentifierBytes,
	); err != nil {
		srv.errorLog.Println("failed writing shutdown reply message: ", err)
		return
	}
}
