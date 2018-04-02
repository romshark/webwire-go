package webwire

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func comparePayload(t *testing.T, expected, actual []byte) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"Payload differs:\n expected: '%s'\n actual:   '%s'",
			string(expected),
			string(actual),
		)
	}
}

func compareMessages(t *testing.T, expected, actual Message) {
	if actual.Client != expected.Client {
		t.Errorf("message.Client differs: %p | %p", expected.Client, actual.Client)
	}
	if actual.msgType != expected.msgType {
		t.Errorf(
			"message.msgType differs:"+
				"\n expected: '%s'\n actual:   '%s'",
			string(expected.msgType),
			string(actual.msgType),
		)
	}
	if !reflect.DeepEqual(actual.id, expected.id) {
		t.Errorf(
			"message.id differs:"+
				"\n expected: '%s'\n actual:   '%s'",
			string(expected.id[:]),
			string(actual.id[:]),
		)
	}
	if actual.Name != expected.Name {
		t.Errorf(
			"message.Name differs: %s | %s",
			expected.Name,
			actual.Name,
		)
	}
	if actual.Payload.Encoding != expected.Payload.Encoding {
		t.Errorf(
			"message.Payload.Encoding differs: %v | %v",
			expected.Payload.Encoding,
			actual.Payload.Encoding,
		)
	}
	comparePayload(t, expected.Payload.Data, actual.Payload.Data)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Deep equality assertion failed")
	}
}

func genRndMsgID() (randID [8]byte) {
	rand.Read(randID[:])
	return randID
}

func genRndName() string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = " !\"#$%&'()*+,-./0123456789:;<=>?@" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"[\\]^_`" +
		"abcdefghijklmnopqrstuvwxyz" +
		"{|}~"
	const min = 1
	const max = 255
	nameBytes := make([]byte, min+rand.Intn(max-min))
	for i := range nameBytes {
		nameBytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(nameBytes)
}

