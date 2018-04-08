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
				_ *webwire.Client,
				msg *webwire.Message,
			) (webwire.Payload, error) {
				if currentStep == 1 && msg.Name != "" {
					t.Errorf(
						"Expected unnamed request, got: '%s'",
						msg.Name,
					)
				}
				if currentStep == 2 && msg.Name != shortestPossibleName {
					t.Errorf("Expected shortest possible "+
						"request name, got: '%s'",
						msg.Name,
					)
				}
				if currentStep == 3 && msg.Name != longestPossibleName {
					t.Errorf("Expected longest possible "+
						"request name, got: '%s'",
						msg.Name,
					)
				}

				return webwire.Payload{}, nil
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
	_, err := client.connection.Request("", webwire.Payload{
		Data: []byte("dummy"),
	})
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
		webwire.Payload{Data: []byte("dummy")},
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
		webwire.Payload{Data: []byte("dummy")},
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
}
