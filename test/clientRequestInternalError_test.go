package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestInternalError tests returning of non-ReqErr errors
// from the request handler
func TestClientRequestInternalError(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		webwire.ServerOptions{
			Hooks: webwire.Hooks{
				OnRequest: func(_ context.Context) (webwire.Payload, error) {
					// Fail the request by returning a non-ReqErr error
					return webwire.Payload{}, fmt.Errorf(
						"don't worry, this internal error is expected",
					)
				},
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send request and await reply
	reply, reqErr := client.Request("", webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("dummydata"),
	})

	// Verify returned error
	if reqErr == nil {
		t.Fatal("Expected an error, got nil")
	}

	if _, isInternalErr := reqErr.(webwire.ReqInternalErr); !isInternalErr {
		t.Fatalf("Expected an internal server error, got: %v", reqErr)
	}

	if len(reply.Data) > 0 {
		t.Fatalf(
			"Reply should have been empty, but was: '%s' (%d)",
			string(reply.Data),
			len(reply.Data),
		)
	}
}
