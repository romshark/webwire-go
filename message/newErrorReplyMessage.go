package message

import "fmt"

// NewErrorReplyMessage composes a new error reply message
// and returns its binary representation
func NewErrorReplyMessage(
	requestIdent [8]byte,
	code,
	message string,
) (msg []byte) {
	if len(code) < 1 {
		panic(fmt.Errorf(
			"Missing error code while creating a new error reply message",
		))
	} else if len(code) > 255 {
		panic(fmt.Errorf(
			"Invalid error code while creating a new error reply message,"+
				"too long (%d)",
			len(code),
		))
	}

	// Determine total message length
	msg = make([]byte, 10+len(code)+len(message))

	// Write message type flag
	msg[0] = MsgErrorReply

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = requestIdent[i]
	}

	// Write code length flag
	msg[9] = byte(len(code))

	// Write error code
	for i := 0; i < len(code); i++ {
		char := code[i]
		if char < 32 || char > 126 {
			panic(fmt.Errorf(
				"Unsupported character in reply error - error code: %s",
				string(char),
			))
		}
		msg[10+i] = code[i]
	}

	errMessageOffset := 10 + len(code)

	// Write error message
	for i := 0; i < len(message); i++ {
		msg[errMessageOffset+i] = message[i]
	}

	return msg
}
