package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

func (msg *Message) parseReplyUtf16(message []byte) error {
	if len(message) < MsgMinLenReplyUtf16 {
		return fmt.Errorf("Invalid UTF16 reply message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded reply message " +
				"(probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Skip payload if there's none
	if len(message) == MsgMinLenReplyUtf16 {
		msg.Payload = pld.Payload{
			Encoding: pld.Utf16,
		}
		return nil
	}

	// Read payload
	msg.Payload = pld.Payload{
		// Take header padding byte into account
		Data: message[10:],
	}
	return nil
}
