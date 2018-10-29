package message

import "fmt"

// parseHeartbeat parses heartbeat messages
func (msg *Message) parseHeartbeat(message []byte) error {
	if len(message) != 1 {
		return fmt.Errorf("Invalid heartbeat message (len: %d)", len(message))
	}
	return nil
}
