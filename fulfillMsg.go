package webwire

import "github.com/qbeon/webwire-go/message"

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
		msg.MsgIdentifierBytes,
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
