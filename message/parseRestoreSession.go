package message

import (
	"errors"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseRestoreSession parses MsgRequestRestoreSession messages
func (msg *Message) parseRestoreSession() error {
	if msg.MsgBuffer.len < MinLenRequestRestoreSession {
		return errors.New(
			"invalid session restoration request message, too short",
		)
	}

	dat := msg.MsgBuffer.Data()

	// Read identifier
	msg.MsgIdentifierBytes = dat[1:9]
	copy(msg.MsgIdentifier[:], msg.MsgIdentifierBytes)

	// Read payload
	msg.MsgPayload = pld.Payload{
		Data: dat[9:],
	}
	return nil
}
