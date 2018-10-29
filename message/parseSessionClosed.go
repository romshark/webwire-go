package message

import "fmt"

// parseSessionClosed parses MsgSessionClosed messages
func (msg *Message) parseSessionClosed(message []byte) error {
	if len(message) != MsgMinLenSessionClosed {
		return fmt.Errorf(
			"Invalid session closure notification message, too short",
		)
	}
	return nil
}
