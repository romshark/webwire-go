package message

import "errors"

// parseSessionClosed parses MsgNotifySessionClosed messages
func (msg *Message) parseSessionClosed() error {
	if msg.MsgBuffer.len != MinLenNotifySessionClosed {
		return errors.New(
			"invalid session closure notification message, too short",
		)
	}
	return nil
}
