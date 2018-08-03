package message

import (
	"testing"

	pld "github.com/qbeon/webwire-go/payload"
)

/****************************************************************\
	Parser - invalid input (corrupt payload)
\****************************************************************/

// TestMsgParseReplyUtf16CorruptInput tests parsing of
// UTF16 encoded reply message with a corrupted input stream
// (length not divisible by 2)
func TestMsgParseReplyUtf16CorruptInput(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgReplyUtf16}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add header padding byte due to UTF16 encoding
	encoded = append(encoded, byte(0))
	// Add payload
	encoded = append(encoded, payload.Data...)

	// Parse
	if _, err := tryParse(t, encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

// TestMsgParseRequestUtf16CorruptInput tests parsing of a named
// UTF16 encoded request with a corrupted input stream
// (length not divisible by 2)
func TestMsgParseRequestUtf16CorruptInput(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgRequestUtf16}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		encoded = append(encoded, byte(0))
	}
	// Add payload
	encoded = append(encoded, payload.Data...)

	// Parse
	if _, err := tryParse(t, encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

// TestMsgParseSignalUtf16CorruptInput tests parsing of a named
// UTF16 encoded signal with a corrupt unaligned input stream
// (length not divisible by 2)
func TestMsgParseSignalUtf16CorruptInput(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgSignalUtf16}
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		encoded = append(encoded, byte(0))
	}
	// Add payload
	encoded = append(encoded, payload.Data...)

	// Parse
	if _, err := tryParse(t, encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}
