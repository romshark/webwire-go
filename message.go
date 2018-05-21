package webwire

import (
	"fmt"
)

const (
	// MsgMinLenSignal represents the minimum binary/UTF8 encoded signal message length.
	// binary/UTF8 signal message structure:
	//  1. message type (1 byte)
	//  2. name length flag (1 byte)
	//  3. name (n bytes, optional if name length flag is 0)
	//  4. payload (n bytes, at least 1 byte)
	MsgMinLenSignal = int(3)

	// MsgMinLenSignalUtf16 represents the minimum UTF16 encoded signal message length.
	// UTF16 signal message structure:
	//  1. message type (1 byte)
	//  2. name length flag (1 byte)
	//  3. name (n bytes, optional if name length flag is 0)
	//  4. header padding (1 byte, required if name length flag is odd)
	//  5. payload (n bytes, at least 2 bytes)
	MsgMinLenSignalUtf16 = int(4)

	// MsgMinLenRequest represents the minimum binary/UTF8 encoded request message length.
	// binary/UTF8 request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. name length flag (1 byte)
	//  4. name (from 0 to 255 bytes, optional if name length flag is 0)
	//  5. payload (n bytes, at least 1 byte or optional if name len > 0)
	MsgMinLenRequest = int(11)

	// MsgMinLenRequestUtf16 represents the minimum UTF16 encoded request message length.
	// UTF16 request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. name length flag (1 byte)
	//  4. name (n bytes, optional if name length flag is 0)
	//  5. header padding (1 byte, required if name length flag is odd)
	//  6. payload (n bytes, at least 2 bytes)
	MsgMinLenRequestUtf16 = int(11)

	// MsgMinLenReply represents the minimum binary/UTF8 encoded reply message length.
	// binary/UTF8 reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. payload (n bytes, optional or at least 1 byte)
	MsgMinLenReply = int(9)

	// MsgMinLenReplyUtf16 represents the minimum UTF16 encoded reply message length
	// UTF16 reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. header padding (1 byte)
	//  4. payload (n bytes, optional or at least 2 bytes)
	MsgMinLenReplyUtf16 = int(10)

	// MsgMinLenErrorReply represents the minimum error reply message length
	// Error reply message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. error code length flag (1 byte, cannot be 0)
	//  4. error code (from 1 to 255 bytes, length must correspond to the length flag)
	//  5. error message (n bytes, UTF8 encoded, optional)
	MsgMinLenErrorReply = int(11)

	// MsgMinLenRestoreSession represents the minimum session restoration request message length
	// Session restoration request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	//  3. session key (n bytes, 7-bit ASCII encoded, at least 1 byte)
	MsgMinLenRestoreSession = int(10)

	// MsgMinLenCloseSession represents the minimum session destruction request message length
	// Session destruction request message structure:
	//  1. message type (1 byte)
	//  2. message id (8 bytes)
	MsgMinLenCloseSession = int(9)

	// MsgMinLenSessionCreated represents the minimum session creation notification message length
	// Session creation notification message structure:
	//  1. message type (1 byte)
	//  2. session key (n bytes, 7-bit ASCII encoded, at least 1 byte)
	MsgMinLenSessionCreated = int(2)

	// MsgMinLenSessionClosed represents the minimum session creation notification message length
	// Session destruction notification message structure:
	//  1. message type (1 byte)
	MsgMinLenSessionClosed = int(1)
)

