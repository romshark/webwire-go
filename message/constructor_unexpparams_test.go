package message

import (
	"testing"

	"github.com/stretchr/testify/require"

	pld "github.com/qbeon/webwire-go/payload"
)

/****************************************************************\
	Constructors - unexpected parameters (panics)
\****************************************************************/

// TestMsgNewReqMsgNoNameNoPayload tests calling
// the request message constructor without both the name and the payload
func TestMsgNewReqMsgNoNameNoPayload(t *testing.T) {
	require.Panics(t,
		func() {
			NewRequestMessage(
				genRndMsgIdentifier(),
				nil,
				pld.Binary,
				nil,
			)
		},
		"Expected a panic after calling the "+
			" request message constructor without both the name "+
			"and the payload",
	)
}

// TestMsgNewReqMsgNameTooLong tests NewRequestMessage with a too long name
func TestMsgNewReqMsgNameTooLong(t *testing.T) {
	tooLongNamelength := 256
	nameBuf := make([]byte, tooLongNamelength)
	for i := 0; i < tooLongNamelength; i++ {
		nameBuf[i] = 'a'
	}

	require.Panics(t,
		func() {
			NewRequestMessage(
				genRndMsgIdentifier(),
				nameBuf,
				0,
				nil,
			)
		},
		"Expected panic after passing a too long request name",
	)
}

// TestMsgNewReqMsgInvalidCharsetBelowAscii32 tests NewRequestMessage
// with an invalid character input below the ASCII 7 bit 32nd character
func TestMsgNewReqMsgInvalidCharsetBelowAscii32(t *testing.T) {
	// Generate invalid name using a character
	// below the ASCII 7 bit 32nd character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(31)

	require.Panics(t,
		func() {
			NewRequestMessage(
				genRndMsgIdentifier(),
				invalidNameBytes,
				0,
				nil,
			)
		},
		"Expected panic after passing an invalid name character set",
	)
}

// TestMsgNewReqMsgInvalidCharsetAboveAscii126 tests NewRequestMessage
// with an invalid character input above the ASCII 7 bit 126th character
func TestMsgNewReqMsgInvalidCharsetAboveAscii126(t *testing.T) {
	// Generate invalid name using a character
	// above the ASCII 7 bit 126th character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(127)

	require.Panics(t,
		func() {
			NewRequestMessage(
				genRndMsgIdentifier(),
				invalidNameBytes,
				0,
				nil,
			)
		},
		"Expected panic after passing an invalid name character set",
	)
}

// TestMsgNewSigMsgNameTooLong tests NewSignalMessage with a too long name
func TestMsgNewSigMsgNameTooLong(t *testing.T) {
	tooLongNamelength := 256
	nameBuf := make([]byte, tooLongNamelength)
	for i := 0; i < tooLongNamelength; i++ {
		nameBuf[i] = 'a'
	}

	require.Panics(t,
		func() {
			NewSignalMessage(nameBuf, 0, nil)
		},
		"Expected panic after passing a too long signal name",
	)
}

// TestMsgNewSigMsgInvalidCharsetBelowAscii32 tests NewSignalMessage
// with an invalid character input below the ASCII 7 bit 32nd character
func TestMsgNewSigMsgInvalidCharsetBelowAscii32(t *testing.T) {
	// Generate invalid name using a character
	// below the ASCII 7 bit 32nd character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(31)

	require.Panics(t,
		func() {
			NewSignalMessage(invalidNameBytes, 0, nil)
		},
		"Expected panic after passing an invalid name character set",
	)
}

// TestMsgNewSigMsgInvalidCharsetAboveAscii126 tests NewSignalMessage
// with an invalid character input above ASCII 7 bit 126th character
func TestMsgNewSigMsgInvalidCharsetAboveAscii126(t *testing.T) {
	// Generate invalid name using a character
	// above the ASCII 7 bit 126th character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(127)

	require.Panics(t,
		func() {
			NewSignalMessage(invalidNameBytes, 0, nil)
		},
		"Expected panic after passing an invalid name character set",
	)
}

// TestMsgNewSpecialRequestReplyMessageInvalidType tests
// NewSpecialRequestReplyMessage with non-special reply message types
func TestMsgNewSpecialRequestReplyMessageInvalidType(t *testing.T) {
	allTypes := []byte{
		MsgErrorReply,
		MsgSessionCreated,
		MsgSessionClosed,
		MsgCloseSession,
		MsgRestoreSession,
		MsgSignalBinary,
		MsgSignalUtf8,
		MsgSignalUtf16,
		MsgRequestBinary,
		MsgRequestUtf8,
		MsgRequestUtf16,
		MsgReplyBinary,
		MsgReplyUtf8,
		MsgReplyUtf16,
	}

	for _, tp := range allTypes {
		require.Panics(t,
			func() {
				NewSpecialRequestReplyMessage(
					tp,
					genRndMsgIdentifier(),
				)
			},
			"Expected panic after passing "+
				"a non-special request reply message type",
		)
	}
}

// TestMsgNewErrorReplyMessageNoCode tests NewErrorReplyMessage
// with no error code which is invalid.
func TestMsgNewErrorReplyMessageNoCode(t *testing.T) {

	require.Panics(t,
		func() {
			NewErrorReplyMessage(
				genRndMsgIdentifier(),
				"",
				"sample error message",
			)
		},
		"Expected panic when creating an error reply message "+
			"with no error code ",
	)
}

// TestMsgNewErrorReplyMessageCodeTooLong tests NewErrorReplyMessage
// with an error code that's surpassing the 255 character limit.
func TestMsgNewErrorReplyMessageCodeTooLong(t *testing.T) {
	tooLongCode := make([]byte, 256)
	for i := 0; i < 256; i++ {
		tooLongCode[i] = 'a'
	}

	require.Panics(t,
		func() {
			NewErrorReplyMessage(
				genRndMsgIdentifier(),
				string(tooLongCode),
				"sample error message",
			)
		},
		"Expected panic when creating an error reply message "+
			"with no error code ",
	)
}

// TestMsgNewErrorReplyMessageCodeCharsetBelowAscii32 tests NewErrorReplyMessage
// with an invalid character input below the ASCII 7 bit 32nd character
func TestMsgNewErrorReplyMessageCodeCharsetBelowAscii32(t *testing.T) {
	// Generate invalid error code using a character
	// below the ASCII 7 bit 32nd character
	invalidCodeBytes := make([]byte, 1)
	invalidCodeBytes[0] = byte(31)

	require.Panics(t,
		func() {
			NewErrorReplyMessage(
				genRndMsgIdentifier(),
				string(invalidCodeBytes),
				"sample error message",
			)
		},
		"Expected panic after passing an invalid error code "+
			" containing a character below the 32th ASCII 7bit char",
	)
}

// TestMsgNewErrorReplyMessageCodeCharsetAboveAscii126 tests
// NewErrorReplyMessage with an invalid character input
// above the ASCII 7 bit 126th character
func TestMsgNewErrorReplyMessageCodeCharsetAboveAscii126(t *testing.T) {
	// Generate invalid error code using a character
	// above the ASCII 7 bit 126th character
	invalidCodeBytes := make([]byte, 1)
	invalidCodeBytes[0] = byte(127)

	require.Panics(t,
		func() {
			NewErrorReplyMessage(
				genRndMsgIdentifier(),
				string(invalidCodeBytes),
				"sample error message",
			)
		},
		"Expected panic after passing an invalid error code "+
			" containing a character above the 126th ASCII 7bit char",
	)
}
