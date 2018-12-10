package message_test

import (
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestMsgParseAcceptNoSubprotoConf tests parsing of server configuration
// messages with no subprotocol name
func TestMsgParseAcceptNoSubprotoConf(t *testing.T) {
	srvConf := message.ServerConfiguration{
		MajorProtocolVersion: 22,
		MinorProtocolVersion: 33,
		ReadTimeout:          11 * time.Second,
		MessageBufferSize:    8192,
		SubprotocolName:      []byte(nil),
	}

	// Compose encoded message
	buf, err := message.NewAcceptConfMessage(srvConf)
	require.NoError(t, err)
	require.True(t, len(buf) > 0)

	// Parse
	actual := tryParseNoErr(t, buf)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgAcceptConf, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{}, actual.MsgPayload)
	require.Equal(t, srvConf, actual.ServerConfiguration)
}

// TestMsgParseAcceptConf tests parsing of server configuration messages
func TestMsgParseAcceptConf(t *testing.T) {
	srvConf := message.ServerConfiguration{
		MajorProtocolVersion: 22,
		MinorProtocolVersion: 33,
		ReadTimeout:          11 * time.Second,
		MessageBufferSize:    8192,
		SubprotocolName:      []byte("test - subprotocol name"),
	}

	// Compose encoded message
	buf, err := message.NewAcceptConfMessage(srvConf)
	require.NoError(t, err)
	require.True(t, len(buf) > 0)

	// Parse
	actual := tryParseNoErr(t, buf)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgAcceptConf, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{}, actual.MsgPayload)
	require.Equal(t, srvConf, actual.ServerConfiguration)
}
