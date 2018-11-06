package message

import (
	"errors"
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseSignalUtf16 parses MsgSignalUtf16 messages
func (msg *Message) parseSignalUtf16() error {
	if msg.MsgBuffer.len < MsgMinLenSignalUtf16 {
		return errors.New("invalid signal message, too short")
	}

	if msg.MsgBuffer.len%2 != 0 {
		return errors.New(
			"Unaligned UTF16 encoded signal message " +
				"(probably missing header padding)",
		)
	}

	dat := msg.MsgBuffer.Data()

	// Read name length
	nameLen := int(byte(dat[1:2][0]))

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
	if msg.MsgBuffer.len < minMsgSize {
		return fmt.Errorf(
			"invalid signal message, too short for full name (%d) "+
				"and the minimum payload (2)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.MsgName = dat[2 : 2+nameLen]
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
