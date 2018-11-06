package message

import (
	"encoding/binary"
	"errors"
	"time"
)

// parseConf parses MsgConf messages
func (msg *Message) parseConf() error {
	if msg.MsgBuffer.len < MsgMinLenConf {
		return errors.New("invalid msg length, too short")
	}
	msg.ServerConfiguration = ServerConfiguration{
		MajorProtocolVersion: msg.MsgBuffer.buf[1:2][0],
		MinorProtocolVersion: msg.MsgBuffer.buf[2:3][0],
		ReadTimeout: time.Duration(
			binary.LittleEndian.Uint32(msg.MsgBuffer.buf[3:7]),
		) * time.Millisecond,
		MessageBufferSize: binary.LittleEndian.Uint32(msg.MsgBuffer.buf[7:11]),
	}
	return nil
}
