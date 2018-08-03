package message

import (
	"testing"

	pld "github.com/qbeon/webwire-go/payload"
)

/****************************************************************\
	Constructors - invalid input (corrupt name length flags)
\****************************************************************/

// TestMsgNewReplyMsgUtf16CorruptPayload tests NewReplyMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReplyMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewReplyMessage(
		genRndMsgIdentifier(),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
	)
}

// TestMsgNewReqMsgUtf16CorruptPayload tests NewRequestMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReqMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewRequestMessage(
		genRndMsgIdentifier(),
		string(genRndName(1, 255)),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
	)
}

// TestMsgNewSigMsgUtf16CorruptPayload tests NewSignalMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewSigMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewSignalMessage(
		string(genRndName(1, 255)),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
	)
}
