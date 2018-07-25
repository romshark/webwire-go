package message

import "fmt"

// NewSpecialRequestReplyMessage composes a new special request reply message
func NewSpecialRequestReplyMessage(msgType byte, reqIdent [8]byte) []byte {
	switch msgType {
	case MsgInternalError:
		break
	case MsgMaxSessConnsReached:
		break
	case MsgSessionNotFound:
		break
	case MsgSessionsDisabled:
		break
	case MsgReplyShutdown:
		break
	case MsgReplyProtocolError:
		break
	default:
		panic(fmt.Errorf(
			"Message type (%d) doesn't represent a special reply message",
			msgType,
		))
	}

	msg := make([]byte, 9)

	// Write message type flag
	msg[0] = msgType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = reqIdent[i]
	}

	return msg
}
