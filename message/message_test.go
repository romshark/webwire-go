package message

import (
	"bytes"
	"encoding/json"
	"reflect"
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

/****************************************************************\
	Constructors
\****************************************************************/

// TestMsgNewNamelessReqMsg tests NewNamelessRequestMessage
func TestMsgNewNamelessReqMsg(t *testing.T) {
	id := genRndMsgIdentifier()
	// sessionKey := generateSessionKey()
	sessionKey := "somesamplesessionkey"

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgRestoreSession}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add session key to payload
	expected = append(expected, []byte(sessionKey)...)

	actual := NewNamelessRequestMessage(
		MsgRestoreSession,
		id,
		[]byte(sessionKey),
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewEmptyReqMsg tests NewEmptyRequestMessage
func TestMsgNewEmptyReqMsg(t *testing.T) {
	id := genRndMsgIdentifier()

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgCloseSession}
	// Add identifier
	expected = append(expected, id[:]...)

	actual := NewEmptyRequestMessage(MsgCloseSession, id)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReqMsgBinary tests NewRequestMessage
// using default binary payload encoding
func TestMsgNewReqMsgBinary(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgRequestBinary}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload
	// (skip header padding byte, not necessary in case of binary encoding)
	expected = append(expected, payload.Data...)

	actual := NewRequestMessage(
		id,
		string(name),
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReqMsgUtf8 tests NewRequestMessage using UTF8 payload encoding
func TestMsgNewReqMsgUtf8(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgRequestUtf8}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload
	// (skip header padding byte, not necessary in case of UTF8 encoding)
	expected = append(expected, payload.Data...)

	actual := NewRequestMessage(
		id,
		string(name),
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReqMsgUtf16 tests NewRequestMessage using UTF8 payload encoding
func TestMsgNewReqMsgUtf16(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgRequestUtf16}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		expected = append(expected, byte(0))
	}
	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewRequestMessage(
		id,
		string(name),
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReqMsgUtf16OddNameLen tests NewRequestMessage using
// UTF16 payload encoding and a name of odd length
func TestMsgNewReqMsgUtf16OddNameLen(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgRequestUtf16}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add name length flag
	expected = append(expected, byte(3))
	// Add name of odd length
	expected = append(expected, []byte("odd")...)
	// Add header padding
	expected = append(expected, byte(0))
	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewRequestMessage(id, "odd", payload.Encoding, payload.Data)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReplyMsgBinary tests NewReplyMessage
// using default binary payload encoding
func TestMsgNewReplyMsgBinary(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgReplyBinary}
	// Add identifier
	expected = append(expected, id[:]...)

	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewReplyMessage(
		id,
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReplyMsgUtf8 tests NewReplyMessage using UTF8 payload encoding
func TestMsgNewReplyMsgUtf8(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgReplyUtf8}
	// Add identifier
	expected = append(expected, id[:]...)

	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewReplyMessage(
		id,
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReplyMsgUtf16 tests NewReplyMessage using UTF16 payload encoding
func TestMsgNewReplyMsgUtf16(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgReplyUtf16}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add header padding byte (necessary in case of a UTF16 encoded reply)
	expected = append(expected, 0)

	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewReplyMessage(
		id,
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewSigMsgBinary tests NewSignalMessage
// using the default binary encoding
func TestMsgNewSigMsgBinary(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgSignalBinary}
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload (skip header padding byte in case of binary encoding)
	expected = append(expected, payload.Data...)

	actual := NewSignalMessage(
		string(name),
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewSigMsgUtf8 tests NewSignalMessage using UTF8 encoding
func TestMsgNewSigMsgUtf8(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgSignalUtf8}
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add payload (skip header padding byte in case of UTF8 encoding)
	expected = append(expected, payload.Data...)

	actual := NewSignalMessage(
		string(name),
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewSigMsgUtf16 tests NewSignalMessage using UTF16 encoding
func TestMsgNewSigMsgUtf16(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgSignalUtf16}
	// Add name length flag
	expected = append(expected, byte(len(name)))
	// Add name
	expected = append(expected, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		expected = append(expected, byte(0))
	}
	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewSignalMessage(
		string(name),
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewSigMsgUtf16OddNameLen tests NewSignalMessage using UTF16 encoding
// and a name of odd length to ensure a header padding byte is used
func TestMsgNewSigMsgUtf16OddNameLen(t *testing.T) {
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgSignalUtf16}
	// Add name length flag
	expected = append(expected, byte(3))
	// Add name of odd length
	expected = append(expected, []byte("odd")...)
	// Add header padding
	expected = append(expected, byte(0))
	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewSignalMessage(
		"odd",
		payload.Encoding,
		payload.Data,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

/****************************************************************\
	Parser - invalid messages (too short)
\****************************************************************/

// TestMsgParseInvalidMessageTooShort tests parsing of an invalid
// empty message
func TestMsgParseInvalidMessageTooShort(t *testing.T) {
	invalidMessage := make([]byte, 0)

	var actual Message
	if typeDetermined, _ := actual.Parse(invalidMessage); typeDetermined {
		t.Fatalf(
			"Expected type to not be determined " +
				"when parsing empty message",
		)
	}
}

// TestMsgParseInvalidReplyTooShort tests parsing of an invalid
// binary/UTF8 reply message which is too short to be considered valid
func TestMsgParseInvalidReplyTooShort(t *testing.T) {
	lenTooShort := MsgMinLenReply - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgReplyBinary

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid reply message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidReplyUtf16TooShort tests parsing of an invalid
// UTF16 reply message which is too short to be considered valid
func TestMsgParseInvalidReplyUtf16TooShort(t *testing.T) {
	lenTooShort := MsgMinLenReplyUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgReplyUtf16

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid UTF16 reply message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidRequestTooShort tests parsing of an invalid
// binary/UTF8 request message which is too short to be considered valid
func TestMsgParseInvalidRequestTooShort(t *testing.T) {
	lenTooShort := MsgMinLenRequest - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestBinary

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid request message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidRequestUtf16TooShort tests parsing of an invalid
// UTF16 request message which is too short to be considered valid
func TestMsgParseInvalidRequestUtf16TooShort(t *testing.T) {
	lenTooShort := MsgMinLenRequestUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestUtf16

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid UTF16 "+
				"encoded request message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidRestrSessReqTooShort tests parsing of an invalid
// session restoration request message which is too short
// to be considered valid
func TestMsgParseInvalidRestrSessReqTooShort(t *testing.T) {
	lenTooShort := MsgMinLenRestoreSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRestoreSession

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid session restoration "+
				"request message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidSessCloseReqTooShort tests parsing of an invalid
// session destruction request message which is too short
// to be considered valid
func TestMsgParseInvalidSessCloseReqTooShort(t *testing.T) {
	lenTooShort := MsgMinLenCloseSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgCloseSession

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid session destruction "+
				"request message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidSessCreatedSigTooShort tests parsing of an invalid
// session creation notification message which is too short
// to be considered valid
func TestMsgParseInvalidSessCreatedSigTooShort(t *testing.T) {
	lenTooShort := MsgMinLenSessionCreated - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSessionCreated

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid session creation "+
				"notification message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidSignalTooShort tests parsing of an invalid
// binary/UTF8 signal message which is too short to be considered valid
func TestMsgParseInvalidSignalTooShort(t *testing.T) {
	lenTooShort := MsgMinLenSignal - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSignalBinary

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid signal message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidSignalUtf16TooShort tests parsing of an invalid
// UTF16 signal message which is too short to be considered valid
func TestMsgParseInvalidSignalUtf16TooShort(t *testing.T) {
	lenTooShort := MsgMinLenSignalUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSignalUtf16

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid UTF16 signal message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidErrorReplyTooShort tests parsing of an invalid
// error reply message which is too short to be considered valid
func TestMsgParseInvalidErrorReplyTooShort(t *testing.T) {
	lenTooShort := MsgMinLenErrorReply - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgErrorReply

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid error reply message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseInvalidSpecialReplyTooShort tests parsing of an invalid
// special reply message which is too short to be considered valid
func TestMsgParseInvalidSpecialReplyTooShort(t *testing.T) {
	invalidMessage := make([]byte, 8)

	// Internal error is a special reply message type
	invalidMessage[0] = MsgInternalError

	if _, err := tryParse(t, invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid special reply message " +
				"(too short: 8)",
		)
	}
}

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

/****************************************************************\
	Parser - invalid messages (corrupt name length flags)
\****************************************************************/

// TestMsgParseRequestCorruptNameLenFlag tests parsing of a named
// Binary/UTF8 encoded request with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseRequestCorruptNameLenFlag(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(10 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(MsgRequestBinary)
	// Add identifier
	encoded.Write(id[:])

	// Add corrupt name length flag (too big) and skip the name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	if _, err := tryParse(t, encoded.Bytes()); err == nil {
		t.Fatal(
			"Expected Parse to return an error due to corrupt name length flag",
		)
	}
}

// TestMsgParseRequestUtf16CorruptNameLenFlag tests parsing of a named
// UTF16 encoded request with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseRequestUtf16CorruptNameLenFlag(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(10 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(MsgRequestUtf16)
	// Add identifier
	encoded.Write(id[:])

	// Add corrupt name length flag (too big) and skip actual name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	if _, err := tryParse(t, encoded.Bytes()); err == nil {
		t.Fatal(
			"Expected Parse to return an error due to corrupt name length flag",
		)
	}
}

// TestMsgParseSignalCorruptNameLenFlag tests parsing of a named
// Binary/UTF8 encoded signal with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseSignalCorruptNameLenFlag(t *testing.T) {
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(2 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(MsgSignalBinary)

	// Add corrupt name length flag (too big) and skip the name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	if _, err := tryParse(t, encoded.Bytes()); err == nil {
		t.Fatal(
			"Expected Parse to return an error due to corrupt name length flag",
		)
	}
}

// TestMsgParseSignalUtf16CorruptNameLenFlag tests parsing of a named
// UTF16 encoded signal with a corrupted input stream
// (name length flag doesn't correspond to actual name length)
func TestMsgParseSignalUtf16CorruptNameLenFlag(t *testing.T) {
	payload := pld.Payload{
		Encoding: pld.Binary,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	encoded := &bytes.Buffer{}
	encoded.Grow(2 + len(payload.Data))

	// Add type flag
	encoded.WriteByte(MsgSignalUtf16)

	// Add corrupt name length flag (too big) and skip the name field
	encoded.WriteByte(255)

	// Add payload
	encoded.Write(payload.Data)

	// Parse
	if _, err := tryParse(t, encoded.Bytes()); err == nil {
		t.Fatal(
			"Expected Parse to return an error due to corrupt name length flag",
		)
	}
}

/****************************************************************\
	Parser - invalid input (corrupt payload)
\****************************************************************/

// TestMsgParseReplyUtf16CorruptInput tests parsing of
// UTF16 encoded reply message with a corrupted input stream
// (length not divisible by 2)
func TestMsgParseReplyUtf16CorruptInput(t *testing.T) {
	id := genRndMsgIdentifier()
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgReplyUtf16}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add header padding byte due to UTF16 encoding
	encoded = append(encoded, byte(0))
	// Add payload
	encoded = append(encoded, payload.Data...)

	// Parse
	if _, err := tryParse(t, encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

// TestMsgParseRequestUtf16CorruptInput tests parsing of a named
// UTF16 encoded request with a corrupted input stream
// (length not divisible by 2)
func TestMsgParseRequestUtf16CorruptInput(t *testing.T) {
	id := genRndMsgIdentifier()
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgRequestUtf16}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		encoded = append(encoded, byte(0))
	}
	// Add payload
	encoded = append(encoded, payload.Data...)

	// Parse
	if _, err := tryParse(t, encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

// TestMsgParseSignalUtf16CorruptInput tests parsing of a named
// UTF16 encoded signal with a corrupt unaligned input stream
// (length not divisible by 2)
func TestMsgParseSignalUtf16CorruptInput(t *testing.T) {
	name := genRndName(1, 255)
	payload := pld.Payload{
		Encoding: pld.Utf16,
		Data:     []byte("invalid"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgSignalUtf16}
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add header padding if necessary
	if len(name)%2 != 0 {
		encoded = append(encoded, byte(0))
	}
	// Add payload
	encoded = append(encoded, payload.Data...)

	// Parse
	if _, err := tryParse(t, encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

/****************************************************************\
	Constructors - invalid input (corrupt name length flags)
\****************************************************************/

// TestMsgNewReplyMsgUtf16CorruptPayload tests NewReplyMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReplyMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewReplyMessage(
		genRndMsgIdentifier(),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
	)
}

// TestMsgNewReqMsgUtf16CorruptPayload tests NewRequestMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReqMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewRequestMessage(
		genRndMsgIdentifier(),
		string(genRndName(1, 255)),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
	)
}

// TestMsgNewSigMsgUtf16CorruptPayload tests NewSignalMessage
// using UTF16 payload encoding passing corrupt data
// (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewSigMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewSignalMessage(
		string(genRndName(1, 255)),
		pld.Utf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		[]byte("invalid"),
	)
}

/****************************************************************\
	Constructors - unexpected parameters (panics)
\****************************************************************/

// TestMsgNewReqMsgNoNameNoPayload tests calling
// the request message constructor without both the name and the payload
func TestMsgNewReqMsgNoNameNoPayload(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal(
				"Expected a panic after calling the " +
					" request message constructor without both the name " +
					"and the payload",
			)
		}
	}()

	id := genRndMsgIdentifier()

	NewRequestMessage(
		id,
		"",
		pld.Binary,
		nil,
	)
}

// TestMsgNewReqMsgNameTooLong tests NewRequestMessage with a too long name
func TestMsgNewReqMsgNameTooLong(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("Expected panic after passing a too long request name")
		}
	}()

	tooLongNamelength := 256
	nameBuf := make([]byte, tooLongNamelength)
	for i := 0; i < tooLongNamelength; i++ {
		nameBuf[i] = 'a'
	}

	NewRequestMessage(
		genRndMsgIdentifier(),
		string(nameBuf),
		0,
		nil,
	)
}

// TestMsgNewReqMsgInvalidCharsetBelowAscii32 tests NewRequestMessage
// with an invalid character input below the ASCII 7 bit 32nd character
func TestMsgNewReqMsgInvalidCharsetBelowAscii32(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic after passing an invalid name character set",
			)
		}
	}()

	// Generate invalid name using a character
	// below the ASCII 7 bit 32nd character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(31)

	NewRequestMessage(
		genRndMsgIdentifier(),
		string(invalidNameBytes),
		0,
		nil,
	)
}

// TestMsgNewReqMsgInvalidCharsetAboveAscii126 tests NewRequestMessage
// with an invalid character input above the ASCII 7 bit 126th character
func TestMsgNewReqMsgInvalidCharsetAboveAscii126(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic after passing an invalid name character set",
			)
		}
	}()

	// Generate invalid name using a character
	// above the ASCII 7 bit 126th character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(127)

	NewRequestMessage(
		genRndMsgIdentifier(),
		string(invalidNameBytes),
		0,
		nil,
	)
}

// TestMsgNewSigMsgNameTooLong tests NewSignalMessage with a too long name
func TestMsgNewSigMsgNameTooLong(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("Expected panic after passing a too long signal name")
		}
	}()

	tooLongNamelength := 256
	nameBuf := make([]byte, tooLongNamelength)
	for i := 0; i < tooLongNamelength; i++ {
		nameBuf[i] = 'a'
	}

	NewSignalMessage(string(nameBuf), 0, nil)
}

// TestMsgNewSigMsgInvalidCharsetBelowAscii32 tests NewSignalMessage
// with an invalid character input below the ASCII 7 bit 32nd character
func TestMsgNewSigMsgInvalidCharsetBelowAscii32(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic after passing an invalid name character set",
			)
		}
	}()

	// Generate invalid name using a character
	// below the ASCII 7 bit 32nd character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(31)

	NewSignalMessage(string(invalidNameBytes), 0, nil)
}

