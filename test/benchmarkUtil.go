package test

import (
	"context"
	"fmt"

	wwr "github.com/qbeon/webwire-go"
)

// setupBenchmarkServer helps setting up and launching the benchmark server
// together with the hosting HTTP(S) server
func setupBenchmarkServer(
	implementation *serverImpl,
	options wwr.ServerOptions,
) wwr.Server {
	if implementation.onClientConnected == nil {
		implementation.onClientConnected = func(_ wwr.Connection) {}
	}
	if implementation.onClientDisconnected == nil {
		implementation.onClientDisconnected = func(_ wwr.Connection, _ error) {}
	}
	if implementation.onSignal == nil {
		implementation.onSignal = func(
			_ context.Context,
			_ wwr.Connection,
			_ wwr.Message,
		) {
		}
	}
	if implementation.onRequest == nil {
		implementation.onRequest = func(
			_ context.Context,
			_ wwr.Connection,
			_ wwr.Message,
		) (response wwr.Payload, err error) {
			return wwr.Payload{}, nil
		}
	}

	// Use default session manager if no specific one is defined
	if options.SessionManager == nil {
		options.SessionManager = newInMemSessManager()
	}

	// Use default localhost address
	options.Host = "127.0.0.1:0"

	server, err := wwr.NewServer(implementation, options)
	if err != nil {
		panic(fmt.Errorf("failed setting up server instance: %s", err))
	}

	// Run the server in a separate goroutine
	go func() {
		if err := server.Run(); err != nil {
			panic(fmt.Errorf("server failed: %s", err))
		}
	}()

	// Return reference to the server and the address its bound to
	return server
}
