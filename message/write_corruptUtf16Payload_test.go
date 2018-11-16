package message

import (
	"testing"

	"github.com/stretchr/testify/require"

	pld "github.com/qbeon/webwire-go/payload"
)

/****************************************************************\
	Constructors - invalid input (corrupt name length flags)
\****************************************************************/

// TestWriteMsgReplyUtf16CorruptPayload tests WriteMsgReply using UTF16 payload
// encoding passing corrupt data (length not divisible by 2 thus not UTF16
// encoded)
func TestWriteMsgReplyUtf16CorruptPayload(t *testing.T) {
	writer := &testWriter{}
	require.Error(t, WriteMsgReply(
		writer,
		genRndMsgIdentifier(),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
	))
	require.True(t, writer.closed)
}

// TestWriteMsgReqUtf16CorruptPayload tests WriteMsgRequest using UTF16 payload
// encoding passing corrupt data (length not divisible by 2 thus not UTF16
// encoded)
func TestWriteMsgReqUtf16CorruptPayload(t *testing.T) {
	writer := &testWriter{}
	require.Error(t, WriteMsgRequest(
		writer,
		genRndMsgIdentifier(),
		genRndName(1, 255),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
		true,
	))
	require.True(t, writer.closed)
}

// TestWriteMsgSigUtf16CorruptPayload tests WriteMsgSignal using UTF16
// payload encoding passing corrupt data (length not divisible by 2 thus not
// UTF16 encoded)
func TestWriteMsgSigUtf16CorruptPayload(t *testing.T) {
	writer := &testWriter{}
	require.Error(t, WriteMsgSignal(
		writer,
		genRndName(1, 255),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
		true,
	))
	require.True(t, writer.closed)
}
