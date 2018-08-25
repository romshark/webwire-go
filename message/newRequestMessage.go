package message

import (
	"fmt"

	pld "github.com/qbeon/webwire-go/payload"
)

// NewRequestMessage composes a new named request message
// and returns its binary representation
func NewRequestMessage(
	identifier [8]byte,
	name string,
	payloadEncoding pld.Encoding,
	payloadData []byte,
) (msg []byte) {
	// Require either a name, or a payload or both, but don't allow none
	if len(name) < 1 && len(payloadData) < 1 {
		panic(fmt.Errorf(
			"Request message requires either a name, or a payload, or both",
		))
	}

	// Cap name length at 255 bytes
	if len(name) > 255 {
		panic(fmt.Errorf(
			"Unsupported request message name length: %d",
			len(name),
		))
	}

	// Verify payload data validity in case of UTF16 encoding
	if payloadEncoding == pld.Utf16 && len(payloadData)%2 != 0 {
		panic(fmt.Errorf(
			"Invalid UTF16 request payload data length: %d",
			len(payloadData),
		))
	}

	// Determine total message length
	messageSize := 10 + len(name) + len(payloadData)

	// Check if a header padding is necessary.
	// A padding is necessary if the payload is UTF16 encoded
	// but not properly aligned due to a header length not divisible by 2
	headerPadding := false
	if payloadEncoding == pld.Utf16 && len(name)%2 != 0 {
		headerPadding = true
		messageSize++
	}

	msg = make([]byte, messageSize)

	// Write message type flag
	reqType := MsgRequestBinary
	switch payloadEncoding {
	case pld.Utf8:
		reqType = MsgRequestUtf8
	case pld.Utf16:
		reqType = MsgRequestUtf16
	}
	msg[0] = reqType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = identifier[i]
	}

	// Write name length flag
	msg[9] = byte(len(name))

	// Write name
	for i := 0; i < len(name); i++ {
		char := name[i]
		if char < 32 || char > 126 {
			panic(fmt.Errorf(
				"Unsupported character in request name: %s",
				string(char),
			))
		}
		msg[10+i] = char
	}

	// Write header padding byte if the payload requires proper alignment
	payloadOffset := 10 + len(name)
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
