package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// BenchmarkRequestC1_P1K benchmarks a request with a 1 kb payload on a single
// connection
func BenchmarkRequestC1_P1K(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 1024)
	msg := wwr.NewPayload(
		wwr.EncodingUtf8,
		payloadData,
	)

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return msg.Payload(), nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		client.connection.Request(context.Background(), nil, msg)
	}
}

// BenchmarkRequestC1_P1M benchmarks a request with a 1 mb payload on a single
// connection
func BenchmarkRequestC1_P1M(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 1024*1024*1024)
	msg := wwr.NewPayload(
		wwr.EncodingUtf8,
		payloadData,
	)

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return msg.Payload(), nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		client.connection.Request(context.Background(), nil, msg)
	}
}

// BenchmarkRequestC1_P1MBuffered benchmarks a request with a 1 mb payload on a
// single connection
func BenchmarkRequestC1_P1MBuffered(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 1024*1024)
	msg := wwr.NewPayload(
		wwr.EncodingUtf8,
		payloadData,
	)

	const bufferSize uint32 = 2 * 1024 * 1024

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return msg.Payload(), nil
			},
		},
		wwr.ServerOptions{
			WriteBufferSize: bufferSize,
			ReadBufferSize:  bufferSize,
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			WriteBufferSize: bufferSize,
			ReadBufferSize:  bufferSize,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		client.connection.Request(context.Background(), nil, msg)
	}
}

// BenchmarkRequestC1K_P1K benchmarks a request with a 1 kb payload on 1000
// concurrent connections
func BenchmarkRequestC1K_P1K(b *testing.B) {
	concurrentConnections := 1000

	// Preallocate the payload
	payloadData := make([]byte, 1024)
	msg := wwr.NewPayload(
		wwr.EncodingUtf8,
		payloadData,
	)

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return msg.Payload(), nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	clients := make([]*callbackPoweredClient, concurrentConnections)
	for i := 0; i < concurrentConnections; i++ {
		client := newCallbackPoweredClient(
			server.AddressURL(),
			wwrclt.Options{},
			callbackPoweredClientHooks{},
		)
		clients[i] = client
		if err := client.connection.Connect(); err != nil {
			panic(err)
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, c := range clients {
			client := c
			go func() {
				client.connection.Request(context.Background(), nil, msg)
			}()
		}
	}
}
