package webwire_test

import (
	"testing"

	wwr "github.com/qbeon/webwire-go"
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestMsgWrapperGetters tests the getter methods of the message wrapper
// interface
func TestMsgWrapperGetters(t *testing.T) {
	// Create a new wrapped message
	wrappedMsg := wwr.NewMessageWrapper(&msg.Message{
		Type:       msg.MsgRequestBinary,
		Identifier: [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
		Name:       "sample-name",
		Payload: pld.Payload{
			Encoding: pld.Binary,
			Data:     []byte("sample-data"),
		},
	})

	require.Equal(t, [8]byte{1, 2, 3, 4, 5, 6, 7, 8}, wrappedMsg.Identifier())
	require.Equal(t, msg.MsgRequestBinary, wrappedMsg.MessageType())
	require.Equal(t, "sample-name", wrappedMsg.Name())

	pld := wrappedMsg.Payload()
	require.Equal(t, wwr.EncodingBinary, pld.Encoding())
	require.Equal(t, []byte("sample-data"), pld.Data())
}
