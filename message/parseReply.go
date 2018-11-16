package message

import (
	"errors"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseReply parses MsgReplyBinary and MsgReplyUtf8 messages
func (msg *Message) parseReply() error {
	if msg.MsgBuffer.len < MsgMinLenReply {
		return errors.New("invalid reply message, too short")
	}

	dat := msg.MsgBuffer.Data()

	// Read identifier
	var id [8]byte
	copy(id[:], dat[1:9])
	msg.MsgIdentifier = id

	// Skip payload if there's none
	if msg.MsgBuffer.len == MsgMinLenReply {
		return nil
	}

	// Read payload
	msg.MsgPayload = pld.Payload{
		Data: dat[9:],
	}
	return nil
}
