package client

import (
	webwire "github.com/qbeon/webwire-go"
	msg "github.com/qbeon/webwire-go/message"
)

// scanPayloadEncoding tries to find out the encoding type
// of the given payload and return the corresponding message type for it.
// If the payload is missing then it will return binary by default
func scanPayloadEncoding(payload webwire.Payload) byte {
	// Set the payload type if any payload is given,
	// otherwise fallback to binary
	reqType := msg.MsgRequestBinary
	if payload != nil {
		switch payload.Encoding() {
		case webwire.EncodingUtf8:
			reqType = msg.MsgRequestUtf8
		case webwire.EncodingUtf16:
			reqType = msg.MsgRequestUtf16
		}
	}
	return reqType
}
