package message

import (
	"errors"
)

// parseSpecialReplyMessage parses the following message types:
// MsgReplyShutdown, MsgReplyInternalError, MsgReplySessionNotFound,
// MsgReplyMaxSessConnsReached, MsgReplySessionsDisabled
func (msg *Message) parseSpecialReplyMessage() error {
	if msg.MsgBuffer.len < 9 {
		return errors.New("invalid special reply message, too short")
	}

	// Read identifier
	msg.MsgIdentifierBytes = msg.MsgBuffer.Data()[1:9]
	copy(msg.MsgIdentifier[:], msg.MsgIdentifierBytes)

	return nil
}
