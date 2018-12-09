package memchan_test

import (
	"sync"
	"testing"
	"time"

	"github.com/qbeon/webwire-go/payload"

	"github.com/qbeon/webwire-go/message"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/require"
)

func testNewServer(serverOptions ...wwr.ServerOptions) *memchan.Transport {
	server := &memchan.Transport{}

	if len(serverOptions) < 1 {
		serverOptions = []wwr.ServerOptions{
			wwr.ServerOptions{
				MessageBufferSize: 1024,
			},
		}
	}

	server.Initialize(
		serverOptions[0],
		func() bool { return false },
		func(_ wwr.ConnectionOptions, _ wwr.Socket) {},
	)
	return server
}

func createSockets(
	t *testing.T,
	server *memchan.Transport,
) (*memchan.Socket, *memchan.Socket) {
	t.Helper()
	srvSock, cltSock := memchan.NewEntangledSockets(server)
	require.NotNil(t, srvSock)
	require.NotNil(t, cltSock)
	return srvSock, cltSock
}

func ensureSocketClosed(t *testing.T, sock *memchan.Socket) {
	t.Helper()

	require.False(t, sock.IsConnected())

	// Ensure no writer
	writer, err := sock.GetWriter()
	require.Error(t, err)
	require.Nil(t, writer)

	// Ensure read fails
	readErr := sock.Read(message.NewMessage(32), time.Time{})
	require.NotNil(t, readErr)
}

func dial(t *testing.T, srvSock, cltSock *memchan.Socket) {
	t.Helper()

	// Ensure both sockets are yet disconnected
	ensureSocketClosed(t, srvSock)
	ensureSocketClosed(t, cltSock)

	require.NoError(t, cltSock.Dial(time.Time{}))

	// Ensure both sockets are connected
	require.True(t, cltSock.IsConnected())
	require.True(t, srvSock.IsConnected())
}

func generateReqMsgBytes() []byte {
	expectedMsg := make([]byte, 18)
	expectedMsg[0] = message.MsgRequestBinary
	copy(expectedMsg[1:], []byte("00000000"))
	expectedMsg[9] = 0
	copy(expectedMsg[10:], []byte("12345678"))
	return expectedMsg
}

func writeRequestMessage(
	t *testing.T,
	sender *memchan.Socket,
	receiver *memchan.Socket,
) {
	t.Helper()

	// Get writer
	writer, err := sender.GetWriter()
	require.NoError(t, err)
	require.NotNil(t, writer)

	expectedMsg := generateReqMsgBytes()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start reader
	go func() {
		defer wg.Done()

		// Read
		msg := message.NewMessage(32)
		require.Nil(t, receiver.Read(msg, time.Time{}))

		// Compare
		require.Equal(t, expectedMsg, msg.MsgBuffer.Data())
	}()

	// Write
	require.NoError(t, message.WriteMsgRequest(
		writer,
		[]byte("00000000"),
		nil,
		payload.Binary,
		[]byte("12345678"),
		true,
	))

	wg.Wait()
}

// TestClientDial tests dialing
func TestClientDial(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)

	dial(t, srvSock, cltSock)
}

// TestClientNoDial tests socket state before connection establishment
func TestClientNoDial(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)

	ensureSocketClosed(t, cltSock)
	ensureSocketClosed(t, srvSock)
}

// TestClientDialServerFail tests dialing on a server-side socket expecting it
// to fail
func TestClientDialServerFail(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)

	require.Error(t, srvSock.Dial(time.Time{}))

	ensureSocketClosed(t, cltSock)
	ensureSocketClosed(t, srvSock)
}

// TestClientClose tests the closing a socket
func TestClientClose(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	require.NoError(t, cltSock.Close())

	ensureSocketClosed(t, cltSock)
	ensureSocketClosed(t, srvSock)
}

// TestClientRepeatedClose tests closing a socket multiple times in a row
func TestClientRepeatedClose(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	require.NoError(t, cltSock.Close())

	ensureSocketClosed(t, cltSock)
	ensureSocketClosed(t, srvSock)

	require.NoError(t, cltSock.Close())
}

// TestCloseBeforeDial tests the closing a socket before it's even connected
func TestCloseBeforeDial(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)

	require.NoError(t, cltSock.Close())

	ensureSocketClosed(t, cltSock)
	ensureSocketClosed(t, srvSock)
}

// TestServerSocketClose tests closing the server socket
func TestServerSocketClose(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	require.NoError(t, srvSock.Close())

	ensureSocketClosed(t, cltSock)
	ensureSocketClosed(t, srvSock)
}

// TestReconnect tests a dial-close-dial scenario
func TestReconnect(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	require.NoError(t, cltSock.Close())

	ensureSocketClosed(t, cltSock)
	ensureSocketClosed(t, srvSock)

	dial(t, srvSock, cltSock)
}

// TestRepeatedDial tests a repeated dialing
func TestRepeatedDial(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	t.Helper()
	require.Error(t, cltSock.Dial(time.Time{}))
	require.True(t, cltSock.IsConnected())
	require.True(t, srvSock.IsConnected())
}

// TestClientSend tests sending client -> server
func TestClientSend(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	writeRequestMessage(t, cltSock, srvSock)
}

// TestSrvSend tests sending client -> server
func TestSrvSend(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	writeRequestMessage(t, srvSock, cltSock)
}

