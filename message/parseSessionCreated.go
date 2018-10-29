package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// parseSessionCreated parses MsgSessionCreated messages
func (msg *Message) parseSessionCreated(message []byte) error {
	if len(message) < MsgMinLenSessionCreated {
		return fmt.Errorf(
			"Invalid session creation notification message, too short",
		)
	}

	msg.Payload = pld.Payload{
		Data: message[1:],
	}
	return nil
}
