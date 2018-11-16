package message

import (
	"errors"
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseSignal parses MsgSignalBinary and MsgSignalUtf8 messages
func (msg *Message) parseSignal() error {
	if msg.MsgBuffer.len < MsgMinLenSignal {
		return errors.New("invalid signal message, too short")
	}

	dat := msg.MsgBuffer.Data()

	// Read name length
	nameLen := int(byte(dat[1:2][0]))
	payloadOffset := 2 + nameLen

	// Verify total message size to prevent segmentation faults
	// caused by inconsistent flags. This could happen if the specified
	// name length doesn't correspond to the actual name length
	if msg.MsgBuffer.len < MsgMinLenSignal+nameLen {
		return fmt.Errorf(
			"invalid signal message, too short for full name (%d) "+
				"and the minimum payload (1)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.MsgName = dat[2:payloadOffset]
		msg.MsgPayload = pld.Payload{
			Data: dat[payloadOffset:],
		}
	} else {
		// No name present, just payload
		msg.MsgPayload = pld.Payload{
			Data: dat[2:],
		}
	}
	return nil
}