// TestMsgParseInvalidSessCloseReqTooShort tests parsing of an invalid session destruction request
// message which is too short to be considered valid
func TestMsgParseInvalidSessCloseReqTooShort(t *testing.T) {
	lenTooShort := MsgMinLenCloseSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgCloseSession

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid session destruction request message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseCloseSessReq tests parsing of a session destruction request
func TestMsgParseCloseSessReq(t *testing.T) {
	id := genRndMsgID()

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgCloseSession}
	// Add identifier
	encoded = append(encoded, id[:]...)

	// Initialize expected message
	expected := Message{
		msgType: MsgCloseSession,
		id:      id,
		Name:    "",
		Payload: Payload{
			Encoding: EncodingBinary,
			Data:     nil,
		},
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseInvalidRestrSessReqTooShort tests parsing of an invalid session restoration request
// message which is too short to be considered valid
func TestMsgParseInvalidRestrSessReqTooShort(t *testing.T) {
	lenTooShort := MsgMinLenRestoreSession - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRestoreSession

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid session restoration request message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseRestrSessReq tests parsing of a session restoration request
func TestMsgParseRestrSessReq(t *testing.T) {
	id := genRndMsgID()
	sessionKey := generateSessionKey()

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgRestoreSession}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add session key to payload
	encoded = append(encoded, sessionKey[:]...)

	// Initialize expected message with the session key in the payload
	expected := Message{
		msgType: MsgRestoreSession,
		id:      id,
		Name:    "",
		Payload: Payload{
			Encoding: EncodingBinary,
			Data:     []byte(sessionKey),
		},
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseInvalidRequestTooShort tests parsing of an invalid binary/UTF8 request message
// which is too short to be considered valid
func TestMsgParseInvalidRequestTooShort(t *testing.T) {
	lenTooShort := MsgMinLenRequest - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestBinary

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid request message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseRequestBinary tests parsing of a named binary encoded request
func TestMsgParseRequestBinary(t *testing.T) {
	id := genRndMsgID()
	name := genRndName()
	payload := Payload{
		Encoding: EncodingBinary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgRequestBinary}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add payload (skip header padding byte in case of binary encoding)
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		msgType: MsgRequestBinary,
		id:      id,
		Name:    name,
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseRequestUtf8 tests parsing of a named UTF8 encoded request
func TestMsgParseRequestUtf8(t *testing.T) {
	id := genRndMsgID()
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf8,
		Data:     []byte("random utf8 payload"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgRequestUtf8}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add payload (skip header padding byte in case of UTF8 encoding)
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		msgType: MsgRequestUtf8,
		id:      id,
		Name:    name,
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseInvalidRequestUtf16TooShort tests parsing of an invalid utf16 request message
// which is too short to be considered valid
func TestMsgParseInvalidRequestUtf16TooShort(t *testing.T) {
	lenTooShort := MsgMinLenRequestUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestUtf16

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid UTF16 encoded request message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseRequestUtf16 tests parsing of a named UTF16 encoded request
func TestMsgParseRequestUtf16(t *testing.T) {
	id := genRndMsgID()
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
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

	// Initialize expected message
	expected := Message{
		msgType: MsgRequestUtf16,
		id:      id,
		Name:    name,
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseRequestUtf16CorruptInput tests parsing of a named UTF16 encoded request
// with a corrupted input stream (length not divisible by 2)
func TestMsgParseRequestUtf16CorruptInput(t *testing.T) {
	id := genRndMsgID()
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf16,
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
	var actual Message
	if err := actual.Parse(encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

// TestMsgParseInvalidReplyTooShort tests parsing of an invalid binary/UTF8 reply message
// which is too short to be considered valid
func TestMsgParseInvalidReplyTooShort(t *testing.T) {
	lenTooShort := MsgMinLenReply - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgReplyBinary

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid reply message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseReplyBinary tests parsing of binary encoded reply message
func TestMsgParseReplyBinary(t *testing.T) {
	id := genRndMsgID()
	payload := Payload{
		Encoding: EncodingBinary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgReplyBinary}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add payload (skip header padding byte in case of binary encoding)
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		msgType: MsgReplyBinary,
		id:      id,
		Name:    "",
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseReplyUtf8 tests parsing of UTF8 encoded reply message
func TestMsgParseReplyUtf8(t *testing.T) {
	id := genRndMsgID()
	payload := Payload{
		Encoding: EncodingUtf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgReplyUtf8}
	// Add identifier
	encoded = append(encoded, id[:]...)
	// Add payload (skip header padding byte in case of UTF8 encoding)
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		msgType: MsgReplyUtf8,
		id:      id,
		Name:    "",
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseInvalidReplyUtf16TooShort tests parsing of an invalid UTF16 reply message
// which is too short to be considered valid
func TestMsgParseInvalidReplyUtf16TooShort(t *testing.T) {
	lenTooShort := MsgMinLenReplyUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgRequestUtf16

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid UTF16 reply message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseReplyUtf16 tests parsing of UTF16 encoded reply message
func TestMsgParseReplyUtf16(t *testing.T) {
	id := genRndMsgID()
	payload := Payload{
		Encoding: EncodingUtf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
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

	// Initialize expected message
	expected := Message{
		msgType: MsgReplyUtf16,
		id:      id,
		Name:    "",
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseReplyUtf16CorruptInput tests parsing of UTF16 encoded reply message
// with a corrupted input stream (length not divisible by 2)
func TestMsgParseReplyUtf16CorruptInput(t *testing.T) {
	id := genRndMsgID()
	payload := Payload{
		Encoding: EncodingUtf16,
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
	var actual Message
	if err := actual.Parse(encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

// TestMsgParseInvalidSignalTooShort tests parsing of an invalid binary/UTF8 signal message
// which is too short to be considered valid
func TestMsgParseInvalidSignalTooShort(t *testing.T) {
	lenTooShort := MsgMinLenSignal - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSignalBinary

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid signal message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseSignalBinary tests parsing of a named binary encoded signal
func TestMsgParseSignalBinary(t *testing.T) {
	name := genRndName()
	payload := Payload{
		Encoding: EncodingBinary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgSignalBinary}
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add payload (skip header padding byte in case of binary encoding)
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		msgType: MsgSignalBinary,
		id:      [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:    name,
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseSignalUtf8 tests parsing of a named UTF8 encoded signal
func TestMsgParseSignalUtf8(t *testing.T) {
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgSignalUtf8}
	// Add name length flag
	encoded = append(encoded, byte(len(name)))
	// Add name
	encoded = append(encoded, []byte(name)...)
	// Add payload (skip header padding byte in case of UTF8 encoding)
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		msgType: MsgSignalUtf8,
		id:      [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:    name,
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseInvalidSignalUtf16TooShort tests parsing of an invalid UTF16 signal message
// which is too short to be considered valid
func TestMsgParseInvalidSignalUtf16TooShort(t *testing.T) {
	lenTooShort := MsgMinLenSignalUtf16 - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSignalUtf16

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid UTF16 signal message (too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseSignalUtf16 tests parsing of a named UTF16 encoded signal
func TestMsgParseSignalUtf16(t *testing.T) {
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf16,
		Data:     []byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
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

	// Initialize expected message
	expected := Message{
		msgType: MsgSignalUtf16,
		id:      [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:    name,
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgParseSignalUtf16CorruptInput tests parsing of a named UTF16 encoded signal
// with a corrupt unaligned input stream (length not divisible by 2)
func TestMsgParseSignalUtf16CorruptInput(t *testing.T) {
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf16,
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
	var actual Message
	if err := actual.Parse(encoded); err == nil {
		t.Fatal("Expected Parse to return an error due to corrupt input stream")
	}
}

// TestMsgParseInvalidSessCreatedSigTooShort tests parsing of an invalid session creation
// notification message which is too short to be considered valid
func TestMsgParseInvalidSessCreatedSigTooShort(t *testing.T) {
	lenTooShort := MsgMinLenSessionCreated - 1
	invalidMessage := make([]byte, lenTooShort)

	invalidMessage[0] = MsgSessionCreated

	var actual Message
	if err := actual.Parse(invalidMessage); err == nil {
		t.Fatalf(
			"Expected error while parsing invalid session creation notification message "+
				"(too short: %d)",
			lenTooShort,
		)
	}
}

// TestMsgParseSessCreatedSig tests parsing of session created signal
func TestMsgParseSessCreatedSig(t *testing.T) {
	sessionKey := generateSessionKey()
	session := Session{
		Key:      sessionKey,
		Creation: time.Now(),
		Info:     nil,
	}
	marshalledSession, err := json.Marshal(&session)
	if err != nil {
		t.Fatalf("Couldn't marshal session object: %s", err)
	}
	payload := Payload{
		Encoding: EncodingBinary,
		Data:     marshalledSession,
	}

	// Compose encoded message
	// Add type flag
	encoded := []byte{MsgSessionCreated}
	// Add session payload
	encoded = append(encoded, payload.Data...)

	// Initialize expected message
	expected := Message{
		msgType: MsgSessionCreated,
		id:      [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:    "",
		Payload: payload,
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

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
		msgType: MsgSessionClosed,
		id:      [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		Name:    "",
		Payload: Payload{},
	}

	// Parse
	var actual Message
	if err := actual.Parse(encoded); err != nil {
		t.Fatalf("Failed parsing: %s", err)
	}

	// Compare
	compareMessages(t, expected, actual)
}

// TestMsgNewNamelessReqMsg tests the NewNamelessRequestMessage method
func TestMsgNewNamelessReqMsg(t *testing.T) {
	id := genRndMsgID()
	sessionKey := generateSessionKey()

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgRestoreSession}
	// Add identifier
	expected = append(expected, id[:]...)
	// Add session key to payload
	expected = append(expected, []byte(sessionKey)...)

	actual := NewNamelessRequestMessage(MsgRestoreSession, id, []byte(sessionKey))

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewEmptyReqMsg tests the NewEmptyRequestMessage method
func TestMsgNewEmptyReqMsg(t *testing.T) {
	id := genRndMsgID()

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

// TestMsgNewReqMsgNameTooLong tests the NewRequestMessage method with a too long name
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

	NewRequestMessage(genRndMsgID(), string(nameBuf), Payload{})
}

// TestMsgNewReqMsgBinary tests the NewRequestMessage method using default binary payload encoding
func TestMsgNewReqMsgBinary(t *testing.T) {
	id := genRndMsgID()
	name := genRndName()
	payload := Payload{
		Encoding: EncodingBinary,
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
	// Add payload (skip header padding byte, not necessary in case of binary encoding)
	expected = append(expected, payload.Data...)

	actual := NewRequestMessage(id, name, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReqMsgUtf8 tests the NewRequestMessage method using UTF8 payload encoding
func TestMsgNewReqMsgUtf8(t *testing.T) {
	id := genRndMsgID()
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf8,
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
	// Add payload (skip header padding byte, not necessary in case of UTF8 encoding)
	expected = append(expected, payload.Data...)

	actual := NewRequestMessage(id, name, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReqMsgUtf16 tests the NewRequestMessage method using UTF8 payload encoding
func TestMsgNewReqMsgUtf16(t *testing.T) {
	id := genRndMsgID()
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf16,
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

	actual := NewRequestMessage(id, name, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReqMsgUtf16CorruptPayload tests the NewRequestMessage method using UTF16 payload encoding
// passing corrupt data (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReqMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected NewRequestMessage to panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewRequestMessage(genRndMsgID(), genRndName(), Payload{
		Encoding: EncodingUtf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		Data: []byte("invalid"),
	})
}

// TestMsgNewReplyMsgBinary tests the NewReplyMessage method using default binary payload encoding
func TestMsgNewReplyMsgBinary(t *testing.T) {
	id := genRndMsgID()
	payload := Payload{
		Encoding: EncodingBinary,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgReplyBinary}
	// Add identifier
	expected = append(expected, id[:]...)

	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewReplyMessage(id, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReplyMsgUtf8 tests the NewReplyMessage method using UTF8 payload encoding
func TestMsgNewReplyMsgUtf8(t *testing.T) {
	id := genRndMsgID()
	payload := Payload{
		Encoding: EncodingUtf8,
		Data:     []byte("random payload data"),
	}

	// Compose encoded message
	// Add type flag
	expected := []byte{MsgReplyUtf8}
	// Add identifier
	expected = append(expected, id[:]...)

	// Add payload
	expected = append(expected, payload.Data...)

	actual := NewReplyMessage(id, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReplyMsgUtf16 tests the NewReplyMessage method using UTF16 payload encoding
func TestMsgNewReplyMsgUtf16(t *testing.T) {
	id := genRndMsgID()
	payload := Payload{
		Encoding: EncodingUtf16,
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

	actual := NewReplyMessage(id, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewReplyMsgUtf16CorruptPayload tests the NewReplyMessage method
// using UTF16 payload encoding
// passing corrupt data (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewReplyMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected NewReplyMessage to panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewReplyMessage(genRndMsgID(), Payload{
		Encoding: EncodingUtf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		Data: []byte("invalid"),
	})
}

// TestMsgNewSigMsgNameTooLong tests the NewSignalMessage method with a too long name
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

	NewSignalMessage(string(nameBuf), Payload{})
}

// TestMsgNewSigMsgBinary tests the NewSignalMessage method using the default binary encoding
func TestMsgNewSigMsgBinary(t *testing.T) {
	name := genRndName()
	payload := Payload{
		Encoding: EncodingBinary,
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

	actual := NewSignalMessage(name, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewSigMsgUtf8 tests the NewSignalMessage method using UTF8 encoding
func TestMsgNewSigMsgUtf8(t *testing.T) {
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf8,
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

	actual := NewSignalMessage(name, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewSigMsgUtf16 tests the NewSignalMessage method using UTF16 encoding
func TestMsgNewSigMsgUtf16(t *testing.T) {
	name := genRndName()
	payload := Payload{
		Encoding: EncodingUtf16,
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

	actual := NewSignalMessage(name, payload)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Binary results differ:\n%v\n%v", expected, actual)
	}
}

// TestMsgNewSigMsgUtf16CorruptPayload tests the NewSignalMessage method using UTF16 payload encoding
// passing corrupt data (length not divisible by 2 thus not UTF16 encoded)
func TestMsgNewSigMsgUtf16CorruptPayload(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected NewSignalMessage to panic due to corrupt UTF16 payload")
		} else {
			return
		}
	}()

	NewSignalMessage(genRndName(), Payload{
		Encoding: EncodingUtf16,
		// Payload is corrupt, only 7 bytes long, not power 2
		Data: []byte("invalid"),
	})
}
