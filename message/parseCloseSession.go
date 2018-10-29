package message

import "fmt"

// parseCloseSession parses MsgCloseSession messages
func (msg *Message) parseCloseSession(message []byte) error {
	if len(message) != MsgMinLenCloseSession {
		return fmt.Errorf(
			"Invalid session destruction request message, too short",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	return nil
}
