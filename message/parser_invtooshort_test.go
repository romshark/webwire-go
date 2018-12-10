package message_test

import (
	"testing"

	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/require"
)

/****************************************************************\
	Parser - invalid messages (too short)
\****************************************************************/

// TestMsgParseInvalidMessageTooShort tests parsing of an invalid
// empty message
func TestMsgParseInvalidMessageTooShort(t *testing.T) {
	invalidMessage := make([]byte, 0)

	actual := message.NewMessage(1024)
	typeDetermined, _ := actual.ReadBytes(invalidMessage)
	require.False(t,
		typeDetermined,
		"Expected type to not be determined "+
			"when parsing empty message",
	)
}

// TestMsgParseInvalidReplyTooShort tests parsing of an invalid
// binary/UTF8 reply message which is too short to be considered valid
func TestMsgParseInvalidReplyTooShort(t *testing.T) {
	lenTooShort := message.MinLenReply - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgReplyBinary

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid reply message (too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidReplyUtf16TooShort tests parsing of an invalid
// UTF16 reply message which is too short to be considered valid
func TestMsgParseInvalidReplyUtf16TooShort(t *testing.T) {
	lenTooShort := message.MinLenReplyUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgReplyUtf16

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid UTF16 reply message "+
			"(too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidRequestTooShort tests parsing of an invalid
// binary/UTF8 request message which is too short to be considered valid
func TestMsgParseInvalidRequestTooShort(t *testing.T) {
	lenTooShort := message.MinLenRequest - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgRequestBinary

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid request message (too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidRequestUtf16TooShort tests parsing of an invalid
// UTF16 request message which is too short to be considered valid
func TestMsgParseInvalidRequestUtf16TooShort(t *testing.T) {
	lenTooShort := message.MinLenRequestUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgRequestUtf16

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid UTF16 "+
			"encoded request message (too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidRestrSessReqTooShort tests parsing of an invalid
// session restoration request message which is too short
// to be considered valid
func TestMsgParseInvalidRestrSessReqTooShort(t *testing.T) {
	lenTooShort := message.MinLenRequestRestoreSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgRequestRestoreSession

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid session restoration "+
			"request message (too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidSessCloseReqTooShort tests parsing of an invalid
// session destruction request message which is too short
// to be considered valid
func TestMsgParseInvalidSessCloseReqTooShort(t *testing.T) {
	lenTooShort := message.MinLenDoCloseSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgRequestCloseSession

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid session destruction "+
			"request message (too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidSessCreatedSigTooShort tests parsing of an invalid
// session creation notification message which is too short
// to be considered valid
func TestMsgParseInvalidSessCreatedSigTooShort(t *testing.T) {
	lenTooShort := message.MinLenNotifySessionCreated - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgNotifySessionCreated

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid session creation "+
			"notification message (too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidSignalTooShort tests parsing of an invalid
// binary/UTF8 signal message which is too short to be considered valid
func TestMsgParseInvalidSignalTooShort(t *testing.T) {
	lenTooShort := message.MinLenSignal - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgSignalBinary

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid signal message (too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidSignalUtf16TooShort tests parsing of an invalid
// UTF16 signal message which is too short to be considered valid
func TestMsgParseInvalidSignalUtf16TooShort(t *testing.T) {
	lenTooShort := message.MinLenSignalUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgSignalUtf16

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid UTF16 signal message "+
			"(too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidErrorReplyTooShort tests parsing of an invalid
// error reply message which is too short to be considered valid
func TestMsgParseInvalidErrorReplyTooShort(t *testing.T) {
	lenTooShort := message.MinLenReplyError - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = message.MsgReplyError

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid error reply message "+
			"(too short: %d)",
		lenTooShort,
	)
}

// TestMsgParseInvalidSpecialReplyTooShort tests parsing of an invalid
// special reply message which is too short to be considered valid
func TestMsgParseInvalidSpecialReplyTooShort(t *testing.T) {
	invalidMessage := make([]byte, 8)

	// Internal error is a special reply message type
	invalidMessage[0] = message.MsgReplyInternalError

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid special reply message "+
			"(too short: 8)",
	)
}
