package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/payload"
	"github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ServerSetup represents a webwire server setup
type ServerSetup struct {
	Transport wwr.Transport
	Server    wwr.Server
}

// ServerSetupTest represents a webwire server testing setup
type ServerSetupTest struct {
	t *testing.T
	ServerSetup
}

// SetupServer helps setting up and launching the server together with the
// underlying transport
func SetupServer(
	impl *ServerImpl,
	opts wwr.ServerOptions,
	trans wwr.Transport,
) (ServerSetup, error) {
	// Use default session manager if no specific one is defined
	if opts.SessionManager == nil {
		opts.SessionManager = newInMemSessManager()
	}

	// Use the transport layer implementation specified by the CLI arguments
	if trans == nil {
		// Use default configuration
		trans = &memchan.Transport{}
	}

	// Initialize webwire server
	server, err := wwr.NewServer(impl, opts, trans)
	if err != nil {
		return ServerSetup{}, fmt.Errorf(
			"failed setting up server instance: %s",
			err,
		)
	}

	// Run server in a separate goroutine
	go func() {
		if err := server.Run(); err != nil {
			panic(fmt.Errorf("server failed: %s", err))
		}
	}()

	// Return reference to the server and the address its bound to
	return ServerSetup{
		Server:    server,
		Transport: trans,
	}, nil
}

// SetupTestServer creates a new server setup failing the test immediately if
// the anything went wrong
func SetupTestServer(
	t *testing.T,
	impl *ServerImpl,
	opts wwr.ServerOptions,
	trans wwr.Transport,
) ServerSetupTest {
	setup, err := SetupServer(impl, opts, trans)
	require.NoError(t, err)
	return ServerSetupTest{t, setup}
}

// NewDisconnectedClientSocket creates a new raw disconnected client socket
func (setup *ServerSetup) NewDisconnectedClientSocket() (
	wwr.ClientSocket,
	error,
) {
	var sock wwr.ClientSocket
	switch srvTrans := setup.Transport.(type) {
	case *memchan.Transport:
		_, sock = memchan.NewEntangledSockets(srvTrans)
		return sock, nil
	}
	return nil, fmt.Errorf(
		"unexpected server transport implementation: %s",
		reflect.TypeOf(setup.Transport),
	)
}

// NewClientSocket creates a new raw client socket connected to the server
func (setup *ServerSetup) NewClientSocket() (
	wwr.Socket,
	message.ServerConfiguration,
	error,
) {
	sock, err := setup.NewDisconnectedClientSocket()
	if err != nil {
		return nil, message.ServerConfiguration{}, err
	}

	// Establish a connection
	if err := sock.Dial(time.Time{}); err != nil {
		return nil, message.ServerConfiguration{}, fmt.Errorf(
			"memchan dial failed: %s",
			err,
		)
	}

	// Read the server configuration push message
	msg := message.NewMessage(32)
	if err := sock.Read(msg, time.Time{}); err != nil {
		return nil, message.ServerConfiguration{}, fmt.Errorf(
			"couldn't read server configuration push message: %s",
			err,
		)
	}

	return sock, msg.ServerConfiguration, nil
}

// NewClientSocket creates a new raw client socket connected to the server
func (setup *ServerSetupTest) NewClientSocket() (
	wwr.Socket,
	message.ServerConfiguration,
) {
	sock, srvConf, err := setup.ServerSetup.NewClientSocket()
	require.NoError(setup.t, err)
	return sock, srvConf
}

// CompareSessions compares a webwire session
func CompareSessions(t *testing.T, expected, actual *wwr.Session) {
	if actual == nil && expected == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Creation.Unix(), actual.Creation.Unix())
}

// signal sends a signal message
func signal(
	t *testing.T,
	socket wwr.Socket,
	name []byte,
	payload payload.Payload,
) {
	writer, err := socket.GetWriter()
	require.NoError(t, err)
	require.NotNil(t, writer)

	require.NoError(t, message.WriteMsgSignal(
		writer,
		name,
		payload.Encoding,
		payload.Data,
		true,
	))
}

