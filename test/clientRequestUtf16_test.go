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
	testPayload := webwire.Payload{
		Encoding: webwire.EncodingUtf16,
		Data:     []byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	}
	verifyPayload := func(payload webwire.Payload) {
		if payload.Encoding != webwire.EncodingUtf16 {
			t.Errorf("Unexpected payload encoding: %s", payload.Encoding.String())
		}
		if len(testPayload.Data) != len(payload.Data) {
			t.Errorf("Corrupt payload: %s", payload.Data)
		}
		for i := 0; i < len(testPayload.Data); i++ {
			if testPayload.Data[i] != payload.Data[i] {
				t.Errorf("Corrupt payload, mismatching byte at position %d: %s", i, payload.Data)
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
				_ *webwire.Client,
				msg *webwire.Message,
			) (webwire.Payload, error) {

				verifyPayload(msg.Payload)

				return webwire.Payload{
					Encoding: webwire.EncodingUtf16,
					Data:     []byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
				}, nil
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
		nil, nil, nil, nil,
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send request and await reply
	reply, err := client.connection.Request("", webwire.Payload{
		Encoding: webwire.EncodingUtf16,
		Data:     []byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	})
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	verifyPayload(reply)
}
