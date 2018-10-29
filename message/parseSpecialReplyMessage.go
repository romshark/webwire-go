package message

import "fmt"

// parseSpecialReplyMessage parses the following message types:
// MsgReplyShutdown, MsgInternalError, MsgSessionNotFound,
// MsgMaxSessConnsReached, MsgSessionsDisabled, MsgReplyProtocolError
func (msg *Message) parseSpecialReplyMessage(message []byte) error {
	if len(message) < 9 {
		return fmt.Errorf("Invalid special reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.Identifier = id

	return nil
}
