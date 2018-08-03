package message

import (
	"reflect"
	"testing"

	pld "github.com/qbeon/webwire-go/payload"
)

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