const (
	// SERVER

	// MsgErrorReply is sent by the server
	// and represents an error-reply to a previously sent request
	MsgErrorReply = byte(0)

	// MsgReplyShutdown is sent by the server when a request is received during server shutdown
	// and can't therefore be processed
	MsgReplyShutdown = byte(1)

	// MsgInternalError is sent by the server if an unexpected internal error arose during
	// the processing of a request
	MsgInternalError = byte(2)

	// MsgSessionNotFound is sent by the server in response to an unfilfilled session restoration
	// request due to the session not being found
	MsgSessionNotFound = byte(3)

	// MsgMaxSessConnsReached is sent by the server in response to an authentication request
	// when the maximum number of concurrent connections for a certain session was reached
	MsgMaxSessConnsReached = byte(4)

	// MsgSessionsDisabled is sent by the server in response to a session restoration request
	// if sessions are disabled for the target server
	MsgSessionsDisabled = byte(5)

	// MsgSessionCreated is sent by the server
	// to notify the client about the session creation
	MsgSessionCreated = byte(21)

	// MsgSessionClosed is sent by the server
	// to notify the client about the session destruction
	MsgSessionClosed = byte(22)

	// CLIENT

	// MsgCloseSession is sent by the client
	// and represents a request for the destruction of the currently active session
	MsgCloseSession = byte(31)

	// MsgRestoreSession is sent by the client
	// to request session restoration
	MsgRestoreSession = byte(32)

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

// Message represents a WebWire protocol message
type Message struct {
	msgType byte
	id      [8]byte
	Name    string
	Payload Payload
}

// MessageType returns the type of the message
func (msg *Message) MessageType() byte {
	return msg.msgType
}

// Identifier returns the message identifier
func (msg *Message) Identifier() [8]byte {
	return msg.id
}

// RequiresResponse returns true if this type of message
// requires a response to be sent in return
func (msg *Message) RequiresResponse() bool {
	switch msg.msgType {
	case MsgRequestBinary:
		return true
	case MsgRequestUtf8:
		return true
	case MsgRequestUtf16:
		return true
	case MsgRestoreSession:
		return true
	case MsgCloseSession:
		return true
	}
	return false
}

// NewSignalMessage composes a new named signal message and returns its binary representation
func NewSignalMessage(name string, payload Payload) (msg []byte) {
	if len(name) > 255 {
		panic(fmt.Errorf("Unsupported request message name length: %d", len(name)))
	}

	// Verify payload data validity in case of UTF16 encoding
	if payload.Encoding == EncodingUtf16 && len(payload.Data)%2 != 0 {
		panic(fmt.Errorf("Invalid UTF16 signal payload data length: %d", len(payload.Data)))
	}

	// Determine total message length
	messageSize := 2 + len(name) + len(payload.Data)

	// Check if a header padding is necessary.
	// A padding is necessary if the payload is UTF16 encoded
	// but not properly aligned due to a header length not divisible by 2
	headerPadding := false
	if payload.Encoding == EncodingUtf16 && len(name)%2 != 0 {
		headerPadding = true
		messageSize++
	}

	msg = make([]byte, messageSize)

	// Write message type flag
	sigType := MsgSignalBinary
	switch payload.Encoding {
	case EncodingUtf8:
		sigType = MsgSignalUtf8
	case EncodingUtf16:
		sigType = MsgSignalUtf16
	}
	msg[0] = sigType

	// Write name length flag
	msg[1] = byte(len(name))

	// Write name
	for i := 0; i < len(name); i++ {
		char := name[i]
		if char < 32 || char > 126 {
			panic(fmt.Errorf("Unsupported character in request name: %s", string(char)))
		}
		msg[2+i] = char
	}

	// Write header padding byte if the payload requires proper alignment
	payloadOffset := 2 + len(name)
	if headerPadding {
		msg[payloadOffset] = 0
		payloadOffset++
	}

	// Write payload
	for i := 0; i < len(payload.Data); i++ {
		msg[payloadOffset+i] = payload.Data[i]
	}

	return msg
}

// NewRequestMessage composes a new named request message and returns its binary representation
func NewRequestMessage(id [8]byte, name string, payload Payload) (msg []byte) {
	// Require either a name, or a payload or both, but don't allow none
	if len(name) < 1 && len(payload.Data) < 1 {
		panic(fmt.Errorf(
			"Request message requires either a name, or a payload, or both",
		))
	}

	// Cap name length at 255 bytes
	if len(name) > 255 {
		panic(fmt.Errorf("Unsupported request message name length: %d", len(name)))
	}

	// Verify payload data validity in case of UTF16 encoding
	if payload.Encoding == EncodingUtf16 && len(payload.Data)%2 != 0 {
		panic(fmt.Errorf("Invalid UTF16 request payload data length: %d", len(payload.Data)))
	}

	// Determine total message length
	messageSize := 10 + len(name) + len(payload.Data)

	// Check if a header padding is necessary.
	// A padding is necessary if the payload is UTF16 encoded
	// but not properly aligned due to a header length not divisible by 2
	headerPadding := false
	if payload.Encoding == EncodingUtf16 && len(name)%2 != 0 {
		headerPadding = true
		messageSize++
	}

	msg = make([]byte, messageSize)

	// Write message type flag
	reqType := MsgRequestBinary
	switch payload.Encoding {
	case EncodingUtf8:
		reqType = MsgRequestUtf8
	case EncodingUtf16:
		reqType = MsgRequestUtf16
	}
	msg[0] = reqType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = id[i]
	}

	// Write name length flag
	msg[9] = byte(len(name))

	// Write name
	for i := 0; i < len(name); i++ {
		char := name[i]
		if char < 32 || char > 126 {
			panic(fmt.Errorf("Unsupported character in request name: %s", string(char)))
		}
		msg[10+i] = char
	}

	// Write header padding byte if the payload requires proper alignment
	payloadOffset := 10 + len(name)
	if headerPadding {
		msg[payloadOffset] = 0
		payloadOffset++
	}

	// Write payload
	for i := 0; i < len(payload.Data); i++ {
		msg[payloadOffset+i] = payload.Data[i]
	}

	return msg
}

