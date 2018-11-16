package webwire

import "github.com/qbeon/webwire-go/payload"

// PayloadEncoding represents the type of encoding of the message payload
type PayloadEncoding = payload.Encoding

const (
	// EncodingBinary represents unencoded binary data
	EncodingBinary PayloadEncoding = payload.Binary

	// EncodingUtf8 represents UTF8 encoding
	EncodingUtf8 PayloadEncoding = payload.Utf8

	// EncodingUtf16 represents UTF16 encoding
	EncodingUtf16 PayloadEncoding = payload.Utf16
)

// Payload represents an encoded payload
type Payload struct {
	// Encoding represents the encoding type of the payload which is
	// EncodingBinary by default
	Encoding PayloadEncoding

	// Data represents the payload data
	Data []byte
}
