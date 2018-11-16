package message

import (
	"errors"
	"io"
)

// ReadBytes reads and parses the message from the given byte slice
func (msg *Message) ReadBytes(bytes []byte) (typeParsed bool, err error) {
	if len(bytes) < 1 {
		return false, nil
	}
	if len(bytes) > len(msg.MsgBuffer.buf) {
		return false, errors.New("message buffer overflow")
	}
	if !msg.MsgBuffer.IsEmpty() {
		msg.MsgBuffer.Close()
	}
	copy(msg.MsgBuffer.buf[:len(bytes)], bytes)
	msg.MsgBuffer.len = len(bytes)
	return msg.parse()
}

// Read reads and parses the message from the given reader
func (msg *Message) Read(reader io.Reader) (typeParsed bool, err error) {
	if err := msg.MsgBuffer.Read(reader); err != nil {
		return false, err
	}
	return msg.parse()
}
