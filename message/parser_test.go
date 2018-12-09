package message

import (
	"encoding/json"
	"testing"
	"time"

	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

/****************************************************************\
	Parser
\****************************************************************/

// TestMsgParseCloseSessReq tests parsing of a session destruction request
func TestMsgParseCloseSessReq(t *testing.T) {
	id := genRndMsgIdentifier()

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgDoCloseSession}
	// Add identifier
	encoded = append(encoded, id[:]...)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgDoCloseSession, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{
		Encoding: pld.Binary,
		Data:     nil,
	}, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseRestrSessReq tests parsing of a session restoration request
func TestMsgParseRestrSessReq(t *testing.T) {
	id := genRndMsgIdentifier()

	//sessionKey := sess.GenerateSessionKey()
	sessionKey := "somesamplesessionkey"

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgRequestRestoreSession}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add session key to payload
	encoded = append(encoded, sessionKey[:]...)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgRequestRestoreSession, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte(sessionKey),
	}, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseRequestBinary tests parsing of a named binary encoded request
func TestMsgParseRequestBinary(t *testing.T) {
	encoded, id, name, payload := rndRequestMsg(
		MsgRequestBinary,
		1, 255,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgRequestBinary, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseRequestUtf8 tests parsing of a named UTF8 encoded request
func TestMsgParseRequestUtf8(t *testing.T) {
	encoded, id, name, payload := rndRequestMsg(
		MsgRequestUtf8,
		2, 255,
		16, 16,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgRequestUtf8, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseRequestUtf16 tests parsing of a named UTF16 encoded request
func TestMsgParseRequestUtf16(t *testing.T) {
	encoded, id, name, payload := rndRequestMsgUtf16(
		1, 255,
		2, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgRequestUtf16, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseReplyBinary tests parsing of binary encoded reply message
func TestMsgParseReplyBinary(t *testing.T) {
	encoded, id, payload := rndReplyMsg(
		MsgReplyBinary,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgReplyBinary, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseReplyUtf8 tests parsing of UTF8 encoded reply message
func TestMsgParseReplyUtf8(t *testing.T) {
	encoded, id, payload := rndReplyMsg(
		MsgReplyUtf8,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgReplyUtf8, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseReplyUtf16 tests parsing of UTF16 encoded reply message
func TestMsgParseReplyUtf16(t *testing.T) {
	encoded, id, payload := rndReplyMsgUtf16(
		2, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgReplyUtf16, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSignalBinary tests parsing of a named binary encoded signal
func TestMsgParseSignalBinary(t *testing.T) {
	encoded, name, payload := rndSignalMsg(
		MsgSignalBinary,
		1, 255,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgSignalBinary, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSignalUtf8 tests parsing of a named UTF8 encoded signal
func TestMsgParseSignalUtf8(t *testing.T) {
	encoded, name, payload := rndSignalMsg(
		MsgSignalUtf8,
		1, 255,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgSignalUtf8, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSignalUtf16 tests parsing of a named UTF16 encoded signal
func TestMsgParseSignalUtf16(t *testing.T) {
	encoded, name, payload := rndSignalMsgUtf16(
		1, 255,
		2, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgSignalUtf16, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSessCreatedSig tests parsing of session created signal
func TestMsgParseSessCreatedSig(t *testing.T) {
	//sessionKey := generateSessionKey()
	sessionKey := "somesamplesessionkey"
	session := struct {
		Key      string
		Creation time.Time
		Info     interface{}
	}{
		Key:      sessionKey,
		Creation: time.Now(),
		Info:     nil,
	}
	marshalledSession, err := json.Marshal(&session)
	require.NoError(t, err)
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     marshalledSession,
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgNotifySessionCreated}
	// Add session payload
	encoded = append(encoded, payload.Data...)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgNotifySessionCreated, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSessClosedSig tests parsing of session closed signal
func TestMsgParseSessClosedSig(t *testing.T) {
	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgNotifySessionClosed}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgNotifySessionClosed, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{}, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseHeartbeat tests parsing of heartbeat messages
func TestMsgParseHeartbeat(t *testing.T) {
	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgHeartbeat}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgHeartbeat, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{}, actual.MsgPayload)
	require.Equal(t, ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseConf tests parsing of server configuration messages
func TestMsgParseConf(t *testing.T) {
	srvConf := ServerConfiguration{
		MajorProtocolVersion: 22,
		MinorProtocolVersion: 33,
		ReadTimeout:          11 * time.Second,
		MessageBufferSize:    8192,
	}

	// Compose encoded message
	buf, err := NewConfMessage(srvConf)
	require.NoError(t, err)
	require.True(t, len(buf) > 0)

	// Parse
	actual := tryParseNoErr(t, buf)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, MsgAcceptConf, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{}, actual.MsgPayload)
	require.Equal(t, srvConf, actual.ServerConfiguration)
}

// TestMsgParseUnknownMessageType tests parsing of messages
// with unknown message type
func TestMsgParseUnknownMessageType(t *testing.T) {
	msgOfUnknownType := make([]byte, 1)
	msgOfUnknownType[0] = byte(255)

	actual := NewMessage(1024)
	typeDetermined, _ := actual.ReadBytes(msgOfUnknownType)
	require.False(t, typeDetermined, "Expected type not to be determined")
}
