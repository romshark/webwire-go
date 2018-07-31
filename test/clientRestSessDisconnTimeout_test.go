package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRestSessDisconnTimeout tests autoconnect timeout
// when the server is unreachable
func TestClientRestSessDisconnTimeout(t *testing.T) {
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
	err := client.connection.RestoreSession([]byte("inexistentkey"))
	_, isTimeoutErr := err.(wwr.TimeoutErr)
	if !isTimeoutErr || !wwr.IsTimeoutErr(err) {
		t.Fatalf(
			"Expected request timeout error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
