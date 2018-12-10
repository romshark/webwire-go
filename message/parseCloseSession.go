package message

import "errors"

// parseCloseSession parses MsgRequestCloseSession messages
func (msg *Message) parseCloseSession() error {
	if msg.MsgBuffer.len != MinLenDoCloseSession {
		return errors.New(
			"invalid session destruction request message, too short",
		)
	}

	// Read identifier
	msg.MsgIdentifierBytes = msg.MsgBuffer.Data()[1:9]
	copy(msg.MsgIdentifier[:], msg.MsgIdentifierBytes)

	return nil
}
