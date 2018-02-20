package webwire

import (
	"fmt"
)

// Error represents an error returned in case of a wrong request
type Error struct {
	Code string `json:"c"`
	Message string `json:"m,omitempty"`
}

// ContextKey defines all values of the context passed to handlers
type ContextKey int
const (
	// MESSAGE represents the received message message
	MESSAGE ContextKey = iota
)

// MessageType defines protocol message types
type MessageType byte
const (
	// MsgRestoreSession is sent by the client
	// to request session restoration
	MsgRestoreSession = 'r'

	// MsgSessionCreated is sent by the server
	// to notify the client about the session creation
	MsgSessionCreated = 'c'

	// MsgSessionClosed is sent by the server
	// to notify the client about the session destruction
	MsgSessionClosed = 'd'

	// MsgSignal is sent by both the client and the server
	// and represents a one-way signal message that doesn't require a reply
	MsgSignal = 's'

	// MsgRequest is sent by the client
	// and represents a roundtrip to the server requiring a reply
	MsgRequest = 'q'

	// MsgReply is sent by the server
	// and represents a reply to a previously sent request
	MsgReply = 'p'

	// MsgErrorReply is sent by the server
	// and represents an error-reply to a previously sent request
	MsgErrorReply = 'e'

	// MsgCloseSession is sent by the client
	// and requests the destruction of the currently active session
	MsgCloseSession = 'x'
)

// Message represents a WebWire protocol message
type Message struct {
	fulfill func(response []byte)
	fail func(Error)
	msgType MessageType
	id *[]byte

	Payload []byte
	Client *Client
}

// ParseMessage tries to parse the byte slice into a typed message object
func ParseMessage(message []byte) (obj Message, err error) {
	if len(message) < 2 {
		return obj, fmt.Errorf("Invalid message (too short)")
	}
	msgType := message[0:1][0]
	switch msgType {

	// Request message must be [1 (type), 1+ (payload)]
	case MsgSignal:
		if len(message) < 2 {
			return obj, fmt.Errorf("Invalid signal message")
		}
		obj.msgType = MsgSignal
		obj.Payload = message[1:]

	// Request message must be [1 (type), 32 (id), | 1+ (payload)]
	case MsgRequest:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgRequest
		obj.Payload = message[33:]

	case MsgReply:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid response message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgReply
		obj.Payload = message[33:]

	// Request message must be [1 (type), 32 (id), | 1+ (payload)]
	case MsgErrorReply:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid error response message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgErrorReply
		obj.Payload = message[33:]

	case MsgRestoreSession:
		// Session activation request message must be
		// [1 (type), 32 (id), | 1+ (payload)]
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid session activation request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgRestoreSession
		obj.Payload = message[33:]

	// Session destruction request message must be [1 (type), 32 (id)]
	case MsgCloseSession:
		if len(message) != 33 {
			return obj, fmt.Errorf("Invalid session destruction request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgCloseSession

	// Session creation request must be [1 (type), 32 (id), | 1+ (payload)]
	case MsgSessionCreated:
		if len(message) < 34 {
			return obj, fmt.Errorf("Invalid session creation request message")
		}
		id := message[1:33]
		obj.id = &id
		obj.msgType = MsgSessionCreated
		obj.Payload = message[33:]
	
	// Session closure request message must be [1 (type), 32 (id)]
	case MsgSessionClosed:
		if len(message) != 1 {
			return obj, fmt.Errorf("Invalid session closure request message")
		}
		obj.msgType = MsgSessionClosed
		obj.Payload = nil

	// Ignore messages of invalid message type
	default:
		return obj, fmt.Errorf("Invalid message type (%d)", rune(msgType))
	}
	return obj, nil
}
