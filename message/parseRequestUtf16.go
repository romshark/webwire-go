package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseRequestUtf16 parses MsgRequestUtf16 messages
func (msg *Message) parseRequestUtf16(message []byte) error {
	if len(message) < MsgMinLenRequestUtf16 {
		return fmt.Errorf("Invalid request message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded request message " +
				"(probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Read name length
	nameLen := int(byte(message[9:10][0]))

	// Determine minimum required message length.
	// There's at least a 10 byte header and a 2 byte payload expected
	minRequiredMsgSize := 12
	if nameLen > 0 {
		// ...unless a name is given, in which case the payload isn't required
		minRequiredMsgSize = 10 + nameLen
	}

	// A header padding byte is only expected, when there's a payload
	// beyond the name. It's not required if there's just the header and a name
	payloadOffset := 10 + nameLen
	if len(message) > payloadOffset && nameLen%2 != 0 {
		minRequiredMsgSize++
		payloadOffset++
	}

	// Verify total message size to prevent segmentation faults caused
	// by inconsistent flags. This could happen if the specified name length
	// doesn't correspond to the actual name length
	if nameLen > 0 {
		if len(message) < minRequiredMsgSize {
			return fmt.Errorf(
				"Invalid request message, too short for full name (%d)",
				nameLen,
			)
		}

		// Take name into account
		msg.Name = message[10 : 10+nameLen]

		// Read payload if any
		if len(message) > minRequiredMsgSize {
			msg.Payload = pld.Payload{
				Data: message[payloadOffset:],
			}
		}
	} else {
		// No name present, just payload
		msg.Payload = pld.Payload{
			Data: message[10:],
		}
	}

	return nil
}
