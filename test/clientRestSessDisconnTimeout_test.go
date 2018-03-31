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
	client := wwrclt.NewClient(
		"127.0.0.1:65000",
		wwrclt.Options{
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
	)

	// Send request and await reply
	err := client.RestoreSession([]byte("inexistentkey"))
	if _, isReqTimeoutErr := err.(wwr.ReqTimeoutErr); !isReqTimeoutErr {
		t.Fatalf(
			"Expected request timeout error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
