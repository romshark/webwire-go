package webwire

import (
	"encoding/json"
	"fmt"
)

// ContextKey represents the identifiers of objects passed to the handlers context
type ContextKey int

const (
	// Msg identifies the message object stored in the handler context
	Msg ContextKey = iota
)

// Error represents an error returned in case of a wrong request
type Error struct {
	Code    string `json:"c"`
	Message string `json:"m,omitempty"`
}

const (
	// MsgMinLenSignal represents the minimum binary/UTF8 encoded signal message length
	MsgMinLenSignal = int(3)

	// MsgMinLenSignalUtf16 represents the minimum UTF16 encoded signal message length
	MsgMinLenSignalUtf16 = int(4)

	// MsgMinLenRequest represents the minimum binary/UTF8 encoded request message length
	MsgMinLenRequest = int(11)

	// MsgMinLenRequestUtf16 represents the minimum UTF16 encoded request message length
	MsgMinLenRequestUtf16 = int(12)

	// MsgMinLenReply represents the minimum binary/UTF8 encoded reply message length
	MsgMinLenReply = int(10)

	// MsgMinLenReplyUtf16 represents the minimum UTF16 encoded reply message length
	MsgMinLenReplyUtf16 = int(12)

	// MsgMinLenErrorReply represents the minimum error reply message length
	MsgMinLenErrorReply = int(10)

	// MsgMinLenRestoreSession represents the minimum session restoration request message length
	MsgMinLenRestoreSession = int(10)

	// MsgMinLenCloseSession represents the minimum session destruction request message length
	MsgMinLenCloseSession = int(9)

	// MsgMinLenSessionCreated represents the minimum session creation notification message length
	MsgMinLenSessionCreated = int(2)

	// MsgMinLenSessionClosed represents the minimum session creation notification message length
	MsgMinLenSessionClosed = int(1)
)

// PayloadEncoding represents the type of encoding of the message payload
type PayloadEncoding int

const (
	// EncodingBinary represents unencoded binary data
	EncodingBinary PayloadEncoding = iota

	// EncodingUtf8 represents UTF8 encoding
	EncodingUtf8

	// EncodingUtf16 represents UTF16 encoding
	EncodingUtf16
)

