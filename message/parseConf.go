package message

import (
	"encoding/binary"
	"fmt"
	"time"
)

// parseConf parses MsgConf messages
func (msg *Message) parseConf(message []byte) error {
	if len(message) < MsgMinLenConf {
		return fmt.Errorf("invalid msg length, too short")
	}
	msg.ServerConfiguration = ServerConfiguration{
		MajorProtocolVersion: message[1:2][0],
		MinorProtocolVersion: message[2:3][0],
		ReadTimeout: time.Duration(binary.LittleEndian.Uint32(message[3:7])) *
			time.Millisecond,
		ReadBufferSize:  binary.LittleEndian.Uint32(message[7:11]),
		WriteBufferSize: binary.LittleEndian.Uint32(message[11:15]),
	}
	return nil
}
