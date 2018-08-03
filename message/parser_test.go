package message

import (
	"encoding/json"
	"testing"
	"time"

	pld "github.com/qbeon/webwire-go/payload"
)

/****************************************************************\
	Parser
\****************************************************************/

// TestMsgParseCloseSessReq tests parsing of a session destruction request
func TestMsgParseCloseSessReq(t *testing.T) {
	id := genRndMsgIdentifier()

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgCloseSession}
	// Add identifier
	encoded = append(encoded, id[:]...)

	// Initialize expected message
	expected := Message{
		Type:       MsgCloseSession,
		Identifier: id,
		Name:       "",
		Payload: pld.Payload{
			Encoding: pld.Binary,
			Data:     nil,
		},
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseRestrSessReq tests parsing of a session restoration request
func TestMsgParseRestrSessReq(t *testing.T) {
	id := genRndMsgIdentifier()

	//sessionKey := sess.GenerateSessionKey()
	sessionKey := "somesamplesessionkey"

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgRestoreSession}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add session key to payload
	encoded = append(encoded, sessionKey[:]...)

	// Initialize expected message with the session key in the payload
	expected := Message{
		Type:       MsgRestoreSession,
		Identifier: id,
		Name:       "",
		Payload: pld.Payload{
			Encoding: pld.Binary,
			Data:     []byte(sessionKey),
		},
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseRequestBinary tests parsing of a named binary encoded request
func TestMsgParseRequestBinary(t *testing.T) {
	encoded, id, name, payload := rndRequestMsg(
		MsgRequestBinary,
		1, 255,
		1, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:       MsgRequestBinary,
		Identifier: id,
		Name:       string(name),
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseRequestUtf8 tests parsing of a named UTF8 encoded request
func TestMsgParseRequestUtf8(t *testing.T) {
	encoded, id, name, payload := rndRequestMsg(
		MsgRequestUtf8,
		2, 255,
		16, 16,
	)

	// Initialize expected message
	expected := Message{
		Type:       MsgRequestUtf8,
		Identifier: id,
		Name:       string(name),
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseRequestUtf16 tests parsing of a named UTF16 encoded request
func TestMsgParseRequestUtf16(t *testing.T) {
	encoded, id, name, payload := rndRequestMsgUtf16(
		1, 255,
		2, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:       MsgRequestUtf16,
		Identifier: id,
		Name:       string(name),
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseReplyBinary tests parsing of binary encoded reply message
func TestMsgParseReplyBinary(t *testing.T) {
	encoded, id, payload := rndReplyMsg(
		MsgReplyBinary,
		1, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:       MsgReplyBinary,
		Identifier: id,
		Name:       "",
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseReplyUtf8 tests parsing of UTF8 encoded reply message
func TestMsgParseReplyUtf8(t *testing.T) {
	encoded, id, payload := rndReplyMsg(
		MsgReplyUtf8,
		1, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:       MsgReplyUtf8,
		Identifier: id,
		Name:       "",
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseReplyUtf16 tests parsing of UTF16 encoded reply message
func TestMsgParseReplyUtf16(t *testing.T) {
	encoded, id, payload := rndReplyMsgUtf16(
		2, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:       MsgReplyUtf16,
		Identifier: id,
		Name:       "",
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseSignalBinary tests parsing of a named binary encoded signal
func TestMsgParseSignalBinary(t *testing.T) {
	encoded, name, payload := rndSignalMsg(
		MsgSignalBinary,
		1, 255,
		1, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:    MsgSignalBinary,
		Name:    string(name),
		Payload: payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseSignalUtf8 tests parsing of a named UTF8 encoded signal
func TestMsgParseSignalUtf8(t *testing.T) {
	encoded, name, payload := rndSignalMsg(
		MsgSignalUtf8,
		1, 255,
		1, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:    MsgSignalUtf8,
		Name:    string(name),
		Payload: payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseSignalUtf16 tests parsing of a named UTF16 encoded signal
func TestMsgParseSignalUtf16(t *testing.T) {
	encoded, name, payload := rndSignalMsgUtf16(
		1, 255,
		2, 1024*64,
	)

	// Initialize expected message
	expected := Message{
		Type:       MsgSignalUtf16,
		Identifier: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:       string(name),
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
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
	if err != nil {
		t.Fatalf("Couldn't marshal session object: %s", err)
	}
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     marshalledSession,
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgSessionCreated}
	// Add session payload
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		Type:       MsgSessionCreated,
		Identifier: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:       "",
		Payload:    payload,
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseSessClosedSig tests parsing of session sloed signal
func TestMsgParseSessClosedSig(t *testing.T) {
	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgSessionClosed}

	// Initialize expected message
	expected := Message{
		Type:       MsgSessionClosed,
		Identifier: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:       "",
		Payload:    pld.Payload{},
	}

	// Parse
	actual := tryParseNoErr(t, encoded)

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseUnknownMessageType tests parsing of messages
// with unknown message type
func TestMsgParseUnknownMessageType(t *testing.T) {
	msgOfUnknownType := make([]byte, 1)
	msgOfUnknownType[0] = byte(255)

	var actual Message
	if typeDetermined, _ := actual.Parse(msgOfUnknownType); typeDetermined {
		t.Fatalf("Expected type not to be determined")
	}
}
