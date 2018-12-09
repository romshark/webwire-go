package message

import (
	"errors"
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseRequestUtf16 parses MsgRequestUtf16 messages
func (msg *Message) parseRequestUtf16() error {
	if msg.MsgBuffer.len < MinLenRequestUtf16 {
		return errors.New("invalid request message, too short")
	}

	if msg.MsgBuffer.len%2 != 0 {
		return errors.New(
			"unaligned UTF16 encoded request message " +
				"(probably missing header padding)",
		)
	}

	dat := msg.MsgBuffer.Data()

	// Read identifier
	msg.MsgIdentifierBytes = dat[1:9]
	copy(msg.MsgIdentifier[:], msg.MsgIdentifierBytes)

	// Read name length
	nameLen := int(dat[9])

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
	if msg.MsgBuffer.len > payloadOffset && nameLen%2 != 0 {
		minRequiredMsgSize++
		payloadOffset++
	}

	// Verify total message size to prevent segmentation faults caused
	// by inconsistent flags. This could happen if the specified name length
	// doesn't correspond to the actual name length
	if nameLen > 0 {
		if msg.MsgBuffer.len < minRequiredMsgSize {
			return fmt.Errorf(
				"invalid request message, too short for full name (%d)",
				nameLen,
			)
		}

		// Take name into account
		msg.MsgName = dat[10 : 10+nameLen]

		// Read payload if any
		if msg.MsgBuffer.len > minRequiredMsgSize {
			msg.MsgPayload = pld.Payload{
				Data: dat[payloadOffset:],
			}
		}
	} else {
		// No name present, just payload
		msg.MsgPayload = pld.Payload{
			Data: dat[10:],
		}
	}

	return nil
}
