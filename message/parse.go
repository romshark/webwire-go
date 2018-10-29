package message

import (
	pld "github.com/qbeon/webwire-go/payload"
)

// Parse tries to parse the message from a byte slice.
// the returned parsedMsgType is set to false if the message type
// couldn't be determined, otherwise it's set to true.
func (msg *Message) Parse(message []byte) (parsedMsgType bool, err error) {
	if len(message) < 1 {
		return false, nil
	}
	var payloadEncoding pld.Encoding
	msgType := message[0:1][0]

	switch msgType {

	// Server Configuration
	case MsgConf:
		err = msg.parseConf(message)

	// Heartbeat
	case MsgHeartbeat:
		err = msg.parseHeartbeat(message)

	// Request error reply message
	case MsgErrorReply:
		err = msg.parseErrorReply(message)

	// Session creation notification message
	case MsgSessionCreated:
		err = msg.parseSessionCreated(message)

	// Session closure notification message
	case MsgSessionClosed:
		err = msg.parseSessionClosed(message)

	// Session destruction request message
	case MsgCloseSession:
		err = msg.parseCloseSession(message)

	// Signal messages
	case MsgSignalBinary:
		payloadEncoding = pld.Binary
		err = msg.parseSignal(message)
	case MsgSignalUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseSignal(message)
	case MsgSignalUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseSignalUtf16(message)

	// Request messages
	case MsgRequestBinary:
		payloadEncoding = pld.Binary
		err = msg.parseRequest(message)
	case MsgRequestUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseRequest(message)
	case MsgRequestUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseRequestUtf16(message)

	// Reply messages
	case MsgReplyBinary:
		payloadEncoding = pld.Binary
		err = msg.parseReply(message)
	case MsgReplyUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseReply(message)
	case MsgReplyUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseReplyUtf16(message)

	// Session restoration request message
	case MsgRestoreSession:
		err = msg.parseRestoreSession(message)

	// Special reply messages
	case MsgReplyShutdown:
		err = msg.parseSpecialReplyMessage(message)
	case MsgInternalError:
		err = msg.parseSpecialReplyMessage(message)
	case MsgSessionNotFound:
		err = msg.parseSpecialReplyMessage(message)
	case MsgMaxSessConnsReached:
		err = msg.parseSpecialReplyMessage(message)
	case MsgSessionsDisabled:
		err = msg.parseSpecialReplyMessage(message)
	case MsgReplyProtocolError:
		err = msg.parseSpecialReplyMessage(message)

	// Ignore messages of invalid message type
	default:
		return false, nil
	}

	msg.Type = msgType
	msg.Payload.Encoding = payloadEncoding
	return true, err
}
