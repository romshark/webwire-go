package message

// NewNamelessRequestMessage composes a new nameless (initially without a name)
// request message and returns its binary representation
func NewNamelessRequestMessage(
	reqType byte,
	identifier [8]byte,
	binaryPayload []byte,
) (msg []byte) {
	// 9 byte header + n bytes payload
	msg = make([]byte, 9+len(binaryPayload))

	// Write message type flag
	msg[0] = reqType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = identifier[i]
	}

	// Write payload
	for i := 0; i < len(binaryPayload); i++ {
		msg[9+i] = binaryPayload[i]
	}

	return msg
}
