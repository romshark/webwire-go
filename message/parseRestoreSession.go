package message

import (
	"errors"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseRestoreSession parses MsgRestoreSession messages
func (msg *Message) parseRestoreSession() error {
	if msg.MsgBuffer.len < MsgMinLenRestoreSession {
		return errors.New(
			"invalid session restoration request message, too short",
		)
	}

	dat := msg.MsgBuffer.Data()

	// Read identifier
	var id [8]byte
	copy(id[:], dat[1:9])
	msg.MsgIdentifier = id

	// Read payload
	msg.MsgPayload = pld.Payload{
		Data: dat[9:],
	}
	return nil
}