const (
	// SERVER/CLIENT

	// MsgErrorReply is sent by the server
	// and represents an error-reply to a previously sent request
	MsgErrorReply = byte(0)

	// SERVER

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

// Payload represents an encoded message payload
type Payload struct {
	Encoding PayloadEncoding
	Data     []byte
}

// Message represents a WebWire protocol message
type Message struct {
	fulfill func(reply Payload)
	fail    func(Error)

	msgType byte
	id      [8]byte

	Name    string
	Payload Payload
	Client  *Client
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
	// but not properly alligned due to a header length not divisible by 2
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
	// but not properly alligned due to a header length not divisible by 2
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
	// but not properly alligned due to a header length not divisible by 2
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

func (msg *Message) parseSignal(message []byte) error {
	// Minimum UTF16 signal message structure:
	// 1. message type (1 byte)
	// 2. name length flag (1 byte)
	// 3. name (n bytes, required if name length flag is bigger zero)
	// 4. payload (n bytes, at least 1 byte)
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
	// Minimum UTF16 signal message structure:
	// 1. message type (1 byte)
	// 2. name length flag (1 byte)
	// 3. name (n bytes, required if name length flag is bigger zero)
	// 4. header padding (1 byte, present if name length is odd)
	// 5. payload (n bytes, at least 2 bytes)
	if len(message) < MsgMinLenSignalUtf16 {
		return fmt.Errorf("Invalid signal message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unalligned UTF16 encoded signal message (probably missing header padding)",
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
	// Minimum binary/UTF8 request message structure:
	// 1. message type (1 byte)
	// 2. message id (8 bytes)
	// 3. name length flag (1 byte)
	// 4. name (n bytes, optional)
	// 5. payload (n bytes, at least 1 byte)
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

	// Verify total message size to prevent segmentation faults caused by inconsistent flags,
	// this could happen if the specified name length doesn't correspond to the actual name length
	if len(message) < MsgMinLenRequest+nameLen {
		return fmt.Errorf(
			"Invalid request message, too short for full name (%d) and the minimum payload (1)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.Name = string(message[10 : 10+nameLen])
		msg.Payload = Payload{
			Data: message[payloadOffset:],
		}
	} else {
		// No name present, just payload
		msg.Payload = Payload{
			Data: message[10:],
		}
	}

	return nil
}

func (msg *Message) parseRequestUtf16(message []byte) error {
	// Minimum UTF16 request message structure:
	// 1. message type (1 byte)
	// 2. message id (8 bytes)
	// 3. name length flag (1 byte)
	// 4. name (n bytes, required if name length flag is bigger zero)
	// 5. header padding (1 byte, present if name length is odd)
	// 6. payload (n bytes, at least 2 bytes)
	if len(message) < MsgMinLenRequestUtf16 {
		return fmt.Errorf("Invalid request message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unalligned UTF16 encoded request message (probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Read name length
	nameLen := int(byte(message[9:10][0]))

	// Determine minimum required message length
	minMsgSize := MsgMinLenRequestUtf16 + nameLen

	// Check whether a name padding byte is to be expected
	payloadOffset := 10 + nameLen
	if nameLen%2 != 0 {
		minMsgSize++
		payloadOffset++
	}

	// Verify total message size to prevent segmentation faults caused by inconsistent flags,
	// this could happen if the specified name length doesn't correspond to the actual name length
	if len(message) < minMsgSize {
		return fmt.Errorf(
			"Invalid request message, too short for full name (%d) and the minimum payload (2)",
			nameLen,
		)
	}

	if nameLen > 0 {
		// Take name into account
		msg.Name = string(message[10 : 10+nameLen])
		msg.Payload = Payload{
			Data: message[payloadOffset:],
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
	// Minimum binary/UTF8 reply message structure:
	// 1. message type (1 byte)
	// 2. message id (8 bytes)
	// 3. payload (n bytes, at least 1 byte)
	if len(message) < MsgMinLenReply {
		return fmt.Errorf("Invalid reply message, too short")
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

func (msg *Message) parseReplyUtf16(message []byte) error {
	// Minimum UTF16 reply message structure:
	// 1. message type (1 byte)
	// 2. message id (8 bytes)
	// 3. header padding (1 byte)
	// 4. payload (n bytes, at least 2 bytes)
	if len(message) < MsgMinLenReplyUtf16 {
		return fmt.Errorf("Invalid UTF16 reply message, too short")
	}

	if len(message)%2 != 0 {
		return fmt.Errorf(
			"Unalligned UTF16 encoded reply message (probably missing header padding)",
		)
	}

	// Read identifier
	var id [8]byte
	copy(id[:], message[1:9])
	msg.id = id

	// Read payload
	msg.Payload = Payload{
		// Take header padding byte into account
		Data: message[10:],
	}
	return nil
}

func (msg *Message) parseErrorReply(message []byte) error {
	if len(message) < MsgMinLenErrorReply {
		return fmt.Errorf("Invalid error reply message, too short")
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

// Parse tries to parse the message from a byte slice
func (msg *Message) Parse(message []byte) (err error) {
	if len(message) < 1 {
		return fmt.Errorf("Invalid message, too short")
	}
	var payloadEncoding PayloadEncoding
	msgType := message[0:1][0]

	switch msgType {

	// Request message format: [1 (type), 32 (id), | 1+ (payload)]
	case MsgErrorReply:
		err = msg.parseErrorReply(message)

	// Session creation notification format [1 (type), 32 (id), | 1+ (payload)]
	case MsgSessionCreated:
		err = msg.parseSessionCreated(message)

	// Session closure notification message format [1 (type), 32 (id)]
	case MsgSessionClosed:
		err = msg.parseSessionClosed(message)

	// Session destruction request message format [1 (type), 32 (id)]
	case MsgCloseSession:
		err = msg.parseCloseSession(message)

	// Request message format: [1 (type), 1+ (payload)]
	case MsgSignalBinary:
		payloadEncoding = EncodingBinary
		err = msg.parseSignal(message)
	case MsgSignalUtf8:
		payloadEncoding = EncodingUtf8
		err = msg.parseSignal(message)
	case MsgSignalUtf16:
		payloadEncoding = EncodingUtf16
		err = msg.parseSignalUtf16(message)

	// Request message format: [1 (type), 32 (id), | 1+ (payload)]
	case MsgRequestBinary:
		payloadEncoding = EncodingBinary
		err = msg.parseRequest(message)
	case MsgRequestUtf8:
		payloadEncoding = EncodingUtf8
		err = msg.parseRequest(message)
	case MsgRequestUtf16:
		payloadEncoding = EncodingUtf16
		err = msg.parseRequestUtf16(message)

	// Reply message format:
	case MsgReplyBinary:
		payloadEncoding = EncodingBinary
		err = msg.parseReply(message)
	case MsgReplyUtf8:
		payloadEncoding = EncodingUtf8
		err = msg.parseReply(message)
	case MsgReplyUtf16:
		payloadEncoding = EncodingUtf16
		err = msg.parseReplyUtf16(message)

	// Session restoration request message format: [1 (type), 32 (id), | 1+ (payload)]
	case MsgRestoreSession:
		err = msg.parseRestoreSession(message)

	// Ignore messages of invalid message type
	default:
		return fmt.Errorf("Invalid message type (%d)", msgType)
	}

	msg.msgType = msgType
	msg.Payload.Encoding = payloadEncoding
	return err
}

func (msg *Message) createFailCallback(client *Client, srv *Server) {
	msg.fail = func(errObj Error) {
		encoded, err := json.Marshal(errObj)
		if err != nil {
			encoded = []byte("CRITICAL: could not encode error report")
		}

		// Send request failure notification
		header := append([]byte{MsgErrorReply}, msg.id[:]...)
		if err = client.write(append(header, encoded...)); err != nil {
			srv.errorLog.Println("Writing failed:", err)
		}
	}
}

func (msg *Message) createReplyCallback(client *Client, srv *Server) {
	msg.fulfill = func(reply Payload) {
		replyType := MsgReplyBinary
		switch reply.Encoding {
		case EncodingUtf8:
			replyType = MsgReplyUtf8
		case EncodingUtf16:
			replyType = MsgReplyUtf16
		}

		header := append([]byte{replyType}, msg.id[:]...)
		// Send reply
		if err := client.write(append(header, reply.Data...)); err != nil {
			srv.errorLog.Println("Writing failed:", err)
		}
	}
}
