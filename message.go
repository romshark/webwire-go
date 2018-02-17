package webwire

import (
	"fmt"
)

type Error struct {
	Code string `json:"c"`
	Message string `json:"m,omitempty"`
}

type ContextKey int
const (
	MESSAGE ContextKey = iota
)

type MessageType byte
const (
	MsgTyp_SESS_CREATED = 'c'
	MsgTyp_SESS_RESTORE = 'r'
	MsgTyp_SESS_CLOSED = 'd'
	MsgTyp_SIGNAL = 's'
	MsgTyp_REQUEST = 'q'
	MsgTyp_RESPONSE = 'p'
	MsgTyp_ERROR_RESP = 'e'
	MsgTyp_CLOSE_SESSION = 'x'
)

type Message struct {
	fulfill func(response []byte)
	fail func(Error)
	msgType MessageType
	id *[]byte

	Payload []byte
	Client *Client
}

func ConstructMessage(
	fulfill func(response []byte),
	fail func(Error),
	msgType MessageType,
	id *[]byte,
	payload []byte,
	client *Client,
) Message {
	return Message {
		fulfill,
		fail,
		msgType,
		id,
		payload,
		client,
	}
}

func ParseMessage(message []byte) (obj Message, err error) {
	if len(message) < 2 {
		return obj, fmt.Errorf("Invalid message (too short)")
	}
	var msgType byte = message[0:1][0]
	switch msgType {

	// Request message must be [1 (type), 1+ (payload)]
	case MsgTyp_SIGNAL:
		if len(message) < 2 {
			return obj, fmt.Errorf("Invalid signal message")
		}
		obj.msgType = MsgTyp_SIGNAL
		obj.Payload = message[1:]

	// Request message must be [1 (type), 32 (id), | 1+ (payload)]
	case MsgTyp_REQUEST:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgTyp_REQUEST
		obj.Payload = message[33:]

	case MsgTyp_RESPONSE:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid response message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgTyp_RESPONSE
		obj.Payload = message[33:]

	// Request message must be [1 (type), 32 (id), | 1+ (payload)]
	case MsgTyp_ERROR_RESP:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid error response message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgTyp_ERROR_RESP
		obj.Payload = message[33:]

	case MsgTyp_SESS_RESTORE:
		// Session activation request message must be
		// [1 (type), 32 (id), | 1+ (payload)]
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid session activation request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgTyp_SESS_RESTORE
		obj.Payload = message[33:]

	// Session destruction request message must be [1 (type), 32 (id)]
	case MsgTyp_CLOSE_SESSION:
		if len(message) != 33 {
			return obj, fmt.Errorf("Invalid session destruction request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgTyp_CLOSE_SESSION

	// Session creation request must be [1 (type), 32 (id), | 1+ (payload)]
	case MsgTyp_SESS_CREATED:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid session creation request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgTyp_SESS_CREATED
		obj.Payload = message[33:]
	
	// Session closure request message must be [1 (type), 32 (id)]
	case MsgTyp_SESS_CLOSED:
		if len(message) != 1 {
			return obj, fmt.Errorf("Invalid session closure request message")
		}
		obj.msgType = MsgTyp_SESS_CLOSED
		obj.Payload = nil

	// Ignore messages of invalid message type
	default:
		return obj, fmt.Errorf("Invalid message type (%d)", rune(msgType))
	}
	return obj, nil
}
