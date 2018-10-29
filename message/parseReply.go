package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseReply parses MsgReplyBinary and MsgReplyUtf8 messages
func (msg *Message) parseReply(message []byte) error {
	if len(message) < MsgMinLenReply {
		return fmt.Errorf("Invalid reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Skip payload if there's none
	if len(message) == MsgMinLenReply {
		return nil
	}

	// Read payload
	msg.Payload = pld.Payload{
		Data: message[9:],
	}
	return nil
}
