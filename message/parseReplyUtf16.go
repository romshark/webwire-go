package message

import (
	"errors"

	pld "github.com/qbeon/webwire-go/payload"
)

func (msg *Message) parseReplyUtf16() error {
	if msg.MsgBuffer.len < MinLenReplyUtf16 {
		return errors.New("invalid UTF16 reply message, too short")
	}

	if msg.MsgBuffer.len%2 != 0 {
		return errors.New(
			"unaligned UTF16 encoded reply message " +
				"(probably missing header padding)",
		)
	}

	dat := msg.MsgBuffer.Data()

	// Read identifier
	msg.MsgIdentifierBytes = dat[1:9]
	copy(msg.MsgIdentifier[:], msg.MsgIdentifierBytes)

	// Skip payload if there's none
	if msg.MsgBuffer.len == MinLenReplyUtf16 {
		msg.MsgPayload = pld.Payload{
			Encoding: pld.Utf16,
		}
		return nil
	}

	// Read payload
	msg.MsgPayload = pld.Payload{
		// Take header padding byte into account
		Data: dat[10:],
	}

	return nil
}
