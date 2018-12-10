package message

import (
	"encoding/binary"
	"errors"
	"time"
)

// parseAcceptConf parses MsgAcceptConf messages
func (msg *Message) parseAcceptConf() error {
	if msg.MsgBuffer.len < MinLenAcceptConf {
		return errors.New("invalid msg length, too short")
	}
	dat := msg.MsgBuffer.Data()

	subProtocolName := []byte(nil)
	if msg.MsgBuffer.len > MinLenAcceptConf {
		subProtocolName = dat[11:]
	}

	msg.ServerConfiguration = ServerConfiguration{
		MajorProtocolVersion: dat[1:2][0],
		MinorProtocolVersion: dat[2:3][0],
		ReadTimeout: time.Duration(
			binary.LittleEndian.Uint32(dat[3:7]),
		) * time.Millisecond,
		MessageBufferSize: binary.LittleEndian.Uint32(dat[7:11]),
		SubProtocolName:   subProtocolName,
	}
	return nil
}
