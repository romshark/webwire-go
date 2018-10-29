package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseErrorReply parses MsgErrorReply messages writing the error code into the
// name field and the UTF8 encoded error message into the payload
func (msg *Message) parseErrorReply(message []byte) error {
	if len(message) < MsgMinLenErrorReply {
		return fmt.Errorf("Invalid error reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Read error code length flag
	errCodeLen := int(byte(message[9:10][0]))
	errMessageOffset := 10 + errCodeLen

	// Verify error code length (must be at least 1 character long)
	if errCodeLen < 1 {
		return fmt.Errorf(
			"Invalid error reply message, error code length flag is zero",
		)
	}

	// Verify total message size to prevent segmentation faults
	// caused by inconsistent flags. This could happen if the specified
	// error code length doesn't correspond to the actual length
	// of the provided error code.
	// Subtract 1 character already taken into account by MsgMinLenErrorReply
	if len(message) < MsgMinLenErrorReply+errCodeLen-1 {
		return fmt.Errorf(
			"Invalid error reply message, "+
				"too short for specified code length (%d)",
			errCodeLen,
		)
	}

	// Read UTF8 encoded error message into the payload
	msg.Name = string(message[10 : 10+errCodeLen])
	msg.Payload = pld.Payload{
		Encoding: pld.Utf8,
		Data:     message[errMessageOffset:],
	}
	return nil
}
