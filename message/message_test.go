package message_test

import (
	"testing"

	"github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestRequiresReplyCloseSession tests the RequiresReply method
// with a session closure request message
func TestRequiresReplyCloseSession(t *testing.T) {
	writer := &testWriter{}
	require.NoError(t, message.WriteMsgNamelessRequest(
		writer,
		message.MsgRequestCloseSession,
		genRndMsgIdentifier(),
		[]byte{},
	))
	require.True(t, writer.closed)

	msg := message.NewMessage(2048)
	typeParsed, err := msg.ReadBytes(writer.buf)
	require.NoError(t, err)
	require.True(t, typeParsed)

	require.Equal(t, message.MsgRequestCloseSession, msg.MsgType)
	require.True(t,
		msg.RequiresReply(),
		"Expected a session closure request message to require a reply",
	)
}

// TestRequiresReplyRestoreSession tests the RequiresReply method
// with a session restoration request message
func TestRequiresReplyRestoreSession(t *testing.T) {
	writer := &testWriter{}
	require.NoError(t, message.WriteMsgNamelessRequest(
		writer,
		message.MsgRequestRestoreSession,
		genRndMsgIdentifier(),
		[]byte("somesamplesessionkey"),
	))
	require.True(t, writer.closed)

	msg := message.NewMessage(1024)
	typeParsed, err := msg.ReadBytes(writer.buf)
	require.NoError(t, err)
	require.True(t, typeParsed)

	require.True(t,
		msg.RequiresReply(),
		"Expected a session restoration request message to require a reply",
	)
}

// TestRequiresReplyRequestBinary tests the RequiresReply method
// with a binary request message
func TestRequiresReplyRequestBinary(t *testing.T) {
	writer := &testWriter{}
	require.NoError(t, message.WriteMsgRequest(
		writer,
		genRndMsgIdentifier(),
		[]byte("samplename"),
		pld.Binary,
		[]byte("random payload data"),
		true,
	))
	require.True(t, writer.closed)

	msg := message.NewMessage(1024)
	typeParsed, err := msg.ReadBytes(writer.buf)
	require.NoError(t, err)
	require.True(t, typeParsed)

	require.True(t,
		msg.RequiresReply(),
		"Expected a binary request message to require a reply",
	)
}

// TestRequiresReplyRequestUtf8 tests the RequiresReply method
// with a UTF8 encoded request message
func TestRequiresReplyRequestUtf8(t *testing.T) {
	writer := &testWriter{}
	require.NoError(t, message.WriteMsgRequest(
		writer,
		genRndMsgIdentifier(),
		[]byte("samplename"),
		pld.Utf8,
		[]byte("random payload data"),
		true,
	))
	require.True(t, writer.closed)

	msg := message.NewMessage(1024)
	typeParsed, err := msg.ReadBytes(writer.buf)
	require.NoError(t, err)
	require.True(t, typeParsed)

	require.True(t,
		msg.RequiresReply(),
		"Expected a UTF8 request message to require a reply",
	)
}

// TestRequiresReplyRequestUtf16 tests the RequiresReply method
// with a UTF16 encoded request message
func TestRequiresReplyRequestUtf16(t *testing.T) {
	writer := &testWriter{}
	require.NoError(t, message.WriteMsgRequest(
		writer,
		genRndMsgIdentifier(),
		[]byte("samplename"),
		pld.Utf16,
		[]byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
		true,
	))
	require.True(t, writer.closed)

	msg := message.NewMessage(1024)
	typeParsed, err := msg.ReadBytes(writer.buf)
	require.NoError(t, err)
	require.True(t, typeParsed)

	require.True(t,
		msg.RequiresReply(),
		"Expected a UTF16 request message to require a reply",
	)
}
