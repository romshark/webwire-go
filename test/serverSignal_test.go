package test

import (
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerSignal tests server-side signals
func TestServerSignal(t *testing.T) {
	finished := sync.WaitGroup{}
	finished.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(
				_ wwr.ConnectionOptions,
				conn wwr.Connection,
			) {
				// Send unnamed signal with binary payload
				assert.NoError(t, conn.Signal(nil, wwr.Payload{
					Encoding: wwr.EncodingBinary,
					Data:     []byte("binary"),
				}))

				// Send unnamed signal with UTF8 encoded payload
				assert.NoError(t, conn.Signal(nil, wwr.Payload{
					Encoding: wwr.EncodingUtf8,
					Data:     []byte("üникод"),
				}))

				// Send unnamed signal with UTF16 encoded payload
				assert.NoError(t, conn.Signal(nil, wwr.Payload{
					Encoding: wwr.EncodingUtf16,
					Data:     []byte{100, 200, 110, 210},
				}))

				// Send named signal with binary payload
				assert.NoError(t, conn.Signal([]byte("bin_sig"), wwr.Payload{
					Encoding: wwr.EncodingBinary,
					Data:     []byte("binary"),
				}))

				// Send named signal with UTF8 encoded payload
				assert.NoError(t, conn.Signal([]byte("utf8_sig"), wwr.Payload{
					Encoding: wwr.EncodingUtf8,
					Data:     []byte("üникод"),
				}))

				// Send named signal with UTF16 encoded payload
				assert.NoError(t, conn.Signal([]byte("utf16_sig"), wwr.Payload{
					Encoding: wwr.EncodingUtf16,
					Data:     []byte{100, 200, 110, 210},
				}))

				finished.Done()
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	/* unnamed signals */

	// Expect unnamed binary signal message
	msgBin := message.NewMessage(64)
	sock.Read(msgBin, time.Time{})

	require.Equal(t, message.MsgSignalBinary, msgBin.MsgType)
	require.Equal(t, wwr.EncodingBinary, msgBin.PayloadEncoding())
	require.Equal(t, []byte("binary"), msgBin.Payload())
	require.Equal(t, []byte(nil), msgBin.MsgName)

	// Expect unnamed UTF8 signal message
	msgUtf8 := message.NewMessage(64)
	sock.Read(msgUtf8, time.Time{})

	require.Equal(t, message.MsgSignalUtf8, msgUtf8.MsgType)
	require.Equal(t, wwr.EncodingUtf8, msgUtf8.PayloadEncoding())
	require.Equal(t, []byte("üникод"), msgUtf8.Payload())
	require.Equal(t, []byte(nil), msgUtf8.MsgName)

	// Expect unnamed UTF16 signal message
	msgUtf16 := message.NewMessage(64)
	sock.Read(msgUtf16, time.Time{})

	require.Equal(t, message.MsgSignalUtf16, msgUtf16.MsgType)
	require.Equal(t, wwr.EncodingUtf16, msgUtf16.PayloadEncoding())
	require.Equal(t, []byte{100, 200, 110, 210}, msgUtf16.Payload())
	require.Equal(t, []byte(nil), msgUtf16.MsgName)

	/* named signals */

	// Expect unnamed binary signal message
	namedMsgBin := message.NewMessage(64)
	sock.Read(namedMsgBin, time.Time{})

	require.Equal(t, message.MsgSignalBinary, namedMsgBin.MsgType)
	require.Equal(t, wwr.EncodingBinary, namedMsgBin.PayloadEncoding())
	require.Equal(t, []byte("binary"), namedMsgBin.Payload())
	require.Equal(t, []byte("bin_sig"), namedMsgBin.MsgName)

	// Expect unnamed UTF8 signal message
	namedMsgUtf8 := message.NewMessage(64)
	sock.Read(namedMsgUtf8, time.Time{})

	require.Equal(t, message.MsgSignalUtf8, namedMsgUtf8.MsgType)
	require.Equal(t, wwr.EncodingUtf8, namedMsgUtf8.PayloadEncoding())
	require.Equal(t, []byte("üникод"), namedMsgUtf8.Payload())
	require.Equal(t, []byte("utf8_sig"), namedMsgUtf8.MsgName)

	// Expect unnamed UTF16 signal message
	namedMsgUtf16 := message.NewMessage(64)
	sock.Read(namedMsgUtf16, time.Time{})

	require.Equal(t, message.MsgSignalUtf16, namedMsgUtf16.MsgType)
	require.Equal(t, wwr.EncodingUtf16, namedMsgUtf16.PayloadEncoding())
	require.Equal(t, []byte{100, 200, 110, 210}, namedMsgUtf16.Payload())
	require.Equal(t, []byte("utf16_sig"), namedMsgUtf16.MsgName)

	finished.Wait()
}
