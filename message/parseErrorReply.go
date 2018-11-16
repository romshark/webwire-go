package message

import (
	"errors"
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseErrorReply parses MsgErrorReply messages writing the error code into the
// name field and the UTF8 encoded error message into the payload
func (msg *Message) parseErrorReply() error {
	if msg.MsgBuffer.len < MsgMinLenErrorReply {
		return errors.New("invalid error reply message, too short")
	}

	dat := msg.MsgBuffer.Data()

	// Read identifier
	var id [8]byte
	copy(id[:], dat[1:9])
	msg.MsgIdentifier = id

	// Read error code length flag
	errCodeLen := int(byte(dat[9:10][0]))
	errMessageOffset := 10 + errCodeLen

	// Verify error code length (must be at least 1 character long)
	if errCodeLen < 1 {
		return errors.New(
			"invalid error reply message, error code length flag is zero",
		)
	}

	// Verify total message size to prevent segmentation faults
	// caused by inconsistent flags. This could happen if the specified
	// error code length doesn't correspond to the actual length
	// of the provided error code.
	// Subtract 1 character already taken into account by MsgMinLenErrorReply
	if msg.MsgBuffer.len < MsgMinLenErrorReply+errCodeLen-1 {
		return fmt.Errorf(
			"invalid error reply message, "+
				"too short for specified code length (%d)",
			errCodeLen,
		)
	}

	// Read UTF8 encoded error message into the payload
	msg.MsgName = dat[10 : 10+errCodeLen]
	msg.MsgPayload = pld.Payload{
		Encoding: pld.Utf8,
		Data:     dat[errMessageOffset:],
	}
	return nil
}
