package webwire

import (
	"errors"

	"github.com/qbeon/webwire-go/msgbuf"
	pld "github.com/qbeon/webwire-go/payload"
)

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

// BufferedEncodedPayload represents an implementation of the webwire.Payload
// interface
type BufferedEncodedPayload struct {
	Buffer  *msgbuf.MessageBuffer
	Payload pld.Payload
	Closed  bool
}

// Encoding implements the webwire.Payload interface
func (pld *BufferedEncodedPayload) Encoding() PayloadEncoding {
	if pld.Closed {
		panic("payload read after close")
	}
	return PayloadEncoding(pld.Payload.Encoding)
}

// Data implements the webwire.Payload interface
func (pld *BufferedEncodedPayload) Data() []byte {
	if pld.Closed {
		panic("payload read after close")
	}
	return pld.Payload.Data
}

// Utf8 implements the webwire.Payload interface
func (pld *BufferedEncodedPayload) Utf8() (string, error) {
	if pld.Closed {
		return "", errors.New("payload read after close")
	}
	return pld.Payload.Utf8()
}

// Close implements the webwire.Payload interface
func (pld *BufferedEncodedPayload) Close() {
	if !pld.Closed {
		pld.Buffer.Close()
		pld.Payload.Encoding = 0
		pld.Payload.Data = nil
		pld.Closed = true
	}
}

// NewPayload creates a new WebWire message payload
func NewPayload(encoding PayloadEncoding, data []byte) Payload {
	return &BufferedEncodedPayload{
		Payload: pld.Payload{
			Encoding: encoding,
			Data:     data,
		},
	}
}
