package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseSignalUtf16 parses MsgSignalUtf16 messages
func (msg *Message) parseSignalUtf16(message []byte) error {
	if len(message) < MsgMinLenSignalUtf16 {
		return fmt.Errorf("Invalid signal message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded signal message " +
				"(probably missing header padding)",
		)
	}

	// Read name length
	nameLen := int(byte(message[1:2][0]))

	// Determine minimum required message length
	minMsgSize := MsgMinLenSignalUtf16 + nameLen
	payloadOffset := 2 + nameLen

	// Check whether a name padding byte is to be expected
	if nameLen%2 != 0 {
		minMsgSize++
		payloadOffset++
	}

	// Verify total message size to prevent segmentation faults
	// caused by inconsistent flags. This could happen if the specified
	// name length doesn't correspond to the actual name length
	if len(message) < minMsgSize {
		return fmt.Errorf(
			"Invalid signal message, too short for full name (%d) "+
				"and the minimum payload (2)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.Name = string(message[2 : 2+nameLen])
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
