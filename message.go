package webwire

import (
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
)

type MessageWrapper struct {
	actual *msg.Message
}

// MessageType returns the type of the message
func (wrp *MessageWrapper) MessageType() byte {
	return wrp.actual.Type
}

// Identifier returns the message identifier
func (wrp *MessageWrapper) Identifier() [8]byte {
	return wrp.actual.Identifier
}

// Name returns the name of the message
func (wrp *MessageWrapper) Name() string {
	return wrp.actual.Name
}

// Payload returns the message payload
func (wrp *MessageWrapper) Payload() Payload {
	return &EncodedPayload{
		Payload: pld.Payload{
			Encoding: wrp.actual.Payload.Encoding,
			Data:     wrp.actual.Payload.Data,
		},
	}
}
