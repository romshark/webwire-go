package webwire

import (
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
)

// NewMessageWrapper creates a new Message interface compliant message object
func NewMessageWrapper(message *msg.Message) *MessageWrapper {
	return &MessageWrapper{
		actual: message,
	}
}

// MessageWrapper wraps a msg.Message pointer
// to make it implement the Message interface
type MessageWrapper struct {
	actual *msg.Message
}

// MessageType implements the Message interface
func (wrp *MessageWrapper) MessageType() byte {
	return wrp.actual.Type
}

// Identifier implements the Message interface
func (wrp *MessageWrapper) Identifier() [8]byte {
	return wrp.actual.Identifier
}

// Name implements the Message interface
func (wrp *MessageWrapper) Name() string {
	return wrp.actual.Name
}

// Payload implements the Message interface
func (wrp *MessageWrapper) Payload() Payload {
	return &EncodedPayload{
		Payload: pld.Payload{
			Encoding: wrp.actual.Payload.Encoding,
			Data:     wrp.actual.Payload.Data,
		},
	}
}
