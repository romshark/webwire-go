package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseRequest parses MsgRequestBinary and MsgRequestUtf8 messages
func (msg *Message) parseRequest(message []byte) error {
	if len(message) < MsgMinLenRequest {
		return fmt.Errorf("Invalid request message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Read name length
	nameLen := int(byte(message[9:10][0]))
	payloadOffset := 10 + nameLen

	// Verify total message size to prevent segmentation faults caused
	// by inconsistent flags. This could happen if the specified name length
	// doesn't correspond to the actual name length
	if nameLen > 0 {
		// Subtract one to not require the payload but at least the name
		if len(message) < MsgMinLenRequest+nameLen-1 {
			return fmt.Errorf(
				"Invalid request message, too short for full name (%d)",
				nameLen,
			)
		}

		// Take name into account
		msg.Name = string(message[10 : 10+nameLen])

		// Read payload if any
		if len(message) > MsgMinLenRequest+nameLen-1 {
			msg.Payload = pld.Payload{
				Data: message[payloadOffset:],
			}
		}
	} else {
		// No name present, expect just the payload to be in place
		msg.Payload = pld.Payload{
			Data: message[10:],
		}
	}

	return nil
}
