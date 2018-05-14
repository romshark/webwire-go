package webwire

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

// NewServer creates a new WebWire server instance
func NewServer(
	implementation ServerImplementation,
	opts ServerOptions,
) (instance WebwireServer, err error) {
	if implementation == nil {
		panic(fmt.Errorf("A headed webwire server requires a server implementation, got nil"))
	}

	opts.SetDefaults()

	sessionsEnabled := false
	if opts.Sessions == Enabled {
		sessionsEnabled = true
	}

	srv := &server{
		impl:              implementation,
		sessionManager:    opts.SessionManager,
		sessionKeyGen:     opts.SessionKeyGenerator,
		sessionInfoParser: opts.SessionInfoParser,

		// State
		addr:            nil,
		shutdown:        false,
		shutdownRdy:     make(chan bool),
		currentOps:      0,
		opsLock:         sync.Mutex{},
		clients:         make([]*Client, 0),
		clientsLock:     &sync.Mutex{},
		sessionsEnabled: sessionsEnabled,
		sessionRegistry: newSessionRegistry(opts.MaxSessionConnections),

		// Internals
		connUpgrader: newConnUpgrader(),
		warnLog:      opts.WarnLog,
		errorLog:     opts.ErrorLog,
	}

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
