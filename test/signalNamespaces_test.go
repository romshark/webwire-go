package test

import (
	"context"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestSignalNamespaces correct handling of namespaced signals
func TestSignalNamespaces(t *testing.T) {
	unnamedSignalArrived := NewPending(1, 1*time.Second, true)
	shortestNameSignalArrived := NewPending(1, 1*time.Second, true)
	longestNameSignalArrived := NewPending(1, 1*time.Second, true)
	currentStep := 1

	shortestPossibleName := "s"
	buf := make([]rune, 255)
	for i := 0; i < 255; i++ {
		buf[i] = 'x'
	}
	longestPossibleName := string(buf)

	// Initialize server
	_, addr := setupServer(
		t,
		webwire.Options{
			Hooks: webwire.Hooks{
				OnSignal: func(ctx context.Context) {
					msg := ctx.Value(webwire.Msg).(webwire.Message)

					if currentStep == 1 && msg.Name != "" {
						t.Errorf("Expected unnamed signal, got: '%s'", msg.Name)
					}
					if currentStep == 2 && msg.Name != shortestPossibleName {
						t.Errorf("Expected shortest possible signal name, got: '%s'", msg.Name)
					}
					if currentStep == 3 && msg.Name != longestPossibleName {
						t.Errorf("Expected longest possible signal name, got: '%s'", msg.Name)
					}

					switch currentStep {
					case 1:
						unnamedSignalArrived.Done()
					case 2:
						shortestNameSignalArrived.Done()
					case 3:
						longestNameSignalArrived.Done()
					}
				},
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	/*****************************************************************\
		Step 1 - Unnamed signal name
	\*****************************************************************/
	// Send unnamed signal
	err := client.Signal("", webwire.Payload{Data: []byte("dummy")})
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
	err = client.Signal(shortestPossibleName, webwire.Payload{Data: []byte("dummy")})
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
	err = client.Signal(longestPossibleName, webwire.Payload{Data: []byte("dummy")})
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	if err := longestNameSignalArrived.Wait(); err != nil {
		t.Fatal("Signal with longest name didn't arrive")
	}
}
