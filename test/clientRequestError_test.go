package test

import (
	"context"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
)

// TestClientRequestError verifies returned request errors properly
// fail the request on the client
func TestClientRequestError(t *testing.T) {
	expectedRequestPayload := []byte("webwire_test_REQUEST_payload")
	expectedReplyError := webwire.Error{
		Code:    "SAMPLE_ERROR",
		Message: "Sample error message",
	}

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		nil, nil, nil,
		func(_ context.Context) ([]byte, *webwire.Error) {
			// Fail the request by returning an error
			err := expectedReplyError
			return nil, &err
		},
		nil, nil, nil, nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil, nil, nil,
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	// Send request and await reply
	reply, err := client.Request(expectedRequestPayload)

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

	if len(reply) > 0 {
		t.Fatalf(
			"Reply should have been empty, but was: '%s' (%d)",
			string(reply),
			len(reply),
		)
	}
}
