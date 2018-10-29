package message

import (
	"time"

	pld "github.com/qbeon/webwire-go/payload"
)

const (
	// MsgMinLenSignal represents the minimum length
	// of binary/UTF8 encoded signal messages.
	// binary/UTF8 signal message structure:
	//  1. message type (1 byte)
	//  2. name length flag (1 byte)
	//  3. name (n bytes, optional if name length flag is 0)
	//  4. payload (n bytes, at least 1 byte)
	MsgMinLenSignal = int(3)

	// MsgMinLenSignalUtf16 represents the minimum length
	// of UTF16 encoded signal messages.
	// UTF16 signal message structure:
	//  1. message type (1 byte)
	//  2. name length flag (1 byte)
	//  3. name (n bytes, optional if name length flag is 0)
	//  4. header padding (1 byte, required if name length flag is odd)
	//  5. payload (n bytes, at least 2 bytes)
	MsgMinLenSignalUtf16 = int(4)

	// MsgMinLenRequest represents the minimum length
	// of binary/UTF8 encoded request messages.
	// binary/UTF8 request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. name length flag (1 byte)
	//  4. name (from 0 to 255 bytes, optional if name length flag is 0)
	//  5. payload (n bytes, at least 1 byte or optional if name len > 0)
	MsgMinLenRequest = int(11)

	// MsgMinLenRequestUtf16 represents the minimum length
	// of UTF16 encoded request messages.
	// UTF16 request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. name length flag (1 byte)
	//  4. name (n bytes, optional if name length flag is 0)
	//  5. header padding (1 byte, required if name length flag is odd)
	//  6. payload (n bytes, at least 2 bytes)
	MsgMinLenRequestUtf16 = int(11)

	// MsgMinLenReply represents the minimum length
	// of binary/UTF8 encoded reply messages.
	// binary/UTF8 reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. payload (n bytes, optional or at least 1 byte)
	MsgMinLenReply = int(9)

	// MsgMinLenReplyUtf16 represents the minimum length
	// of UTF16 encoded reply messages.
	// UTF16 reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. header padding (1 byte)
	//  4. payload (n bytes, optional or at least 2 bytes)
	MsgMinLenReplyUtf16 = int(10)

	// MsgMinLenErrorReply represents the minimum length
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
	MsgMinLenErrorReply = int(11)

	// MsgMinLenRestoreSession represents the minimum length
	// of session restoration request messages.
	// Session restoration request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. session key (n bytes, 7-bit ASCII encoded, at least 1 byte)
	MsgMinLenRestoreSession = int(10)

	// MsgMinLenCloseSession represents the minimum length
	// of session destruction request messages.
	// Session destruction request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	MsgMinLenCloseSession = int(9)

	// MsgMinLenSessionCreated represents the minimum length
	// of session creation notification messages.
	// Session creation notification message structure:
	//  1. message type (1 byte)
	//  2. session key (n bytes, 7-bit ASCII encoded, at least 1 byte)
	MsgMinLenSessionCreated = int(2)

	// MsgMinLenSessionClosed represents the minimum length
	// of session creation notification messages.
	// Session destruction notification message structure:
	//  1. message type (1 byte)
	MsgMinLenSessionClosed = int(1)

	// MsgMinLenConf represents the minimum length
	// of an endpoint metadata message.
	//  1. message type (1 byte)
	//  2. major protocol version (1 byte)
	//  3. minor protocol version (1 byte)
	//  4. read timeout in milliseconds (4 byte)
	//  5. read buffer size in bytes (4 byte)
	//  6. write buffer size in bytes (4 byte)
	MsgMinLenConf = int(15)
)

const (
	// SERVER

	// MsgErrorReply is sent by the server
	// and represents an error-reply to a previously sent request
	MsgErrorReply = byte(0)

	// MsgReplyShutdown is sent by the server when a request is received
	// during server shutdown and can't therefore be processed
	MsgReplyShutdown = byte(1)

	// MsgInternalError is sent by the server if an unexpected internal error
	// arose during the processing of a request
	MsgInternalError = byte(2)

	// MsgSessionNotFound is sent by the server in response to an unfulfilled
	// session restoration request due to the session not being found
	MsgSessionNotFound = byte(3)

	// MsgMaxSessConnsReached is sent by the server in response to
	// an authentication request when the maximum number
	// of concurrent connections for a certain session was reached
	MsgMaxSessConnsReached = byte(4)

	// MsgSessionsDisabled is sent by the server in response to
	// a session restoration request
	// if sessions are disabled for the target server
	MsgSessionsDisabled = byte(5)

	// MsgReplyProtocolError is sent by the server in response to an invalid
	// message violating the protocol
	MsgReplyProtocolError = byte(6)

	// MsgSessionCreated is sent by the server
	// to notify the client about the session creation
	MsgSessionCreated = byte(21)

	// MsgSessionClosed is sent by the server
	// to notify the client about the session destruction
	MsgSessionClosed = byte(22)

	// MsgConf is sent by the server right after the handshake and includes
	// the server exposed configurations
	MsgConf = byte(23)

	// CLIENT

	// MsgCloseSession is sent by the client
	// and represents a request for the destruction
	// of the currently active session
	MsgCloseSession = byte(31)

	// MsgRestoreSession is sent by the client
	// to request session restoration
	MsgRestoreSession = byte(32)

	// MsgHeartbeat is sent by the client to acknowledge the server about the
	// activity of the connection to prevent it from shutting the connection
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

// ServerConfiguration represents the MsgConf payload data
type ServerConfiguration struct {
	MajorProtocolVersion byte
	MinorProtocolVersion byte
	ReadTimeout          time.Duration
	ReadBufferSize       uint32
	WriteBufferSize      uint32
}

// Message represents a WebWire protocol message
type Message struct {
	Type       byte
	Identifier [8]byte
	Name       string
	Payload    pld.Payload

	// ServerConfiguration is only initialized for MsgConf type messages
	ServerConfiguration ServerConfiguration
}

// RequiresReply returns true if a message of this type requires a reply,
// otherwise returns false.
func (msg *Message) RequiresReply() bool {
	switch msg.Type {
	case MsgCloseSession:
		fallthrough
	case MsgRestoreSession:
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
