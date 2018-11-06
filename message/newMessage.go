package message

// NewMessage creates a new buffered message instance
func NewMessage(bufferSize uint32) *Message {
	return &Message{
		MsgBuffer: Buffer{
			buf: make([]byte, bufferSize),
			len: 0,
		},
	}
}
