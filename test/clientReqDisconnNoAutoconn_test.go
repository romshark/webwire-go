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
			Autoconnect:           wwrclt.OptDisabled,
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		nil, nil, nil, nil,
	)

	// Send request and await reply
	_, err := client.connection.Request(
		"",
		wwr.Payload{Data: []byte("testdata")},
	)
	if _, isDisconnErr := err.(wwr.DisconnectedErr); !isDisconnErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