// TestMsgNewSigMsgInvalidCharsetAboveAscii126 tests NewSignalMessage
// with an invalid character input above ASCII 7 bit 126th character
func TestMsgNewSigMsgInvalidCharsetAboveAscii126(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic after passing an invalid name character set",
			)
		}
	}()

	// Generate invalid name using a character
	// above the ASCII 7 bit 126th character
	invalidNameBytes := make([]byte, 1)
	invalidNameBytes[0] = byte(127)

	NewSignalMessage(string(invalidNameBytes), 0, nil)
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
		func(msgType byte) {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatalf(
						"Expected panic after passing " +
							"a non-special request reply message type",
					)
				}
			}()
			NewSpecialRequestReplyMessage(MsgErrorReply, genRndMsgIdentifier())
		}(tp)
	}
}

// TestMsgNewErrorReplyMessageNoCode tests NewErrorReplyMessage
// with no error code which is invalid.
func TestMsgNewErrorReplyMessageNoCode(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic when creating an error reply message " +
					"with no error code ",
			)
		}
	}()

	NewErrorReplyMessage(genRndMsgIdentifier(), "", "sample error message")
}

// TestMsgNewErrorReplyMessageCodeTooLong tests NewErrorReplyMessage
// with an error code that's surpassing the 255 character limit.
func TestMsgNewErrorReplyMessageCodeTooLong(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic when creating an error reply message " +
					"with no error code ",
			)
		}
	}()

	tooLongCode := make([]byte, 256)
	for i := 0; i < 256; i++ {
		tooLongCode[i] = 'a'
	}

	NewErrorReplyMessage(
		genRndMsgIdentifier(),
		string(tooLongCode),
		"sample error message",
	)
}

