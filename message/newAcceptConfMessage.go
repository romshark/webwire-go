package message

import (
	"encoding/binary"
	"fmt"
	"time"
)

// NewAcceptConfMessage composes a server configuration message and writes it to the
// given buffer
func NewAcceptConfMessage(conf ServerConfiguration) ([]byte, error) {
	buf := make([]byte, MinLenAcceptConf+len(conf.SubProtocolName))

	buf[0] = byte(MsgAcceptConf)
	buf[1] = byte(conf.MajorProtocolVersion)
	buf[2] = byte(conf.MinorProtocolVersion)

	readTimeoutMs := conf.ReadTimeout / time.Millisecond
	if readTimeoutMs > 4294967295 {
		return nil, fmt.Errorf(
			"read timeout (milliseconds) overflow in server conf message (%s)",
			conf.ReadTimeout.String(),
		)
	} else if readTimeoutMs < 0 {
		return nil, fmt.Errorf(
			"negative read timeout (milliseconds) in server conf message (%d)",
			readTimeoutMs,
		)
	}

	binary.LittleEndian.PutUint32(buf[3:7], uint32(readTimeoutMs))
	binary.LittleEndian.PutUint32(buf[7:11], conf.MessageBufferSize)

	copy(buf[11:], conf.SubProtocolName)

	return buf, nil
}