// NewReplyMessage composes a new reply message and returns its binary representation
func NewReplyMessage(requestID [8]byte, payload Payload) (msg []byte) {
	// Determine total message length
	messageSize := 9 + len(payload.Data)

	// Verify payload data validity in case of UTF16 encoding
	if payload.Encoding == EncodingUtf16 && len(payload.Data)%2 != 0 {
		panic(fmt.Errorf("Invalid UTF16 reply payload data length: %d", len(payload.Data)))
	}

	// Check if a header padding is necessary.
	// A padding is necessary if the payload is UTF16 encoded
	// but not properly aligned due to a header length not divisible by 2
	headerPadding := false
	if payload.Encoding == EncodingUtf16 {
		headerPadding = true
		messageSize++
	}

	msg = make([]byte, messageSize)

	// Write message type flag
	reqType := MsgReplyBinary
	switch payload.Encoding {
	case EncodingUtf8:
		reqType = MsgReplyUtf8
	case EncodingUtf16:
		reqType = MsgReplyUtf16
	}
	msg[0] = reqType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = requestID[i]
	}

	// Write header padding byte if the payload requires proper alignment
	payloadOffset := 9
	if headerPadding {
		msg[payloadOffset] = 0
		payloadOffset++
	}

	// Write payload
	for i := 0; i < len(payload.Data); i++ {
		msg[payloadOffset+i] = payload.Data[i]
	}

	return msg
}

// NewErrorReplyMessage composes a new error reply message
// and returns its binary representation
func NewErrorReplyMessage(
	requestIdent [8]byte,
	code,
	message string,
) (msg []byte) {
	if len(code) < 1 {
		panic(fmt.Errorf(
			"Missing error code while creating a new error reply message",
		))
	} else if len(code) > 255 {
		panic(fmt.Errorf(
			"Invalid error code while creating a new error reply message,"+
				"too long (%d)",
			len(code),
		))
	}

	// Determine total message length
	msg = make([]byte, 10+len(code)+len(message))

	// Write message type flag
	msg[0] = MsgErrorReply

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = requestIdent[i]
	}

	// Write code length flag
	msg[9] = byte(len(code))

	// Write error code
	for i := 0; i < len(code); i++ {
		char := code[i]
		if char < 32 || char > 126 {
			panic(fmt.Errorf(
				"Unsupported character in reply error - error code: %s",
				string(char),
			))
		}
		msg[10+i] = code[i]
	}

	errMessageOffset := 10 + len(code)

	// Write error message
	for i := 0; i < len(message); i++ {
		msg[errMessageOffset+i] = message[i]
	}

	return msg
}

