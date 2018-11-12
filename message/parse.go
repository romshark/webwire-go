package message

import (
	"errors"

	pld "github.com/qbeon/webwire-go/payload"
)

// parse tries to parse the message from a byte slice.
// the returned parsedMsgType is set to false if the message type
// couldn't be determined, otherwise it's set to true.
func (msg *Message) parse() (parsedMsgType bool, err error) {
	if msg.MsgBuffer.IsEmpty() {
		return false, nil
	}
	var payloadEncoding pld.Encoding
	msgType := msg.MsgBuffer.buf[0:1][0]

	switch msgType {

	// Server Configuration
	case MsgConf:
		err = msg.parseConf()

	// Heartbeat
	case MsgHeartbeat:
		err = msg.parseHeartbeat()

	// Request error reply message
	case MsgErrorReply:
		err = msg.parseErrorReply()

	// Session creation notification message
	case MsgSessionCreated:
		err = msg.parseSessionCreated()

	// Session closure notification message
	case MsgSessionClosed:
		err = msg.parseSessionClosed()

	// Session destruction request message
	case MsgCloseSession:
		err = msg.parseCloseSession()

	// Signal messages
	case MsgSignalBinary:
		payloadEncoding = pld.Binary
		err = msg.parseSignal()
	case MsgSignalUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseSignal()
	case MsgSignalUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseSignalUtf16()

	// Request messages
	case MsgRequestBinary:
		payloadEncoding = pld.Binary
		err = msg.parseRequest()
	case MsgRequestUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseRequest()
	case MsgRequestUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseRequestUtf16()

	// Reply messages
	case MsgReplyBinary:
		payloadEncoding = pld.Binary
		err = msg.parseReply()
	case MsgReplyUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseReply()
	case MsgReplyUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseReplyUtf16()

	// Session restoration request message
	case MsgRestoreSession:
		err = msg.parseRestoreSession()

	// Special reply messages
	case MsgReplyShutdown:
		err = msg.parseSpecialReplyMessage()
	case MsgInternalError:
		err = msg.parseSpecialReplyMessage()
	case MsgSessionNotFound:
		err = msg.parseSpecialReplyMessage()
	case MsgMaxSessConnsReached:
		err = msg.parseSpecialReplyMessage()
	case MsgSessionsDisabled:
		err = msg.parseSpecialReplyMessage()

	// Ignore messages of invalid message type
	default:
		return false, errors.New("invalid message type")
	}

	msg.MsgType = msgType
	msg.MsgPayload.Encoding = payloadEncoding
	return true, err
}
