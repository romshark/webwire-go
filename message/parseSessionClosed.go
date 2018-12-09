package message

import "errors"

// parseSessionClosed parses MsgNotifySessionClosed messages
func (msg *Message) parseSessionClosed() error {
	if msg.MsgBuffer.len != MsgMinLenSessionClosed {
		return errors.New(
			"invalid session closure notification message, too short",
		)
	}
	return nil
}
