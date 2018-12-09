package requestmanager

import (
	"github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
)

// reply represents an implementation of the Reply interface
type reply struct {
	msg *message.Message
}

// PayloadEncoding implements the Reply interface
func (rp *reply) PayloadEncoding() webwire.PayloadEncoding {
	return rp.msg.MsgPayload.Encoding
}

// Payload implements the Reply interface
func (rp *reply) Payload() []byte {
	return rp.msg.MsgPayload.Data
}

// PayloadUtf8 implements the Reply interface
func (rp *reply) PayloadUtf8() ([]byte, error) {
	return rp.msg.MsgPayload.Utf8()
}

// Close implements the Reply interface
func (rp *reply) Close() {
	rp.msg.Close()
}