// NewNamelessRequestMessage composes a new nameless (initially without a name) request message
// and returns its binary representation
func NewNamelessRequestMessage(reqType byte, id [8]byte, payload []byte) (msg []byte) {
	// 9 byte header + n bytes payload
	msg = make([]byte, 9+len(payload))

	// Write message type flag
	msg[0] = reqType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = id[i]
	}

	// Write payload
	for i := 0; i < len(payload); i++ {
		msg[9+i] = payload[i]
	}

	return msg
}

// NewEmptyRequestMessage composes a new request message consisting only of the type and identifier
// and returns its binary representation
func NewEmptyRequestMessage(msgType byte, id [8]byte) (msg []byte) {
	msg = make([]byte, 9)

	// Write message type flag
	msg[0] = msgType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = id[i]
	}

	return msg
}

// NewSpecialRequestReplyMessage composes a new special request reply message
func NewSpecialRequestReplyMessage(msgType byte, reqIdent [8]byte) []byte {
	switch msgType {
	case MsgInternalError:
		break
	case MsgMaxSessConnsReached:
		break
	case MsgSessionNotFound:
		break
	case MsgSessionsDisabled:
		break
	case MsgReplyShutdown:
		break
	default:
		panic(fmt.Errorf(
			"Message type (%d) doesn't represent a special reply message",
			msgType,
		))
	}

	msg := make([]byte, 9)

	// Write message type flag
	msg[0] = msgType

	// Write request identifier
	for i := 0; i < 8; i++ {
		msg[1+i] = reqIdent[i]
	}

	return msg
}

