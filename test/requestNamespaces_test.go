package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestRequestNamespaces tests correct handling of namespaced requests
func TestRequestNamespaces(t *testing.T) {
	currentStep := 1

	shortestPossibleName := "s"
	buf := make([]rune, 255)
	for i := 0; i < 255; i++ {
		buf[i] = 'x'
	}
	longestPossibleName := string(buf)

	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ webwire.Connection,
				msg webwire.Message,
			) (webwire.Payload, error) {
				msgName := msg.Name()
				if currentStep == 1 && msgName != "" {
					t.Errorf(
						"Expected unnamed request, got: '%s'",
						msgName,
					)
				}
				if currentStep == 2 && msgName != shortestPossibleName {
					t.Errorf("Expected shortest possible "+
						"request name, got: '%s'",
						msgName,
					)
				}
				if currentStep == 3 && msgName != longestPossibleName {
					t.Errorf("Expected longest possible "+
						"request name, got: '%s'",
						msgName,
					)
				}

				return nil, nil
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

	/*****************************************************************\
		Step 1 - Unnamed request name
	\*****************************************************************/
	// Send unnamed request
	_, err := client.connection.Request("", webwire.NewPayload(
		webwire.EncodingBinary,
		[]byte("dummy"),
	))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	/*****************************************************************\
		Step 2 - Shortest possible request name
	\*****************************************************************/
	currentStep = 2
	// Send request with the shortest possible name
	_, err = client.connection.Request(
		shortestPossibleName,
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	/*****************************************************************\
		Step 3 - Longest possible request name
	\*****************************************************************/
	currentStep = 3
	// Send request with the longest possible name
	_, err = client.connection.Request(
		longestPossibleName,
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
}
