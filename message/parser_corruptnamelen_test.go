package message_test

import (
	"bytes"
	"testing"

	"github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

/****************************************************************\
	Parser - invalid messages (corrupt name length flags)
\****************************************************************/

// TestMsgParseRequestCorruptNameLenFlag tests parsing of a named
// Binary/UTF8 encoded request with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseRequestCorruptNameLenFlag(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(10 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(message.MsgRequestBinary)
	// Add identifier
	encoded.Write(id[:])

	// Add corrupt name length flag (too big) and skip the name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	_, err := tryParse(t, encoded.Bytes())
	require.Error(t,
		err,
		"Expected Parse to return an error due to corrupt name length flag",
	)
}

// TestMsgParseRequestUtf16CorruptNameLenFlag tests parsing of a named
// UTF16 encoded request with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseRequestUtf16CorruptNameLenFlag(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(10 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(message.MsgRequestUtf16)
	// Add identifier
	encoded.Write(id[:])

	// Add corrupt name length flag (too big) and skip actual name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	_, err := tryParse(t, encoded.Bytes())
	require.Error(t,
		err,
		"Expected Parse to return an error due to corrupt name length flag",
	)
}

// TestMsgParseSignalCorruptNameLenFlag tests parsing of a named
// Binary/UTF8 encoded signal with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseSignalCorruptNameLenFlag(t *testing.T) {
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(2 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(message.MsgSignalBinary)

	// Add corrupt name length flag (too big) and skip the name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	_, err := tryParse(t, encoded.Bytes())
	require.Error(t,
		err,
		"Expected Parse to return an error due to corrupt name length flag",
	)
}

// TestMsgParseSignalUtf16CorruptNameLenFlag tests parsing of a named
// UTF16 encoded signal with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseSignalUtf16CorruptNameLenFlag(t *testing.T) {
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(2 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(message.MsgSignalUtf16)

	// Add corrupt name length flag (too big) and skip the name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	_, err := tryParse(t, encoded.Bytes())
	require.Error(t,
		err,
		"Expected Parse to return an error due to corrupt name length flag",
	)
}
