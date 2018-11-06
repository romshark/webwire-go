package message

import "fmt"

// parseHeartbeat parses heartbeat messages
func (msg *Message) parseHeartbeat() error {
	if msg.MsgBuffer.len != 1 {
		return fmt.Errorf(
			"invalid heartbeat message (len: %d)",
			msg.MsgBuffer.len,
		)
	}
	return nil
}
