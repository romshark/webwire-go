package message

import (
	"errors"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseSessionCreated parses MsgSessionCreated messages
func (msg *Message) parseSessionCreated() error {
	if msg.MsgBuffer.len < MsgMinLenSessionCreated {
		return errors.New(
			"invalid session creation notification message, too short",
		)
	}

	msg.MsgPayload = pld.Payload{
		Data: msg.MsgBuffer.Data()[1:],
	}

	return nil
}
