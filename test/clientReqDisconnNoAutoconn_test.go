package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientReqDisconnNoAutoconn tests disconnected error
// when trying to send a request while the server is unreachable
// and autoconn is disabled
func TestClientReqDisconnNoAutoconn(t *testing.T) {
	// Initialize client
	client := newCallbackPoweredClient(
		"127.0.0.1:65000",
		wwrclt.Options{
			Autoconnect:           wwr.Disabled,
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
	)

	// Try to send a request and expect a DisconnectedErr error
	_, err := client.connection.Request(
		"",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	)
	if _, isDisconnErr := err.(wwr.DisconnectedErr); !isDisconnErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
