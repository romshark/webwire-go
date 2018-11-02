package message

import (
	"testing"

	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestRequiresReplyCloseSession tests the RequiresReply method
// with a session closure request message
func TestRequiresReplyCloseSession(t *testing.T) {
	msg := &Message{}
	_, err := msg.Parse(NewEmptyRequestMessage(
		MsgCloseSession,
		genRndMsgIdentifier(),
	))
	require.NoError(t, err)

	require.True(t,
		msg.RequiresReply(),
		"Expected a session closure request message to require a reply",
	)
}

// TestRequiresReplyRestoreSession tests the RequiresReply method
// with a session restoration request message
func TestRequiresReplyRestoreSession(t *testing.T) {
	msg := &Message{}
	_, err := msg.Parse(NewNamelessRequestMessage(
		MsgRestoreSession,
		genRndMsgIdentifier(),
		[]byte("somesamplesessionkey"),
	))
	require.NoError(t, err)

	require.True(t,
		msg.RequiresReply(),
		"Expected a session restoration request message to require a reply",
	)
}

// TestRequiresReplyRequestBinary tests the RequiresReply method
// with a binary request message
func TestRequiresReplyRequestBinary(t *testing.T) {
	msg := &Message{}
	_, err := msg.Parse(NewRequestMessage(
		genRndMsgIdentifier(),
		[]byte("samplename"),
		pld.Binary,
		[]byte("random payload data"),
	))
	require.NoError(t, err)

	require.True(t,
		msg.RequiresReply(),
		"Expected a binary request message to require a reply",
	)
}

// TestRequiresReplyRequestUtf8 tests the RequiresReply method
// with a UTF8 encoded request message
func TestRequiresReplyRequestUtf8(t *testing.T) {
	msg := &Message{}
	_, err := msg.Parse(NewRequestMessage(
		genRndMsgIdentifier(),
		[]byte("samplename"),
		pld.Utf8,
		[]byte("random payload data"),
	))
	require.NoError(t, err)

	require.True(t,
		msg.RequiresReply(),
		"Expected a UTF8 request message to require a reply",
	)
}

// TestRequiresReplyRequestUtf16 tests the RequiresReply method
// with a UTF16 encoded request message
func TestRequiresReplyRequestUtf16(t *testing.T) {
	msg := &Message{}
	_, err := msg.Parse(NewRequestMessage(
		genRndMsgIdentifier(),
		[]byte("samplename"),
		pld.Utf16,
		[]byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	))
	require.NoError(t, err)

	require.True(t,
		msg.RequiresReply(),
		"Expected a UTF16 request message to require a reply",
	)
}
