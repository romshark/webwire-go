package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/****************************************************************\
	Parser - invalid messages (too short)
\****************************************************************/

// TestMsgParseInvalidMessageTooShort tests parsing of an invalid
// empty message
func TestMsgParseInvalidMessageTooShort(t *testing.T) {
	invalidMessage := make([]byte, 0)

	actual := NewMessage(1024)
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
	lenTooShort := MinLenReply - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgReplyBinary

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
	lenTooShort := MinLenReplyUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgReplyUtf16

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
	lenTooShort := MinLenRequest - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestBinary

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
	lenTooShort := MinLenRequestUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestUtf16

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
	lenTooShort := MinLenRequestRestoreSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestRestoreSession

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
	lenTooShort := MinLenDoCloseSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgDoCloseSession

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
	lenTooShort := MinLenNotifySessionCreated - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgNotifySessionCreated

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
	lenTooShort := MinLenSignal - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSignalBinary

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
	lenTooShort := MinLenSignalUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSignalUtf16

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
	lenTooShort := MinLenReplyError - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgReplyError

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
	invalidMessage[0] = MsgReplyInternalError

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid special reply message "+
			"(too short: 8)",
	)
}
