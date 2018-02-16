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
	SESS_CREATION = 'c'
	SESS_RESTORE = 'r'
	SESS_CLOSURE = 'd'
	SIGNAL = 's'
	REQUEST = 'q'
	RESPONSE = 'p'
	ERROR_RESP = 'e'
)

type Message struct {
	fulfill func(response []byte)
	fail func(Error)
	msgType MessageType
	id *[]byte
	payload []byte
	session *Session
}

func ConstructMessage(
	fulfill func(response []byte),
	fail func(Error),
	msgType MessageType,
	id *[]byte,
	payload []byte,
	session *Session,
) Message {
	return Message {
		fulfill,
		fail,
		msgType,
		id,
		payload,
		session,
	}
}

func (msg *Message) Payload() []byte {
	return msg.payload
}

func (msg *Message) Session() *Session {
	return msg.session
}

func ParseMessage(message []byte) (obj Message, err error) {
	if len(message) < 2 {
		return obj, fmt.Errorf("Invalid message (too short)")
	}
	var msgType byte = message[0:1][0]
	switch msgType {

	// Request message must be [1 (type), 1+ (payload)]
	case SIGNAL:
		if len(message) < 2 {
			return obj, fmt.Errorf("Invalid signal message")
		}
		obj.msgType = SIGNAL
		obj.payload = message[1:]

	// Request message must be [1 (type), 32 (id), | 1+ (payload)]
	case REQUEST:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = REQUEST
		obj.payload = message[33:]

	case RESPONSE:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid response message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = RESPONSE
		obj.payload = message[33:]

	// Request message must be [1 (type), 32 (id), | 1+ (payload)]
	case ERROR_RESP:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid error response message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = ERROR_RESP
		obj.payload = message[33:]

	case SESS_RESTORE:
		// Session activation request message must be
		// [1 (type), 32 (id), | 1+ (payload)]
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid session activation request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = SESS_RESTORE
		obj.payload = message[33:]

	// Session creation request must be [1 (type), 32 (id), | 1+ (payload)]
	case SESS_CREATION:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid session creation request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = SESS_CREATION
		obj.payload = message[33:]
	
	// Session closure request message must be [1 (type), 32 (id)]
	case SESS_CLOSURE:
		if len(message) < 33 {
			return obj, fmt.Errorf("Invalid session closure request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = SESS_CLOSURE
		obj.payload = nil

	// Ignore messages of invalid message type
	default:
		return obj, fmt.Errorf("Invalid message type (%d)", rune(msgType))
	}
	return obj, nil
}
