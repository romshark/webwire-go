package message

import (
	"errors"
)

// parseCloseSession parses MsgDoCloseSession messages
func (msg *Message) parseCloseSession() error {
	if msg.MsgBuffer.len != MsgMinLenCloseSession {
		return errors.New(
			"invalid session destruction request message, too short",
		)
	}

	// Read identifier
	msg.MsgIdentifierBytes = msg.MsgBuffer.Data()[1:9]
	copy(msg.MsgIdentifier[:], msg.MsgIdentifierBytes)

	return nil
}
