package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// Parse tries to parse the message from a byte slice.
// the returned parsedMsgType is set to false if the message type
// couldn't be determined, otherwise it's set to true.
func (msg *Message) Parse(message []byte) (parsedMsgType bool, err error) {
	if len(message) < 1 {
		return false, nil
	}
	var payloadEncoding pld.Encoding
	msgType := message[0:1][0]

	switch msgType {

	// Request error reply message
	case MsgErrorReply:
		err = msg.parseErrorReply(message)

	// Session creation notification message
	case MsgSessionCreated:
		err = msg.parseSessionCreated(message)

	// Session closure notification message
	case MsgSessionClosed:
		err = msg.parseSessionClosed(message)

	// Session destruction request message
	case MsgCloseSession:
		err = msg.parseCloseSession(message)

	// Signal messages
	case MsgSignalBinary:
		payloadEncoding = pld.Binary
		err = msg.parseSignal(message)
	case MsgSignalUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseSignal(message)
	case MsgSignalUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseSignalUtf16(message)

	// Request messages
	case MsgRequestBinary:
		payloadEncoding = pld.Binary
		err = msg.parseRequest(message)
	case MsgRequestUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseRequest(message)
	case MsgRequestUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseRequestUtf16(message)

	// Reply messages
	case MsgReplyBinary:
		payloadEncoding = pld.Binary
		err = msg.parseReply(message)
	case MsgReplyUtf8:
		payloadEncoding = pld.Utf8
		err = msg.parseReply(message)
	case MsgReplyUtf16:
		payloadEncoding = pld.Utf16
		err = msg.parseReplyUtf16(message)

	// Session restoration request message
	case MsgRestoreSession:
		err = msg.parseRestoreSession(message)

	// Special reply messages
	case MsgReplyShutdown:
		err = msg.parseSpecialReplyMessage(message)
	case MsgInternalError:
		err = msg.parseSpecialReplyMessage(message)
	case MsgSessionNotFound:
		err = msg.parseSpecialReplyMessage(message)
	case MsgMaxSessConnsReached:
		err = msg.parseSpecialReplyMessage(message)
	case MsgSessionsDisabled:
		err = msg.parseSpecialReplyMessage(message)
	case MsgReplyProtocolError:
		err = msg.parseSpecialReplyMessage(message)

	// Ignore messages of invalid message type
	default:
		return false, nil
	}

	msg.Type = msgType
	msg.Payload.Encoding = payloadEncoding
	return true, err
}

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

func (msg *Message) parseRequestUtf16(message []byte) error {
	if len(message) < MsgMinLenRequestUtf16 {
		return fmt.Errorf("Invalid request message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded request message " +
				"(probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Read name length
	nameLen := int(byte(message[9:10][0]))

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
	if len(message) > payloadOffset && nameLen%2 != 0 {
		minRequiredMsgSize++
		payloadOffset++
	}

	// Verify total message size to prevent segmentation faults caused
	// by inconsistent flags. This could happen if the specified name length
	// doesn't correspond to the actual name length
	if nameLen > 0 {
		if len(message) < minRequiredMsgSize {
			return fmt.Errorf(
				"Invalid request message, too short for full name (%d)",
				nameLen,
			)
		}

		// Take name into account
		msg.Name = string(message[10 : 10+nameLen])

		// Read payload if any
		if len(message) > minRequiredMsgSize {
			msg.Payload = pld.Payload{
				Data: message[payloadOffset:],
			}
		}
	} else {
		// No name present, just payload
		msg.Payload = pld.Payload{
			Data: message[10:],
		}
	}

	return nil
}

func (msg *Message) parseReply(message []byte) error {
	if len(message) < MsgMinLenReply {
		return fmt.Errorf("Invalid reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Skip payload if there's none
	if len(message) == MsgMinLenReply {
		return nil
	}

	// Read payload
	msg.Payload = pld.Payload{
		Data: message[9:],
	}
	return nil
}

func (msg *Message) parseReplyUtf16(message []byte) error {
	if len(message) < MsgMinLenReplyUtf16 {
		return fmt.Errorf("Invalid UTF16 reply message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded reply message " +
				"(probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Skip payload if there's none
	if len(message) == MsgMinLenReplyUtf16 {
		msg.Payload = pld.Payload{
			Encoding: pld.Utf16,
		}
		return nil
	}

	// Read payload
	msg.Payload = pld.Payload{
		// Take header padding byte into account
		Data: message[10:],
	}
	return nil
}

// parseErrorReply parses the given message assuming it's an error reply message
// parsing the error code into the name field
// and the UTF8 encoded error message into the payload
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

func (msg *Message) parseRestoreSession(message []byte) error {
	if len(message) < MsgMinLenRestoreSession {
		return fmt.Errorf(
			"Invalid session restoration request message, too short",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	// Read payload
	msg.Payload = pld.Payload{
		Data: message[9:],
	}
	return nil
}

func (msg *Message) parseCloseSession(message []byte) error {
	if len(message) != MsgMinLenCloseSession {
		return fmt.Errorf(
			"Invalid session destruction request message, too short",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	return nil
}

func (msg *Message) parseSessionCreated(message []byte) error {
	if len(message) < MsgMinLenSessionCreated {
		return fmt.Errorf(
			"Invalid session creation notification message, too short",
		)
	}

	msg.Payload = pld.Payload{
		Data: message[1:],
	}
	return nil
}

func (msg *Message) parseSessionClosed(message []byte) error {
	if len(message) != MsgMinLenSessionClosed {
		return fmt.Errorf(
			"Invalid session closure notification message, too short",
		)
	}
	return nil
}

func (msg *Message) parseSpecialReplyMessage(message []byte) error {
	if len(message) < 9 {
		return fmt.Errorf("Invalid special reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	return nil
}
