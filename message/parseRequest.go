package message

import (
	"errors"
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseRequest parses MsgRequestBinary and MsgRequestUtf8 messages
func (msg *Message) parseRequest() error {
	if msg.MsgBuffer.len < MsgMinLenRequest {
		return errors.New("invalid request message, too short")
	}

	dat := msg.MsgBuffer.Data()

	// Read identifier
	var id [8]byte
	copy(id[:], dat[1:9])
	msg.MsgIdentifier = id

	// Read name length
	nameLen := int(byte(dat[9:10][0]))
	payloadOffset := 10 + nameLen

	// Verify total message size to prevent segmentation faults caused
	// by inconsistent flags. This could happen if the specified name length
	// doesn't correspond to the actual name length
	if nameLen > 0 {
		// Subtract one to not require the payload but at least the name
		if msg.MsgBuffer.len < MsgMinLenRequest+nameLen-1 {
			return fmt.Errorf(
				"invalid request message, too short for full name (%d)",
				nameLen,
			)
		}

		// Take name into account
		msg.MsgName = dat[10 : 10+nameLen]

		// Read payload if any
		if msg.MsgBuffer.len > MsgMinLenRequest+nameLen-1 {
			msg.MsgPayload = pld.Payload{
				Data: dat[payloadOffset:],
			}
		}
	} else {
		// No name present, expect just the payload to be in place
		msg.MsgPayload = pld.Payload{
			Data: dat[10:],
		}
	}

	return nil
}
