package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// NewReplyMessage composes a new reply message and returns its binary representation
func NewReplyMessage(
	requestIdentifier [8]byte,
	payloadEncoding pld.Encoding,
	payloadData []byte,
) (msg []byte) {
	// Determine total message length
	messageSize := 9 + len(payloadData)

	// Verify payload data validity in case of UTF16 encoding
	if payloadEncoding == pld.Utf16 && len(payloadData)%2 != 0 {
		panic(fmt.Errorf("Invalid UTF16 reply payload data length: %d", len(payloadData)))
	}

	// Check if a header padding is necessary.
	// A padding is necessary if the payload is UTF16 encoded
	// but not properly aligned due to a header length not divisible by 2
	headerPadding := false
	if payloadEncoding == pld.Utf16 {
		headerPadding = true
		messageSize++
	}

	msg = make([]byte, messageSize)

	// Write message type flag
	reqType := MsgReplyBinary
	switch payloadEncoding {
	case pld.Utf8:
		reqType = MsgReplyUtf8
	case pld.Utf16:
		reqType = MsgReplyUtf16
	}
	msg[0] = reqType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = requestIdentifier[i]
	}

	// Write header padding byte if the payload requires proper alignment
	payloadOffset := 9
	if headerPadding {
		msg[payloadOffset] = 0
		payloadOffset++
	}

	// Write payload
	for i := 0; i < len(payloadData); i++ {
		msg[payloadOffset+i] = payloadData[i]
	}

	return msg
}