func (msg *Message) parseSignal(message []byte) error {
	if len(message) < MsgMinLenSignal {
		return fmt.Errorf("Invalid signal message, too short")
	}

	// Read name length
	nameLen := int(byte(message[1:2][0]))
	payloadOffset := 2 + nameLen

	// Verify total message size to prevent segmentation faults caused by inconsistent flags,
	// this could happen if the specified name length doesn't correspond to the actual name length
	if len(message) < MsgMinLenSignal+nameLen {
		return fmt.Errorf(
			"Invalid signal message, too short for full name (%d) and the minimum payload (1)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.Name = string(message[2:payloadOffset])
		msg.Payload = Payload{
			Data: message[payloadOffset:],
		}
	} else {
		// No name present, just payload
		msg.Payload = Payload{
			Data: message[2:],
		}
	}
	return nil
}

func (msg *Message) parseSignalUtf16(message []byte) error {
	if len(message) < MsgMinLenSignalUtf16 {
		return fmt.Errorf("Invalid signal message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded signal message (probably missing header padding)",
		)
	}

	// Read name length
	nameLen := int(byte(message[1:2][0]))

	// Determine minimum required message length
	minMsgSize := MsgMinLenSignalUtf16 + nameLen
	payloadOffset := 2 + nameLen

	// Check whether a name padding byte is to be expected
	if nameLen%2 != 0 {
		minMsgSize++
		payloadOffset++
	}

	// Verify total message size to prevent segmentation faults caused by inconsistent flags,
	// this could happen if the specified name length doesn't correspond to the actual name length
	if len(message) < minMsgSize {
		return fmt.Errorf(
			"Invalid signal message, too short for full name (%d) and the minimum payload (2)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.Name = string(message[2 : 2+nameLen])
		msg.Payload = Payload{
			Data: message[payloadOffset:],
		}
	} else {
		// No name present, just payload
		msg.Payload = Payload{
			Data: message[2:],
		}
	}
	return nil
}

func (msg *Message) parseRequest(message []byte) error {
	if len(message) < MsgMinLenRequest {
		return fmt.Errorf("Invalid request message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Read name length
	nameLen := int(byte(message[9:10][0]))
	payloadOffset := 10 + nameLen

	// Verify total message size to prevent segmentation faults caused
	// by inconsistent flags, this could happen if the specified name length
	// doesn't correspond to the actual name length
	if nameLen > 0 {
		// Subtract one to not require the payload but at least the name
		if len(message) < MsgMinLenRequest+nameLen-1 {
			return fmt.Errorf(
				"Invalid request message, too short for full name (%d)",
				nameLen,
			)
		}

		// Take name into account
		msg.Name = string(message[10 : 10+nameLen])

		// Read payload if any
		if len(message) > MsgMinLenRequest+nameLen-1 {
			msg.Payload = Payload{
				Data: message[payloadOffset:],
			}
		}
	} else {
		// No name present, expect just the payload to be in place
		msg.Payload = Payload{
			Data: message[10:],
		}
	}

	return nil
}

func (msg *Message) parseRequestUtf16(message []byte) error {
	if len(message) < MsgMinLenRequestUtf16 {
		return fmt.Errorf("Invalid request message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded request message (probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Read name length
	nameLen := int(byte(message[9:10][0]))

	// Determine minimum required message length.
	// There's at least a 10 byte header and a 2 byte payload expected
	minRequiredMsgSize := 12
	if nameLen > 0 {
		// ...unless a name is given, in which case the payload isn't required
		minRequiredMsgSize = 10 + nameLen
	}

	// A header padding byte is only expected, when there's a payload
	// beyond the name. It's not required if there's just the header and a name
	payloadOffset := 10 + nameLen
	if len(message) > payloadOffset && nameLen%2 != 0 {
		minRequiredMsgSize++
		payloadOffset++
	}

	// Verify total message size to prevent segmentation faults caused
	// by inconsistent flags, this could happen if the specified name length
	// doesn't correspond to the actual name length
	if nameLen > 0 {
		if len(message) < minRequiredMsgSize {
			return fmt.Errorf(
				"Invalid request message, too short for full name (%d)",
				nameLen,
			)
		}

		// Take name into account
		msg.Name = string(message[10 : 10+nameLen])

		// Read payload if any
		if len(message) > minRequiredMsgSize {
			msg.Payload = Payload{
				Data: message[payloadOffset:],
			}
		}
	} else {
		// No name present, just payload
		msg.Payload = Payload{
			Data: message[10:],
		}
	}

	return nil
}

func (msg *Message) parseReply(message []byte) error {
	if len(message) < MsgMinLenReply {
		return fmt.Errorf("Invalid reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Skip payload if there's none
	if len(message) == MsgMinLenReply {
		return nil
	}

	// Read payload
	msg.Payload = Payload{
		Data: message[9:],
	}
	return nil
}

func (msg *Message) parseReplyUtf16(message []byte) error {
	if len(message) < MsgMinLenReplyUtf16 {
		return fmt.Errorf("Invalid UTF16 reply message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unaligned UTF16 encoded reply message (probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Skip payload if there's none
	if len(message) == MsgMinLenReplyUtf16 {
		msg.Payload = Payload{
			Encoding: EncodingUtf16,
		}
		return nil
	}

	// Read payload
	msg.Payload = Payload{
		// Take header padding byte into account
		Data: message[10:],
	}
	return nil
}

// parseErrorReply parses the given message assuming it's an error reply message
// parsing the error code into the name field and the UTF8 encoded error message into the payload
func (msg *Message) parseErrorReply(message []byte) error {
	if len(message) < MsgMinLenErrorReply {
		return fmt.Errorf("Invalid error reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Read error code length flag
	errCodeLen := int(byte(message[9:10][0]))
	errMessageOffset := 10 + errCodeLen

	// Verify error code length (must be at least 1 character long)
	if errCodeLen < 1 {
		return fmt.Errorf("Invalid error reply message, error code length flag is zero")
	}

	// Verify total message size to prevent segmentation faults caused by inconsistent flags,
	// this could happen if the specified error code length
	// doesn't correspond to the actual length of the provided error code.
	// Subtract 1 character already taken into account by MsgMinLenErrorReply
	if len(message) < MsgMinLenErrorReply+errCodeLen-1 {
		return fmt.Errorf(
			"Invalid error reply message, too short for specified code length (%d)",
			errCodeLen,
		)
	}

	// Read UTF8 encoded error message into the payload
	msg.Name = string(message[10 : 10+errCodeLen])
	msg.Payload = Payload{
		Encoding: EncodingUtf8,
		Data:     message[errMessageOffset:],
	}
	return nil
}

func (msg *Message) parseRestoreSession(message []byte) error {
	if len(message) < MsgMinLenRestoreSession {
		return fmt.Errorf("Invalid session restoration request message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Read payload
	msg.Payload = Payload{
		Data: message[9:],
	}
	return nil
}

func (msg *Message) parseCloseSession(message []byte) error {
	if len(message) != MsgMinLenCloseSession {
		return fmt.Errorf("Invalid session destruction request message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	return nil
}

func (msg *Message) parseSessionCreated(message []byte) error {
	if len(message) < MsgMinLenSessionCreated {
		return fmt.Errorf("Invalid session creation notification message, too short")
	}

	msg.Payload = Payload{
		Data: message[1:],
	}
	return nil
}

func (msg *Message) parseSessionClosed(message []byte) error {
	if len(message) != MsgMinLenSessionClosed {
		return fmt.Errorf("Invalid session closure notification message, too short")
	}
	return nil
}

func (msg *Message) parseSpecialReplyMessage(message []byte) error {
	if len(message) < 9 {
		return fmt.Errorf("Invalid special reply message, too short")
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	return nil
}

// Parse tries to parse the message from a byte slice
func (msg *Message) Parse(message []byte) (err error) {
	if len(message) < 1 {
		return fmt.Errorf("Invalid message, too short")
	}
	var payloadEncoding PayloadEncoding
	msgType := message[0:1][0]

	switch msgType {

	// Request error reply message
	case MsgErrorReply:
		err = msg.parseErrorReply(message)

	// Session creation notification message
	case MsgSessionCreated:
		err = msg.parseSessionCreated(message)

	// Session closure notification message
	case MsgSessionClosed:
		err = msg.parseSessionClosed(message)

	// Session destruction request message
	case MsgCloseSession:
		err = msg.parseCloseSession(message)

	// Signal messages
	case MsgSignalBinary:
		payloadEncoding = EncodingBinary
		err = msg.parseSignal(message)
	case MsgSignalUtf8:
		payloadEncoding = EncodingUtf8
		err = msg.parseSignal(message)
	case MsgSignalUtf16:
		payloadEncoding = EncodingUtf16
		err = msg.parseSignalUtf16(message)

	// Request messages
	case MsgRequestBinary:
		payloadEncoding = EncodingBinary
		err = msg.parseRequest(message)
	case MsgRequestUtf8:
		payloadEncoding = EncodingUtf8
		err = msg.parseRequest(message)
	case MsgRequestUtf16:
		payloadEncoding = EncodingUtf16
		err = msg.parseRequestUtf16(message)

	// Reply messages
	case MsgReplyBinary:
		payloadEncoding = EncodingBinary
		err = msg.parseReply(message)
	case MsgReplyUtf8:
		payloadEncoding = EncodingUtf8
		err = msg.parseReply(message)
	case MsgReplyUtf16:
		payloadEncoding = EncodingUtf16
		err = msg.parseReplyUtf16(message)

	// Session restoration request message
	case MsgRestoreSession:
		err = msg.parseRestoreSession(message)

	// Special reply messages
	case MsgReplyShutdown:
		err = msg.parseSpecialReplyMessage(message)
	case MsgInternalError:
		err = msg.parseSpecialReplyMessage(message)
	case MsgSessionNotFound:
		err = msg.parseSpecialReplyMessage(message)
	case MsgMaxSessConnsReached:
		err = msg.parseSpecialReplyMessage(message)
	case MsgSessionsDisabled:
		err = msg.parseSpecialReplyMessage(message)

	// Ignore messages of invalid message type
	default:
		return fmt.Errorf("Invalid message type (%d)", msgType)
	}

	msg.msgType = msgType
	msg.Payload.Encoding = payloadEncoding
	return err
}