// request performs a synchronous request. Blocks until a reply is received
func request(
	t *testing.T,
	socket wwr.Socket,
	replyBufferSize uint32,
	name []byte,
	payload payload.Payload,
) *message.Message {
	writer, err := socket.GetWriter()
	require.NoError(t, err)
	require.NotNil(t, writer)

	reply := message.NewMessage(replyBufferSize)

	require.NoError(t, message.WriteMsgRequest(
		writer,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		name,
		payload.Encoding,
		payload.Data,
		true,
	))

	require.Nil(t, socket.Read(reply, time.Time{}))

	return reply
}

// requestSuccess performs a synchronous request and expects it to succeed.
// Blocks until a reply is received
func requestSuccess(
	t *testing.T,
	socket wwr.Socket,
	replyBufferSize uint32,
	name []byte,
	pld payload.Payload,
) *message.Message {
	reply := request(t, socket, replyBufferSize, name, pld)

	// Verify reply message type
	switch pld.Encoding {
	case payload.Binary:
		require.Equal(t, message.MsgReplyBinary, reply.MsgType)
		require.Equal(t, payload.Binary, reply.MsgPayload.Encoding)
	case payload.Utf8:
		require.Equal(t, message.MsgReplyUtf8, reply.MsgType)
		require.Equal(t, payload.Utf8, reply.MsgPayload.Encoding)
	case payload.Utf16:
		require.Equal(t, message.MsgReplyUtf16, reply.MsgType)
		require.Equal(t, payload.Utf16, reply.MsgPayload.Encoding)
	default:
		panic("unexpected payload encoding type")
	}

	return reply
}

// requestRestoreSession performs a synchronous session restoration request.
// Blocks until a reply is received
func requestRestoreSession(
	t *testing.T,
	socket wwr.Socket,
	sessionKey []byte,
) *message.Message {
	writer, err := socket.GetWriter()
	require.NoError(t, err)
	require.NotNil(t, writer)

	reply := message.NewMessage(1024)

	require.NoError(t, message.WriteMsgNamelessRequest(
		writer,
		message.MsgRequestRestoreSession,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		sessionKey,
	))

	require.Nil(t, socket.Read(reply, time.Time{}))

	return reply
}

// requestRestoreSessionSuccess performs a synchronous session restoration
// request and expects it to succeed. Blocks until a reply is received
func requestRestoreSessionSuccess(
	t *testing.T,
	socket wwr.Socket,
	sessionKey []byte,
) *message.Message {
	reply := requestRestoreSession(t, socket, sessionKey)
	require.Equal(t, message.MsgReplyUtf8, reply.MsgType)
	return reply
}

// requestCloseSession performs a synchronous session closure request. Blocks
// until a reply is received
func requestCloseSession(t *testing.T, socket wwr.Socket) *message.Message {
	writer, err := socket.GetWriter()
	require.NoError(t, err)
	require.NotNil(t, writer)

	reply := message.NewMessage(32)

	require.NoError(t, message.WriteMsgNamelessRequest(
		writer,
		message.MsgDoCloseSession,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		nil,
	))

	require.Nil(t, socket.Read(reply, time.Time{}))

	return reply
}

// requestCloseSessionSuccess performs a synchronous session closure request and
// expects it to succeed. Blocks until a reply is received
func requestCloseSessionSuccess(
	t *testing.T,
	socket wwr.Socket,
) *message.Message {
	reply := requestCloseSession(t, socket)
	require.Equal(t, message.MsgReplyBinary, reply.MsgType)
	return reply
}

// readSessionCreated reads a session creation notification message
func readSessionCreated(t *testing.T, sock wwr.Socket) *message.Message {
	// Expect session creation notification message
	msg := message.NewMessage(1024)
	require.Nil(t, sock.Read(msg, time.Time{}))
	require.Equal(t, message.MsgNotifySessionCreated, msg.MsgType)
	return msg
}

// readSessionClosed reads a session closure notification message
func readSessionClosed(t *testing.T, sock wwr.Socket) *message.Message {
	// Expect session creation notification message
	msg := message.NewMessage(1024)
	require.Nil(t, sock.Read(msg, time.Time{}))
	require.Equal(t, message.MsgNotifySessionClosed, msg.MsgType)
	return msg
}