// TestReadDeadline tests reading with a deadline
func TestReadDeadline(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	deadline := time.Now().Add(50 * time.Millisecond)

	// Read for 50 milliseconds
	err := srvSock.Read(message.NewMessage(32), deadline)

	require.NotNil(t, err)
	require.False(t, err.IsCloseErr())
}

// TestCloseWhileRead tests closing a socket while it's reading
func TestCloseWhileRead(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start reader
	go func() {
		defer wg.Done()

		// Read
		err := srvSock.Read(message.NewMessage(32), time.Time{})

		require.NotNil(t, err)
		require.True(t, err.IsCloseErr())
	}()

	// Wait until the reader is reading
	time.Sleep(50 * time.Millisecond)

	// Close the reading socket
	require.NoError(t, srvSock.Close())

	ensureSocketClosed(t, srvSock)
	ensureSocketClosed(t, cltSock)

	wg.Wait()
}

// TestCloseWhileOpponentRead tests closing the socket while the opposite socket
// is reading
func TestCloseWhileOpponentRead(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start reader
	go func() {
		defer wg.Done()

		// Read on client
		err := cltSock.Read(message.NewMessage(32), time.Time{})

		require.NotNil(t, err)
		require.True(t, err.IsCloseErr())
	}()

	// Wait until the reader is reading
	time.Sleep(50 * time.Millisecond)

	// Close the opposite socket
	require.NoError(t, srvSock.Close())

	ensureSocketClosed(t, srvSock)
	ensureSocketClosed(t, cltSock)

	wg.Wait()
}

// TestCloseWhileWrite tests closing a socket while it's writing
func TestCloseWhileWrite(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start writer
	go func() {
		defer wg.Done()

		// Write
		writer, err := srvSock.GetWriter()
		require.NoError(t, err)
		require.NotNil(t, writer)

		require.Error(t, message.WriteMsgRequest(
			writer,
			[]byte("00000000"),
			nil,
			payload.Binary,
			[]byte("12345678"),
			true,
		))
	}()

	// Wait until the writer has finished writing
	time.Sleep(100 * time.Millisecond)

	// Close the writing socket
	require.NoError(t, srvSock.Close())

	ensureSocketClosed(t, srvSock)
	ensureSocketClosed(t, cltSock)

	// Ensure the reader doesn't receive the message
	msg := message.NewMessage(32)
	err := cltSock.Read(msg, time.Time{})
	require.NotNil(t, err)

	wg.Wait()
}

// TestCloseWhileOpponentWrite tests closing a socket while the opposite socket
// is writing
func TestCloseWhileOpponentWrite(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start writer
	go func() {
		defer wg.Done()

		// Write
		writer, err := srvSock.GetWriter()
		require.NoError(t, err)
		require.NotNil(t, writer)

		require.Error(t, message.WriteMsgRequest(
			writer,
			[]byte("00000000"),
			nil,
			payload.Binary,
			[]byte("12345678"),
			true,
		))
	}()

	// Wait until the writer has finished writing
	time.Sleep(100 * time.Millisecond)

	// Close the receiver socket
	require.NoError(t, cltSock.Close())

	ensureSocketClosed(t, srvSock)
	ensureSocketClosed(t, cltSock)

	// Ensure the reader doesn't receive the message
	msg := message.NewMessage(32)
	err := cltSock.Read(msg, time.Time{})
	require.NotNil(t, err)

	wg.Wait()
}

// TestRemoveAddr tests the RemoteAddr getter method
func TestRemoveAddr(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	testRemoteAddr := func(sock *memchan.Socket) {
		addr := sock.RemoteAddr()
		require.NotNil(t, addr)
		require.Equal(t, "memchan", addr.Network())
	}

	testRemoteAddr(srvSock)
	testRemoteAddr(cltSock)
}

// TestClientConcurrentReadWrite tests sending client -> server
func TestClientConcurrentReadWrite(t *testing.T) {
	server := testNewServer()
	srvSock, cltSock := createSockets(t, server)
	dial(t, srvSock, cltSock)

	concurrentReaders := 128
	concurrentWriters := concurrentReaders
	expectedMsg := generateReqMsgBytes()

	readerFinished := sync.WaitGroup{}
	readerFinished.Add(concurrentReaders)
	writersFinished := sync.WaitGroup{}
	writersFinished.Add(concurrentWriters)

	// Start multiple concurrent reader goroutines
	for r := 0; r < concurrentReaders; r++ {
		go func() {
			defer readerFinished.Done()

			// Read
			msg := message.NewMessage(32)
			require.Nil(t, srvSock.Read(msg, time.Time{}))

			// Compare
			require.Equal(t, expectedMsg, msg.MsgBuffer.Data())
		}()
	}

	// Start multiple concurrent writer goroutines
	for i := 0; i < concurrentWriters; i++ {
		go func() {
			defer writersFinished.Done()

			// Get writer
			writer, err := cltSock.GetWriter()
			require.NoError(t, err)
			require.NotNil(t, writer)

			// Write message
			require.NoError(t, message.WriteMsgRequest(
				writer,
				[]byte("00000000"),
				nil,
				payload.Binary,
				[]byte("12345678"),
				true,
			))
		}()
	}

	readerFinished.Wait()
	writersFinished.Wait()
}
