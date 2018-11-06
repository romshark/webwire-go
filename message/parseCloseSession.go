package message

import "errors"

// parseCloseSession parses MsgCloseSession messages
func (msg *Message) parseCloseSession() error {
	if msg.MsgBuffer.len != MsgMinLenCloseSession {
		return errors.New(
			"invalid session destruction request message, too short",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], msg.MsgBuffer.buf[1:9])
	msg.MsgIdentifier = id

	return nil
}
