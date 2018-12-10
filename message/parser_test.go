package message_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"
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
	encoded := []byte{message.MsgRequestCloseSession}
	// Add identifier
	encoded = append(encoded, id[:]...)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgRequestCloseSession, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{
		Encoding: pld.Binary,
		Data:     nil,
	}, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseRestrSessReq tests parsing of a session restoration request
func TestMsgParseRestrSessReq(t *testing.T) {
	id := genRndMsgIdentifier()

	//sessionKey := sess.GenerateSessionKey()
	sessionKey := "somesamplesessionkey"

	// Compose encoded message
	// Add type flag
	encoded := []byte{message.MsgRequestRestoreSession}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add session key to payload
	encoded = append(encoded, sessionKey[:]...)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgRequestRestoreSession, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte(sessionKey),
	}, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseRequestBinary tests parsing of a named binary encoded request
func TestMsgParseRequestBinary(t *testing.T) {
	encoded, id, name, payload := rndRequestMsg(
		message.MsgRequestBinary,
		1, 255,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgRequestBinary, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseRequestUtf8 tests parsing of a named UTF8 encoded request
func TestMsgParseRequestUtf8(t *testing.T) {
	encoded, id, name, payload := rndRequestMsg(
		message.MsgRequestUtf8,
		2, 255,
		16, 16,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgRequestUtf8, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
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
	require.Equal(t, message.MsgRequestUtf16, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseReplyBinary tests parsing of binary encoded reply message
func TestMsgParseReplyBinary(t *testing.T) {
	encoded, id, payload := rndReplyMsg(
		message.MsgReplyBinary,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgReplyBinary, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseReplyUtf8 tests parsing of UTF8 encoded reply message
func TestMsgParseReplyUtf8(t *testing.T) {
	encoded, id, payload := rndReplyMsg(
		message.MsgReplyUtf8,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgReplyUtf8, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
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
	require.Equal(t, message.MsgReplyUtf16, actual.MsgType)
	require.Equal(t, id, actual.MsgIdentifier[:])
	require.Equal(t, id, actual.MsgIdentifierBytes)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSignalBinary tests parsing of a named binary encoded signal
func TestMsgParseSignalBinary(t *testing.T) {
	encoded, name, payload := rndSignalMsg(
		message.MsgSignalBinary,
		1, 255,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgSignalBinary, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSignalUtf8 tests parsing of a named UTF8 encoded signal
func TestMsgParseSignalUtf8(t *testing.T) {
	encoded, name, payload := rndSignalMsg(
		message.MsgSignalUtf8,
		1, 255,
		1, 1024*64,
	)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgSignalUtf8, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
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
	require.Equal(t, message.MsgSignalUtf16, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Equal(t, name, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
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
	encoded := []byte{message.MsgNotifySessionCreated}
	// Add session payload
	encoded = append(encoded, payload.Data...)

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgNotifySessionCreated, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, payload, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseSessClosedSig tests parsing of session closed signal
func TestMsgParseSessClosedSig(t *testing.T) {
	// Compose encoded message
	// Add type flag
	encoded := []byte{message.MsgNotifySessionClosed}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgNotifySessionClosed, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{}, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseHeartbeat tests parsing of heartbeat messages
func TestMsgParseHeartbeat(t *testing.T) {
	// Compose encoded message
	// Add type flag
	encoded := []byte{message.MsgHeartbeat}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	require.NotNil(t, actual.MsgBuffer)
	require.Equal(t, message.MsgHeartbeat, actual.MsgType)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, actual.MsgIdentifierBytes)
	require.Equal(t, [8]byte{}, actual.MsgIdentifier)
	require.Nil(t, actual.MsgName)
	require.Equal(t, pld.Payload{}, actual.MsgPayload)
	require.Equal(t, message.ServerConfiguration{}, actual.ServerConfiguration)
}

// TestMsgParseUnknownMessageType tests parsing of messages
// with unknown message type
func TestMsgParseUnknownMessageType(t *testing.T) {
	msgOfUnknownType := make([]byte, 1)
	msgOfUnknownType[0] = byte(255)

	actual := message.NewMessage(1024)
	typeDetermined, _ := actual.ReadBytes(msgOfUnknownType)
	require.False(t, typeDetermined, "Expected type not to be determined")
}
