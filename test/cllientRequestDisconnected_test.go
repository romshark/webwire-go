package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRequestDisconnected tests sending requests on disconnected clients
func TestClientRequestDisconnected(t *testing.T) {
	// Initialize webwire server given only the request
	_, addr := setupServer(
		t,
		wwr.ServerOptions{},
	)

	// Initialize client
	client := wwrclt.NewClient(
		addr,
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	// Send request and await reply
	_, err := client.Request("", wwr.Payload{Data: []byte("testdata")})
	if _, isDisconnErr := err.(wwrclt.DisconnectedErr); !isDisconnErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
