package message

import (
	"time"

	pld "github.com/qbeon/webwire-go/payload"
)

const (
	// MinLenSignal represents the minimum length
	// of binary/UTF8 encoded signal messages.
	// binary/UTF8 signal message structure:
	//  1. message type (1 byte)
	//  2. name length flag (1 byte)
	//  3. name (n bytes, optional if name length flag is 0)
	//  4. payload (n bytes, at least 1 byte)
	MinLenSignal = int(3)

	// MinLenSignalUtf16 represents the minimum length
	// of UTF16 encoded signal messages.
	// UTF16 signal message structure:
	//  1. message type (1 byte)
	//  2. name length flag (1 byte)
	//  3. name (n bytes, optional if name length flag is 0)
	//  4. header padding (1 byte, required if name length flag is odd)
	//  5. payload (n bytes, at least 2 bytes)
	MinLenSignalUtf16 = int(4)

	// MinLenRequest represents the minimum length
	// of binary/UTF8 encoded request messages.
	// binary/UTF8 request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. name length flag (1 byte)
	//  4. name (from 0 to 255 bytes, optional if name length flag is 0)
	//  5. payload (n bytes, at least 1 byte or optional if name len > 0)
	MinLenRequest = int(11)

	// MinLenRequestUtf16 represents the minimum length
	// of UTF16 encoded request messages.
	// UTF16 request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. name length flag (1 byte)
	//  4. name (n bytes, optional if name length flag is 0)
	//  5. header padding (1 byte, required if name length flag is odd)
	//  6. payload (n bytes, at least 2 bytes)
	MinLenRequestUtf16 = int(11)

	// MinLenReply represents the minimum length
	// of binary/UTF8 encoded reply messages.
	// binary/UTF8 reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. payload (n bytes, optional or at least 1 byte)
	MinLenReply = int(9)

	// MinLenReplyUtf16 represents the minimum length
	// of UTF16 encoded reply messages.
	// UTF16 reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. header padding (1 byte)
	//  4. payload (n bytes, optional or at least 2 bytes)
	MinLenReplyUtf16 = int(10)

	// MinLenReplyError represents the minimum length
	// of error reply messages.
	// Error reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. error code length flag (1 byte, cannot be 0)
	//  4. error code (
	//    from 1 to 255 bytes,
	//    length must correspond to the length flag
	//  )
	//  5. error message (n bytes, UTF8 encoded, optional)
	MinLenReplyError = int(11)

	// MinLenRequestRestoreSession represents the minimum length
	// of session restoration request messages.
	// Session restoration request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. session key (n bytes, 7-bit ASCII encoded, at least 1 byte)
	MinLenRequestRestoreSession = int(10)

	// MinLenDoCloseSession represents the minimum length
	// of session destruction request messages.
	// Session destruction request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	MinLenDoCloseSession = int(9)

	// MinLenNotifySessionCreated represents the minimum length
	// of session creation notification messages.
	// Session creation notification message structure:
	//  1. message type (1 byte)
	//  2. session key (n bytes, 7-bit ASCII encoded, at least 1 byte)
	MinLenNotifySessionCreated = int(2)

	// MinLenNotifySessionClosed represents the minimum length
	// of session creation notification messages.
	// Session destruction notification message structure:
	//  1. message type (1 byte)
	MinLenNotifySessionClosed = int(1)

	// MinLenAcceptConf represents the minimum length
	// of an endpoint metadata message.
	//  1. message type (1 byte)
	//  2. major protocol version (1 byte)
	//  3. minor protocol version (1 byte)
	//  4. read timeout in milliseconds (4 byte)
	//  5. message buffer size in bytes (4 byte)
	//  6. sub-protocol name (0+ bytes)
	MinLenAcceptConf = int(11)
)

