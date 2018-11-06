package message

import "errors"

// parseSpecialReplyMessage parses the following message types:
// MsgReplyShutdown, MsgInternalError, MsgSessionNotFound,
// MsgMaxSessConnsReached, MsgSessionsDisabled
func (msg *Message) parseSpecialReplyMessage() error {
	if msg.MsgBuffer.len < 9 {
		return errors.New("invalid special reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], msg.MsgBuffer.buf[1:9])
	msg.MsgIdentifier = id

	return nil
}