// TestMsgNewErrorReplyMessageCodeCharsetBelowAscii32 tests NewErrorReplyMessage
// with an invalid character input below the ASCII 7 bit 32nd character
func TestMsgNewErrorReplyMessageCodeCharsetBelowAscii32(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic after passing an invalid error code " +
					" containing a character below the 32th ASCII 7bit char",
			)
		}
	}()

	// Generate invalid error code using a character
	// below the ASCII 7 bit 32nd character
	invalidCodeBytes := make([]byte, 1)
	invalidCodeBytes[0] = byte(31)

	NewErrorReplyMessage(
		genRndMsgIdentifier(),
		string(invalidCodeBytes),
		"sample error message",
	)
}

// TestMsgNewErrorReplyMessageCodeCharsetAboveAscii126 tests
// NewErrorReplyMessage with an invalid character input
// above the ASCII 7 bit 126th character
func TestMsgNewErrorReplyMessageCodeCharsetAboveAscii126(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf(
				"Expected panic after passing an invalid error code " +
					" containing a character above the 126th ASCII 7bit char",
			)
		}
	}()

	// Generate invalid error code using a character
	// above the ASCII 7 bit 126th character
	invalidCodeBytes := make([]byte, 1)
	invalidCodeBytes[0] = byte(127)

	NewErrorReplyMessage(
		genRndMsgIdentifier(),
		string(invalidCodeBytes),
		"sample error message",
	)
}
