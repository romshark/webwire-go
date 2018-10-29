package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseRestoreSession parses MsgRestoreSession messages
func (msg *Message) parseRestoreSession(message []byte) error {
	if len(message) < MsgMinLenRestoreSession {
		return fmt.Errorf(
			"Invalid session restoration request message, too short",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Read payload
	msg.Payload = pld.Payload{
		Data: message[9:],
	}
	return nil
}
