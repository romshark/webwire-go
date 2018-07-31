package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestUtf16 tests requests with UTF16 encoded payloads
func TestClientRequestUtf16(t *testing.T) {
	testPayload := webwire.NewPayload(
		webwire.EncodingUtf16,
		[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	)
	verifyPayload := func(payload webwire.Payload) {
		payloadEncoding := payload.Encoding()
		if payloadEncoding != webwire.EncodingUtf16 {
			t.Errorf("Unexpected payload encoding: %s", payloadEncoding.String())
		}
		testPayloadData := testPayload.Data()
		payloadData := payload.Data()
		if len(testPayloadData) != len(payloadData) {
			t.Errorf("Corrupt payload: %s", payloadData)
		}
		for i := 0; i < len(testPayloadData); i++ {
			if testPayloadData[i] != payloadData[i] {
				t.Errorf("Corrupt payload, mismatching byte at position %d: %s", i, payloadData)
				return
			}
		}
	}

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ webwire.Connection,
				msg webwire.Message,
			) (webwire.Payload, error) {

				verifyPayload(msg.Payload())

				return webwire.NewPayload(
					webwire.EncodingUtf16,
					[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
				), nil
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
			webwire.EncodingUtf16,
			[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
		),
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	verifyPayload(reply)
}
