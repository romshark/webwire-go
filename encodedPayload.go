package webwire

import pld "github.com/qbeon/webwire-go/payload"

// PayloadEncoding represents the type of encoding of the message payload
type PayloadEncoding = pld.Encoding

const (
	// EncodingBinary represents unencoded binary data
	EncodingBinary PayloadEncoding = pld.Binary

	// EncodingUtf8 represents UTF8 encoding
	EncodingUtf8 = pld.Utf8

	// EncodingUtf16 represents UTF16 encoding
	EncodingUtf16 = pld.Utf16
)

// EncodedPayload represents an encoded message payload
// and implements the WebWire payload interface
type EncodedPayload struct {
	Payload pld.Payload
}

// Encoding implements the WebWire payload interface
func (pld *EncodedPayload) Encoding() PayloadEncoding {
	return PayloadEncoding(pld.Payload.Encoding)
}

// Data implements the WebWire payload interface
func (pld *EncodedPayload) Data() []byte {
	return pld.Payload.Data
}

// Utf8 implements the WebWire payload interface
func (pld *EncodedPayload) Utf8() (string, error) {
	return pld.Payload.Utf8()
}

// NewPayload creates a new WebWire message payload
func NewPayload(encoding PayloadEncoding, data []byte) Payload {
	return &EncodedPayload{
		Payload: pld.Payload{
			Encoding: encoding,
			Data:     data,
		},
	}
}
