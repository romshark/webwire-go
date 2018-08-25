package message

// NewEmptyRequestMessage composes a new request message
// consisting only of the type and identifier
// and returns its binary representation
func NewEmptyRequestMessage(msgType byte, id [8]byte) (msg []byte) {
	msg = make([]byte, 9)

	// Write message type flag
	msg[0] = msgType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = id[i]
	}

	return msg
}
