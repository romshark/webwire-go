package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestError verifies returned request errors properly
// fail the request on the client
func TestClientRequestError(t *testing.T) {
	expectedRequestPayload := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("webwire_test_REQUEST_payload"),
	}
	expectedReplyError := webwire.Error{
		Code:    "SAMPLE_ERROR",
		Message: "Sample error message",
	}

	// Initialize webwire server given only the request
	_, addr := setupServer(
		t,
		webwire.ServerOptions{
			Hooks: webwire.Hooks{
				OnRequest: func(_ context.Context) (webwire.Payload, error) {
					// Fail the request by returning an error
					return webwire.Payload{}, expectedReplyError
				},
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send request and await reply
	reply, err := client.Request("", expectedRequestPayload)

	// Verify returned error
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	if err.Code != expectedReplyError.Code {
		t.Fatalf(
			"Unexpected error code: '%s' (%d)",
			err.Code,
			len(err.Code),
		)
	}

	if err.Message != expectedReplyError.Message {
		t.Fatalf(
			"Unexpected error message: '%s' (%d)",
			err.Message,
			len(err.Message),
		)
	}

	if len(reply.Data) > 0 {
		t.Fatalf(
			"Reply should have been empty, but was: '%s' (%d)",
			string(reply.Data),
			len(reply.Data),
		)
	}
}
