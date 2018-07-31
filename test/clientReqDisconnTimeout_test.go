package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientReqDisconnTimeout tests request timeout
// when the server is unreachable and autoconnect is enabled
func TestClientReqDisconnTimeout(t *testing.T) {
	// Initialize client
	client := newCallbackPoweredClient(
		"127.0.0.1:65000",
		wwrclt.Options{
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
	)

	// Send request and await reply
	_, err := client.connection.Request(
		context.Background(),
		"",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	)
	_, isTimeoutErr := err.(wwr.TimeoutErr)
	if !isTimeoutErr || !wwr.IsTimeoutErr(err) {
		t.Fatalf(
			"Expected request timeout error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
