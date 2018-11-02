package message

import (
	"testing"

	"github.com/stretchr/testify/require"

	pld "github.com/qbeon/webwire-go/payload"
)

/****************************************************************\
	Constructors - invalid input (corrupt name length flags)
\****************************************************************/

// TestMsgNewReplyMsgUtf16CorruptPayload tests NewReplyMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReplyMsgUtf16CorruptPayload(t *testing.T) {
	require.Panics(t, func() {
		NewReplyMessage(
			genRndMsgIdentifier(),
			pld.Utf16,
			// Payload is corrupt, only 7 bytes long, not power 2
			[]byte("invalid"),
		)
	})
}

// TestMsgNewReqMsgUtf16CorruptPayload tests NewRequestMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReqMsgUtf16CorruptPayload(t *testing.T) {
	require.Panics(t, func() {
		NewRequestMessage(
			genRndMsgIdentifier(),
			genRndName(1, 255),
			pld.Utf16,
			// Payload is corrupt, only 7 bytes long, not power 2
			[]byte("invalid"),
		)
	})
}

// TestMsgNewSigMsgUtf16CorruptPayload tests NewSignalMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewSigMsgUtf16CorruptPayload(t *testing.T) {
	require.Panics(t, func() {
		NewSignalMessage(
			genRndName(1, 255),
			pld.Utf16,
			// Payload is corrupt, only 7 bytes long, not power 2
			[]byte("invalid"),
		)
	})
}
