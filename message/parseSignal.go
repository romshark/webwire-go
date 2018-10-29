package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseSignal parses MsgSignalBinary and MsgSignalUtf8 messages
func (msg *Message) parseSignal(message []byte) error {
	if len(message) < MsgMinLenSignal {
		return fmt.Errorf("Invalid signal message, too short")
	}

	// Read name length
	nameLen := int(byte(message[1:2][0]))
	payloadOffset := 2 + nameLen

	// Verify total message size to prevent segmentation faults
	// caused by inconsistent flags. This could happen if the specified
	// name length doesn't correspond to the actual name length
	if len(message) < MsgMinLenSignal+nameLen {
		return fmt.Errorf(
			"Invalid signal message, too short for full name (%d) "+
				"and the minimum payload (1)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.Name = string(message[2:payloadOffset])
		msg.Payload = pld.Payload{
			Data: message[payloadOffset:],
		}
	} else {
		// No name present, just payload
		msg.Payload = pld.Payload{
			Data: message[2:],
		}
	}
	return nil
}
