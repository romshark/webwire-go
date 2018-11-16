package client

import (
	"github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
)

// determineMsgTypeBasedOnEncoding returns the corresponding message type for
// the given encoding type.
func determineMsgTypeBasedOnEncoding(encoding webwire.PayloadEncoding) byte {
	// Set the payload type if any payload is given,
	// otherwise fallback to binary
	switch encoding {
	case webwire.EncodingUtf8:
		return message.MsgRequestUtf8
	case webwire.EncodingUtf16:
		return message.MsgRequestUtf16
	default:
		return message.MsgRequestBinary
	}
}
