package webwire

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Error represents an error returned in case of a wrong request
type Error struct {
	Code    string `json:"c"`
	Message string `json:"m,omitempty"`
}

// ContextKey defines all values of the context passed to handlers
type ContextKey int

const (
	// MESSAGE represents the received message
	MESSAGE ContextKey = iota
)

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
	fail    func(Error)
	msgType rune
	id      *[]byte

	Payload []byte
	Client  *Client
}

func (msg *Message) parseSignal(message []byte) error {
	if len(message) < 2 {
		return fmt.Errorf("Invalid signal message")
	}
	msg.msgType = MsgSignal
	msg.Payload = message[1:]
	return nil
}

func (msg *Message) parseRequest(message []byte) error {
	if len(message) < 34 {
		return fmt.Errorf("Invalid request message")
	}
	id := message[1:33]
	msg.id = &id
	msg.msgType = MsgRequest
	msg.Payload = message[33:]
	return nil
}

func (msg *Message) parseReply(message []byte) error {
	if len(message) < 34 {
		return fmt.Errorf("Invalid response message")
	}
	id := message[1:33]
	msg.id = &id
	msg.msgType = MsgReply
	msg.Payload = message[33:]
	return nil
}

func (msg *Message) parseErrorReply(message []byte) error {
	if len(message) < 34 {
		return fmt.Errorf("Invalid error reply message")
	}
	id := message[1:33]
	msg.id = &id
	msg.msgType = MsgErrorReply
	msg.Payload = message[33:]
	return nil
}

func (msg *Message) parseRestoreSession(message []byte) error {
	if len(message) < 34 {
		return fmt.Errorf("Invalid session activation request message")
	}
	id := message[1:33]
	msg.id = &id
	msg.msgType = MsgRestoreSession
	msg.Payload = message[33:]
	return nil
}

func (msg *Message) parseCloseSession(message []byte) error {
	if len(message) != 33 {
		return fmt.Errorf("Invalid session destruction request message")
	}
	id := message[1:33]
	msg.id = &id
	msg.msgType = MsgCloseSession
	return nil
}

func (msg *Message) parseSessionCreated(message []byte) error {
	if len(message) < 34 {
		return fmt.Errorf("Invalid session creation request message")
	}
	id := message[1:33]
	msg.id = &id
	msg.msgType = MsgSessionCreated
	msg.Payload = message[33:]
	return nil
}

func (msg *Message) parseSessionClosed(message []byte) error {
	if len(message) != 1 {
		return fmt.Errorf("Invalid session closure request message")
	}
	msg.msgType = MsgSessionClosed
	msg.Payload = nil
	return nil
}

// Parse tries to parse the message from a byte slice
func (msg *Message) Parse(message []byte) error {
	if len(message) < 2 {
		return fmt.Errorf("Invalid message (too short)")
	}
	msgType := message[0:1][0]

	switch msgType {

	// Request message must be: [1 (type), 1+ (payload)]
	case MsgSignal:
		return msg.parseSignal(message)

	// Request message must be: [1 (type), 32 (id), | 1+ (payload)]
	case MsgRequest:
		return msg.parseRequest(message)

	case MsgReply:
		return msg.parseReply(message)

	// Request message must be: [1 (type), 32 (id), | 1+ (payload)]
	case MsgErrorReply:
		return msg.parseErrorReply(message)

	// Session restoration request message must be: [1 (type), 32 (id), | 1+ (payload)]
	case MsgRestoreSession:
		return msg.parseRestoreSession(message)

	// Session destruction request message must be [1 (type), 32 (id)]
	case MsgCloseSession:
		return msg.parseCloseSession(message)

	// Session creation request must be [1 (type), 32 (id), | 1+ (payload)]
	case MsgSessionCreated:
		return msg.parseSessionCreated(message)

	// Session closure request message must be [1 (type), 32 (id)]
	case MsgSessionClosed:
		return msg.parseSessionClosed(message)

	// Ignore messages of invalid message type
	default:
		return fmt.Errorf("Invalid message type (%d)", rune(msgType))
	}
}

func (msg *Message) createFailCallback(client *Client, srv *Server) {
	msg.fail = func(errObj Error) {
		encoded, err := json.Marshal(errObj)
		if err != nil {
			encoded = []byte("CRITICAL: could not encode error report")
		}

		// Send request failure notification
		header := append([]byte("e"), *msg.id...)
		err = client.write(
			websocket.TextMessage,
			append(header, encoded...),
		)
		if err != nil {
			srv.errorLog.Println("Writing failed:", err)
		}
	}
}

func (msg *Message) createReplyCallback(client *Client, srv *Server) {
	msg.fulfill = func(response []byte) {
		// Send response
		header := append([]byte("p"), *msg.id...)
		err := client.write(
			websocket.TextMessage,
			append(header, response...),
		)
		if err != nil {
			srv.errorLog.Println("Writing failed:", err)
		}
	}
}
