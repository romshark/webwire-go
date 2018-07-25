package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestSignalNamespaces tests correct handling of namespaced signals
func TestSignalNamespaces(t *testing.T) {
	unnamedSignalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	shortestNameSignalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	longestNameSignalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
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
			onSignal: func(
				_ context.Context,
				_ *webwire.Client,
				msg webwire.Message,
			) {
				msgName := msg.Name()
				if currentStep == 1 && msgName != "" {
					t.Errorf(
						"Expected unnamed signal, got: '%s'",
						msgName,
					)
				}
				if currentStep == 2 &&
					msgName != shortestPossibleName {
					t.Errorf(
						"Expected shortest possible signal name, got: '%s'",
						msgName,
					)
				}
				if currentStep == 3 &&
					msgName != longestPossibleName {
					t.Errorf(
						"Expected longest possible signal name, got: '%s'",
						msgName,
					)
				}

				switch currentStep {
				case 1:
					unnamedSignalArrived.Progress(1)
				case 2:
					shortestNameSignalArrived.Progress(1)
				case 3:
					longestNameSignalArrived.Progress(1)
				}
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
		Step 1 - Unnamed signal name
	\*****************************************************************/
	// Send unnamed signal
	err := client.connection.Signal("", webwire.NewPayload(
		webwire.EncodingBinary,
		[]byte("dummy"),
	))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	if err := unnamedSignalArrived.Wait(); err != nil {
		t.Fatal("Unnamed signal didn't arrive")
	}

	/*****************************************************************\
		Step 2 - Shortest possible request name
	\*****************************************************************/
	currentStep = 2
	// Send request with the shortest possible name
	err = client.connection.Signal(
		shortestPossibleName,
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	if err := shortestNameSignalArrived.Wait(); err != nil {
		t.Fatal("Signal with shortest name didn't arrive")
	}

	/*****************************************************************\
		Step 3 - Longest possible request name
	\*****************************************************************/
	currentStep = 3
	// Send request with the longest possible name
	err = client.connection.Signal(
		longestPossibleName,
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	if err := longestNameSignalArrived.Wait(); err != nil {
		t.Fatal("Signal with longest name didn't arrive")
	}
}
