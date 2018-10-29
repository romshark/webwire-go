package message

import (
	"encoding/binary"
	"fmt"
	"time"
)

// NewConfMessage composes a server configuration message
func NewConfMessage(conf ServerConfiguration) ([]byte, error) {
	msg := make([]byte, 15)

	msg[0] = byte(MsgConf)
	msg[1] = byte(conf.MajorProtocolVersion)
	msg[2] = byte(conf.MinorProtocolVersion)

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

	binary.LittleEndian.PutUint32(msg[3:7], uint32(readTimeoutMs))
	binary.LittleEndian.PutUint32(msg[7:11], conf.ReadBufferSize)
	binary.LittleEndian.PutUint32(msg[11:15], conf.WriteBufferSize)

	return msg, nil
}
