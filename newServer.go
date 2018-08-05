package webwire

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"golang.org/x/sync/semaphore"
)

// NewServer creates a new headed WebWire server instance
// with a built-in HTTP server hosting it
func NewServer(
	implementation ServerImplementation,
	opts ServerOptions,
) (instance Server, err error) {
	opts.SetDefaults()

	instance, err = NewHeadlessServer(implementation, opts)
	if err != nil {
		return nil, err
	}

	srv := instance.(*server)
	srv.options = opts

	// Initialize HTTP server
	srv.httpServer = &http.Server{
		Addr:    opts.Address,
		Handler: srv,
	}

	// Determine final address
	if opts.Address == "" {
		opts.Address = ":http"
	}

	// Initialize TCP/IP listener
	srv.listener, err = net.Listen("tcp", opts.Address)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up TCP/IP listener: %s", err)
	}

	srv.addr = srv.listener.Addr()

	return srv, nil
}

// NewHeadlessServer creates a new headless WebWire server instance
// which relies on an external HTTP server to host it
func NewHeadlessServer(
	implementation ServerImplementation,
	opts ServerOptions,
) (instance Server, err error) {
	if implementation == nil {
		panic(fmt.Errorf("A headed webwire server requires a server implementation, got nil"))
	}

	opts.SetDefaults()

	sessionsEnabled := false
	if opts.Sessions == Enabled {
		sessionsEnabled = true
	}

	return &server{
		impl:              implementation,
		sessionManager:    opts.SessionManager,
		sessionKeyGen:     opts.SessionKeyGenerator,
		sessionInfoParser: opts.SessionInfoParser,

		// State
		addr:        nil,
		options:     opts,
		stopping:    0,
		currentOps:  0,
		shutdownRdy: make(chan bool),
		handlerSlots: semaphore.NewWeighted(
			int64(opts.MaxConcurrentHandlers),
		),
		connections:     make([]*connection, 0),
		connectionsLock: &sync.Mutex{},
		sessionsEnabled: sessionsEnabled,
		sessionRegistry: newSessionRegistry(opts.MaxSessionConnections),

		// Internals
		connUpgrader: newConnUpgrader(),
		warnLog:      opts.WarnLog,
		errorLog:     opts.ErrorLog,
	}, nil
}
