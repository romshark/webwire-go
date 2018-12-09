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

	invalidMessage[0] = MsgNotifySessionClosed

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid session closed message "+
			"(too long: %d)",
		lenTooLong,
	)
}

// TestMsgParseInvalidHeartbeatTooLong tests parsing of an invalid heartbeat
// message which is too long to be considered valid
func TestMsgParseInvalidHeartbeatTooLong(t *testing.T) {
	lenTooLong := 2
	invalidMessage := make([]byte, lenTooLong)

	invalidMessage[0] = MsgHeartbeat

	_, err := tryParse(t, invalidMessage)
	require.Error(t,
		err,
		"Expected error while parsing invalid heartbeat message "+
			"(too long: %d)",
		lenTooLong,
	)
}
