package message

import (
	"testing"

	pld "github.com/qbeon/webwire-go/payload"
)

// TestRequiresReplyCloseSession tests the RequiresReply method
// with a session closure request message
func TestRequiresReplyCloseSession(t *testing.T) {
	msg := &Message{}
	if _, err := msg.Parse(NewEmptyRequestMessage(
		MsgCloseSession,
		genRndMsgIdentifier(),
	)); err != nil {
		t.Fatal(err)
	}

	if !msg.RequiresReply() {
		t.Fatalf(
			"Expected a session closure request message to require a reply",
		)
	}
}

// TestRequiresReplyRestoreSession tests the RequiresReply method
// with a session restoration request message
func TestRequiresReplyRestoreSession(t *testing.T) {
	msg := &Message{}
	if _, err := msg.Parse(NewNamelessRequestMessage(
		MsgRestoreSession,
		genRndMsgIdentifier(),
		[]byte("somesamplesessionkey"),
	)); err != nil {
		t.Fatal(err)
	}

	if !msg.RequiresReply() {
		t.Fatalf(
			"Expected a session restoration request message to require a reply",
		)
	}
}

// TestRequiresReplyRequestBinary tests the RequiresReply method
// with a binary request message
func TestRequiresReplyRequestBinary(t *testing.T) {
	msg := &Message{}
	if _, err := msg.Parse(NewRequestMessage(
		genRndMsgIdentifier(),
		"samplename",
		pld.Binary,
		[]byte("random payload data"),
	)); err != nil {
		t.Fatal(err)
	}

	if !msg.RequiresReply() {
		t.Fatalf("Expected a binary request message to require a reply")
	}
}

// TestRequiresReplyRequestUtf8 tests the RequiresReply method
// with a UTF8 encoded request message
func TestRequiresReplyRequestUtf8(t *testing.T) {
	msg := &Message{}
	if _, err := msg.Parse(NewRequestMessage(
		genRndMsgIdentifier(),
		"samplename",
		pld.Utf8,
		[]byte("random payload data"),
	)); err != nil {
		t.Fatal(err)
	}

	if !msg.RequiresReply() {
		t.Fatalf("Expected a UTF8 request message to require a reply")
	}
}

// TestRequiresReplyRequestUtf16 tests the RequiresReply method
// with a UTF16 encoded request message
func TestRequiresReplyRequestUtf16(t *testing.T) {
	msg := &Message{}
	if _, err := msg.Parse(NewRequestMessage(
		genRndMsgIdentifier(),
		"samplename",
		pld.Utf16,
		[]byte{'r', 0, 'a', 0, 'n', 0, 'd', 0, 'o', 0, 'm', 0},
	)); err != nil {
		t.Fatal(err)
	}

	if !msg.RequiresReply() {
		t.Fatalf("Expected a UTF16 request message to require a reply")
	}
}
