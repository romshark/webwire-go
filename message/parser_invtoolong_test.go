package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/****************************************************************\
	Parser - invalid messages (too long)
\****************************************************************/

// TestMsgParseInvalidSessionClosedTooLong tests parsing of an invalid
// session closed notification message which is too long to be considered valid
func TestMsgParseInvalidSessionClosedTooLong(t *testing.T) {
	lenTooLong := MsgMinLenSessionClosed + 1
	invalidMessage := make([]byte, lenTooLong)

	invalidMessage[0] = MsgSessionClosed

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid session closed message "+
			"(too long: %d)",
		lenTooLong,
	)
}
