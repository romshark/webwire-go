package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestError tests server-side request errors properly
// failing the client-side requests
func TestClientRequestError(t *testing.T) {
	expectedReplyError := webwire.ReqErr{
		Code:    "SAMPLE_ERROR",
		Message: "Sample error message",
	}

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ webwire.Connection,
				_ webwire.Message,
			) (webwire.Payload, error) {
				// Fail the request by returning an error
				return nil, expectedReplyError
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		"",
		webwire.NewPayload(
			webwire.EncodingUtf8,
			[]byte("webwire_test_REQUEST_payload"),
		),
	)

	// Verify returned error
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	if err, isReqErr := err.(webwire.ReqErr); isReqErr {
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
	} else {
		t.Fatalf("Unexpected request failure: %v", err)
	}

	if reply != nil {
		t.Fatalf(
			"Reply should have been empty, but was: '%s' (%d)",
			string(reply.Data()),
			len(reply.Data()),
		)
	}
}
