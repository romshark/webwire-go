package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRequestCancel tests canceling of fired requests
func TestClientRequestCancel(t *testing.T) {
	requestFinished := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				time.Sleep(2 * time.Second)
				return nil, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 5 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	cancelableCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Send request and await reply
	go func() {
		reply, err := client.connection.Request(cancelableCtx, "test", nil)
		if err == nil {
			t.Error("Expected a canceled-error, got nil")
		}
		if reply != nil {
			t.Errorf("Expected nil reply, got: %v", reply)
		}
		_, isCanceledErr := err.(wwr.CanceledErr)
		if !isCanceledErr || !wwr.IsCanceledErr(err) {
			t.Errorf(
				"Expected a canceled-error, got: (%s) %s",
				err,
				reflect.TypeOf(err),
			)
		}
		requestFinished.Progress(1)
	}()

	// Cancel the context some time after sending the request
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Wait for the requestor goroutine to finish
	if err := requestFinished.Wait(); err != nil {
		t.Fatal("Test timed out")
	}
}
