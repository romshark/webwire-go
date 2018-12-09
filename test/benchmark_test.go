package test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
)

// BenchmarkRequestC1_P16 benchmarks a request with a 1 kb payload on a single
// connection
func BenchmarkRequestC1_P16(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 16)
	msgPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	setup, err := SetupServer(
		&ServerImpl{
			Request: func(
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
		nil, // Use default transport implementation
	)
	if err != nil {
		log.Fatalf("couldn't setup server: %s", err)
	}

	// Initialize client
	sock, _, err := setup.NewClientSocket()
	if err != nil {
		panic(err)
	}
	ident := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	reply := message.NewMessage(32)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Write request
		writer, err := sock.GetWriter()
		if err != nil {
			panic(err)
		}

		message.WriteMsgRequest(
			writer,
			ident,
			nil, // No name
			msgPayload.Encoding,
			msgPayload.Data,
			false,
		)

		// Await reply
		if err := sock.Read(reply, time.Time{}); err != nil {
			panic(err)
		}
	}
}

// BenchmarkRequestC1_P1K benchmarks a request with a 1 kb payload on a single
// connection
func BenchmarkRequestC1_P1K(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 1024)
	msgPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	setup, err := SetupServer(
		&ServerImpl{
			Request: func(
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
		nil, // Use default transport implementation
	)
	if err != nil {
		log.Fatalf("couldn't setup server: %s", err)
	}

	// Initialize client
	sock, _, err := setup.NewClientSocket()
	if err != nil {
		panic(err)
	}
	ident := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	reply := message.NewMessage(1024 + 32)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Write request
		writer, err := sock.GetWriter()
		if err != nil {
			panic(err)
		}

		message.WriteMsgRequest(
			writer,
			ident,
			nil, // No name
			msgPayload.Encoding,
			msgPayload.Data,
			false,
		)

		// Await reply
		if err := sock.Read(reply, time.Time{}); err != nil {
			panic(err)
		}
	}
}

// BenchmarkRequestC1_P1M benchmarks a request with a 1 mb payload on a single
// connection
func BenchmarkRequestC1_P1M(b *testing.B) {
	// Preallocate the payload
	payloadData := make([]byte, 1024*1024)
	msgPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	setup, err := SetupServer(
		&ServerImpl{
			Request: func(
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
			MessageBufferSize: 1024*1024 + 32,
		},
		nil, // Use default transport implementation
	)
	if err != nil {
		log.Fatalf("couldn't setup server: %s", err)
	}

	// Initialize client
	sock, _, err := setup.NewClientSocket()
	if err != nil {
		panic(err)
	}
	ident := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	reply := message.NewMessage(1024*1024 + 32)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Write request
		writer, err := sock.GetWriter()
		if err != nil {
			panic(err)
		}

		message.WriteMsgRequest(
			writer,
			ident,
			nil, // No name
			msgPayload.Encoding,
			msgPayload.Data,
			false,
		)

		// Await reply
		if err := sock.Read(reply, time.Time{}); err != nil {
			panic(err)
		}
	}
}

// BenchmarkRequestC1K_P1K benchmarks a request with a 1 kb payload on 1000
// concurrent connections
func BenchmarkRequestC1K_P1K(b *testing.B) {
	concurrentConnections := 1000

	// Preallocate the payload
	payloadSize := uint32(1024)
	payloadData := make([]byte, payloadSize)
	msgPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     payloadData,
	}

	// Initialize a webwire server
	setup, err := SetupServer(
		&ServerImpl{
			Request: func(
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
			MessageBufferSize: payloadSize + 32,
		},
		nil, // Use default transport implementation
	)
	if err != nil {
		log.Fatalf("couldn't setup server: %s", err)
	}

	type Client struct {
		sock  wwr.Socket
		reply *message.Message
	}

	// Initialize clients
	clients := make([]Client, concurrentConnections)
	for i := 0; i < concurrentConnections; i++ {
		sock, _, err := setup.NewClientSocket()
		if err != nil {
			panic(err)
		}
		clients[i] = Client{
			sock:  sock,
			reply: message.NewMessage(payloadSize + 32),
		}
	}
	ident := []byte{0, 0, 0, 0, 0, 0, 0, 0}

	wg := sync.WaitGroup{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		wg.Add(concurrentConnections)
		for _, c := range clients {
			client := c
			go func() {
				// Write request
				writer, err := client.sock.GetWriter()
				if err != nil {
					panic(err)
				}

				message.WriteMsgRequest(
					writer,
					ident,
					nil, // No name
					msgPayload.Encoding,
					msgPayload.Data,
					false,
				)

				// Await reply
				if err := client.sock.Read(
					client.reply,
					time.Time{},
				); err != nil {
					panic(err)
				}

				wg.Done()
			}()
		}
		wg.Wait()
	}
}

// BenchmarkRequestSock_C1_P16 benchmarks a request with an 8 byte payload on a
// raw socket connection bypassing the client implementation
func BenchmarkRequestSock_C1_P16(b *testing.B) {
	requestName := ""
	const headerSize = 10
	const payloadSize = 16

	// Compose a binary request message
	payload := make([]byte, payloadSize)
	msgBytes := make([]byte, headerSize+len(requestName)+payloadSize)
	msgBytes[0] = message.MsgRequestBinary
	requestIdent := [8]byte{1, 1, 1, 1, 1, 1, 1, 1}
	copy(msgBytes[1:9], requestIdent[:])
	msgBytes[9] = byte(len(requestName))
	if len(requestName) > 0 {
		copy(msgBytes[headerSize:], []byte(requestName))
	}
	if payloadSize > 0 {
		copy(msgBytes[headerSize+len(requestName):], payload)
	}

	// Initialize a webwire server
	setup, err := SetupServer(
		&ServerImpl{
			Request: func(
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
		nil, // Use default transport implementation
	)
	if err != nil {
		log.Fatalf("couldn't setup server: %s", err)
	}

	// Setup client socket
	socket, _, err := setup.NewClientSocket()
	if err != nil {
		panic(err)
	}

	// Ignore the server configuration push-message
	replyMsg := message.NewMessage(1024)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Get writer
		writer, err := socket.GetWriter()
		if err != nil {
			panic(err)
		}

		// Write the message
		_, writeErr := writer.Write(msgBytes)
		if writeErr != nil {
			panic(writeErr)
		}

		// Flush buffer
		if err := writer.Close(); err != nil {
			panic(err)
		}

		// Await reply
		if err := socket.Read(replyMsg, time.Time{}); err != nil {
			panic(err)
		}
	}
}
