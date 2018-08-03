package message

import "testing"

/****************************************************************\
	Parser - invalid messages (too long)
\****************************************************************/

// TestMsgParseInvalidSessionClosedTooLong tests parsing of an invalid
// session closed notification message which is too long to be considered valid
func TestMsgParseInvalidSessionClosedTooLong(t *testing.T) {
	lenTooLong := MsgMinLenSessionClosed + 1
	invalidMessage := make([]byte, lenTooLong)

	invalidMessage[0] = MsgSessionClosed

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid session closed message "+
				"(too long: %d)",
			lenTooLong,
		)
	}
}
