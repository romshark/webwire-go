package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// BenchmarkRequestC1_P16 benchmarks a request with a 1 kb payload on a single
// connection
func BenchmarkRequestC1_P16(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 16)
	msg := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return wwr.Payload{
					Encoding: msg.PayloadEncoding(),
					Data:     msg.Payload(),
				}, nil
			},
		},
		wwr.ServerOptions{
			MessageBufferSize: 1024,
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			MessageBufferSize: 1024,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		reply, err := client.connection.Request(context.Background(), nil, msg)
		if err != nil {
			panic(err)
		}
		reply.Close()
	}
}

// BenchmarkRequestC1_P1K benchmarks a request with a 1 kb payload on a single
// connection
func BenchmarkRequestC1_P1K(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 1024)
	msg := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return wwr.Payload{
					Encoding: msg.PayloadEncoding(),
					Data:     msg.Payload(),
				}, nil
			},
		},
		wwr.ServerOptions{
			MessageBufferSize: 2048,
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			MessageBufferSize: 2048,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		reply, err := client.connection.Request(context.Background(), nil, msg)
		if err != nil {
			panic(err)
		}
		reply.Close()
	}
}

// BenchmarkRequestC1_P1M benchmarks a request with a 1 mb payload on a single
// connection
func BenchmarkRequestC1_P1M(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 1024*1024)
	msg := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return wwr.Payload{
					Encoding: msg.PayloadEncoding(),
					Data:     msg.Payload(),
				}, nil
			},
		},
		wwr.ServerOptions{
			MessageBufferSize: 1024*1024 + 1024,
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			MessageBufferSize: 1024*1024 + 1024,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		reply, err := client.connection.Request(context.Background(), nil, msg)
		if err != nil {
			panic(err)
		}
		reply.Close()
	}
}

// BenchmarkRequestC1K_P1K benchmarks a request with a 1 kb payload on 1000
// concurrent connections
func BenchmarkRequestC1K_P1K(b *testing.B) {
	concurrentConnections := 1000

	// Preallocate the payload
	payloadData := make([]byte, 1024)
	msg := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	server := setupBenchmarkServer(
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				return wwr.Payload{
					Encoding: msg.PayloadEncoding(),
					Data:     msg.Payload(),
				}, nil
			},
		},
		wwr.ServerOptions{
			MessageBufferSize: 2048,
		},
	)

	// Initialize client
	clients := make([]*callbackPoweredClient, concurrentConnections)
	for i := 0; i < concurrentConnections; i++ {
		client := newCallbackPoweredClient(
			server.Address(),
			wwrclt.Options{
				MessageBufferSize: 2048,
			},
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
				reply, err := client.connection.Request(context.Background(), nil, msg)
				if err != nil {
					panic(err)
				}
				reply.Close()
			}()
		}
	}
}