const (
	// SERVER

	// MsgReplyError is a request reply sent only by the server and represents
	// an error-reply to a previously sent request
	MsgReplyError = byte(0)

	// MsgReplyShutdown is a request reply sent only by the server when a
	// request is received during server shutdown and can't therefore be
	// processed
	MsgReplyShutdown = byte(1)

	// MsgReplyInternalError is a request reply sent only by the server if an
	// unexpected internal error arose during the processing of a request
	MsgReplyInternalError = byte(2)

	// MsgReplySessionNotFound is a session restoration request reply sent only
	// by the server when the requested session was not found
	MsgReplySessionNotFound = byte(3)

	// MsgReplyMaxSessConnsReached is session restoration request reply sent
	// only by the server when the maximum number of concurrent connections for
	// a the requested session was reached
	MsgReplyMaxSessConnsReached = byte(4)

	// MsgReplySessionsDisabled is session restoration request reply sent only
	// by the server when sessions are disabled
	MsgReplySessionsDisabled = byte(5)

	// MsgNotifySessionCreated is a notification signal sent only by the server
	// to notify the client about the creation of a session
	MsgNotifySessionCreated = byte(21)

	// MsgNotifySessionClosed is a notification signal sent only by the server
	// to notify the client about the closure of the currently active session
	MsgNotifySessionClosed = byte(22)

	// MsgAcceptConf is a connection approval push-message sent only by the
	// server right after the handshake and includes the server configurations
	MsgAcceptConf = byte(23)

	// CLIENT

	// MsgRequestCloseSession is session closure command sent only by the client to
	// make the server close the currently active session
	MsgRequestCloseSession = byte(31)

	// MsgRequestRestoreSession is a session restoration request sent only by
	// the client
	MsgRequestRestoreSession = byte(32)

	// MsgHeartbeat is sent only by the client to acknowledge the server about
	// the activity of the connection to prevent it from shutting the connection
	// down on read timeout
	MsgHeartbeat = byte(33)

	// SIGNAL

	// Signals are sent by both the client and the server
	// and represents a one-way signal message that doesn't require a reply

	// MsgSignalBinary represents a signal with binary payload
	MsgSignalBinary = byte(63)

	// MsgSignalUtf8 represents a signal with UTF8 encoded payload
	MsgSignalUtf8 = byte(64)

	// MsgSignalUtf16 represents a signal with UTF16 encoded payload
	MsgSignalUtf16 = byte(65)

	// REQUEST
	// Requests are sent by the client
	// and represents a roundtrip to the server requiring a reply

	// MsgRequestBinary represents a request with binary payload
	MsgRequestBinary = byte(127)

	// MsgRequestUtf8 represents a request with a UTF8 encoded payload
	MsgRequestUtf8 = byte(128)

	// MsgRequestUtf16 represents a request with a UTF16 encoded payload
	MsgRequestUtf16 = byte(129)

	// REPLY
	// Replies are sent by the server
	// and represent a reply to a previously sent request

	// MsgReplyBinary represents a reply with a binary payload
	MsgReplyBinary = byte(191)

	// MsgReplyUtf8 represents a reply with a UTF8 encoded payload
	MsgReplyUtf8 = byte(192)

	// MsgReplyUtf16 represents a reply with a UTF16 encoded payload
	MsgReplyUtf16 = byte(193)
)

// ServerConfiguration represents the MsgAcceptConf payload data
type ServerConfiguration struct {
	MajorProtocolVersion byte
	MinorProtocolVersion byte
	SubProtocolName      []byte
	ReadTimeout          time.Duration
	MessageBufferSize    uint32
}

// Message represents a non-thread-safe WebWire protocol message
type Message struct {
	MsgBuffer          Buffer
	MsgType            byte
	MsgIdentifier      [8]byte
	MsgIdentifierBytes []byte
	MsgName            []byte
	MsgPayload         pld.Payload

	// ServerConfiguration is only initialized for MsgAcceptConf type messages
	ServerConfiguration ServerConfiguration

	onClose func()
}

// RequiresReply returns true if a message of this type requires a reply,
// otherwise returns false.
func (msg *Message) RequiresReply() bool {
	switch msg.MsgType {
	case MsgRequestCloseSession:
		fallthrough
	case MsgRequestRestoreSession:
		fallthrough
	case MsgRequestBinary:
		fallthrough
	case MsgRequestUtf8:
		fallthrough
	case MsgRequestUtf16:
		return true
	}
	return false
}

// Identifier implements the Message interface
func (msg *Message) Identifier() [8]byte {
	if msg.MsgBuffer.IsEmpty() {
		panic("read after close")
	}
	return msg.MsgIdentifier
}

// Name implements the Message interface
func (msg *Message) Name() []byte {
	if msg.MsgBuffer.IsEmpty() {
		panic("read after close")
	}
	return msg.MsgName
}

// PayloadEncoding implements the Message interface
func (msg *Message) PayloadEncoding() pld.Encoding {
	if msg.MsgBuffer.IsEmpty() {
		panic("read after close")
	}
	return msg.MsgPayload.Encoding
}

// Payload implements the Message interface
func (msg *Message) Payload() []byte {
	if msg.MsgBuffer.IsEmpty() {
		panic("read after close")
	}
	return msg.MsgPayload.Data
}

// PayloadUtf8 implements the Message interface
func (msg *Message) PayloadUtf8() ([]byte, error) {
	if msg.MsgBuffer.IsEmpty() {
		panic("read after close")
	}
	return msg.MsgPayload.Utf8()
}

// Close implements the Message interface
func (msg *Message) Close() {
	if msg.MsgBuffer.IsEmpty() {
		return
	}

	msg.MsgBuffer.Close()
	msg.MsgType = 0
	msg.MsgIdentifier = [8]byte{}
	msg.MsgIdentifierBytes = nil
	msg.MsgName = nil
	msg.MsgPayload = pld.Payload{}
	msg.ServerConfiguration = ServerConfiguration{}

	// Call closure callback
	msg.onClose()
}
