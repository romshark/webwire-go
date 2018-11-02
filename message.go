package webwire

import (
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
)

// newMessageWrapper creates a new Message interface compliant message object
func newMessageWrapper(message *msg.Message) *messageWrapper {
	return &messageWrapper{
		actual: message,
	}
}

// messageWrapper wraps a msg.Message pointer
// to make it implement the Message interface
type messageWrapper struct {
	actual *msg.Message
}

// MessageType implements the Message interface
func (wrp *messageWrapper) MessageType() byte {
	return wrp.actual.Type
}

// Identifier implements the Message interface
func (wrp *messageWrapper) Identifier() [8]byte {
	return wrp.actual.Identifier
}

// Name implements the Message interface
func (wrp *messageWrapper) Name() []byte {
	return wrp.actual.Name
}

// Payload implements the Message interface
func (wrp *messageWrapper) Payload() Payload {
	return &BufferedEncodedPayload{
		Payload: pld.Payload{
			Encoding: wrp.actual.Payload.Encoding,
			Data:     wrp.actual.Payload.Data,
		},
	}
}

// Close implements the Message interface
func (wrp *messageWrapper) Close() {
	wrp.actual.Buffer.Close()
}
