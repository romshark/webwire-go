package webwire

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

// NewServer creates a new headed WebWire server instance
// with a built-in HTTP server hosting it
func NewServer(
	implementation ServerImplementation,
	opts ServerOptions,
) (Server, error) {
	opts.SetDefaults()

	headlessInstance, err := NewHeadlessServer(implementation, opts)
	if err != nil {
		return nil, err
	}

	srv := headlessInstance.(*server)
	srv.options = opts

	// Initialize HTTP server
	srv.httpServer = &http.Server{
		Addr:    opts.Host,
		Handler: srv,
	}

	// Determine final address
	if opts.Host == "" {
		opts.Host = ":http"
	}

	// Initialize TCP/IP listener
	srv.listener, err = net.Listen("tcp", opts.Host)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up TCP/IP listener: %s", err)
	}

	srv.addr = url.URL{
		Scheme: "http",
		Host:   srv.listener.Addr().String(),
		Path:   "/",
	}

	return srv, nil
}

// NewServerSecure creates a new headed WebWire server instance
// with a built-in HTTPS server hosting it
func NewServerSecure(
	implementation ServerImplementation,
	opts ServerOptions,
	certFilePath,
	keyFilePath string,
	TLSConfig *tls.Config,
) (Server, error) {
	opts.SetDefaults()

	if TLSConfig == nil {
		TLSConfig = &tls.Config{}
	}

	headlessInstance, err := NewHeadlessServer(implementation, opts)
	if err != nil {
		return nil, err
	}

	srv := headlessInstance.(*server)
	srv.options = opts
	srv.certFilePath = certFilePath
	srv.keyFilePath = keyFilePath

	// Initialize HTTPS server
	srv.httpServer = &http.Server{
		Addr:      opts.Host,
		Handler:   srv,
		TLSConfig: TLSConfig,
	}

	// Determine final address
	if opts.Host == "" {
		opts.Host = ":http"
	}

	// Initialize TCP/IP listener
	srv.listener, err = net.Listen("tcp", opts.Host)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up TCP/IP listener: %s", err)
	}

	srv.addr = url.URL{
		Scheme: "https",
		Host:   srv.listener.Addr().String(),
		Path:   "/",
	}

	return srv, nil
}

// NewHeadlessServer creates a new headless WebWire server instance
// which relies on an external HTTP server to host it
func NewHeadlessServer(
	implementation ServerImplementation,
	opts ServerOptions,
) (instance HeadlessServer, err error) {
	if implementation == nil {
		panic(fmt.Errorf(
			"server instance requires an implementation, got nil",
		))
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
		addr:            url.URL{},
		options:         opts,
		shutdown:        false,
		shutdownRdy:     make(chan bool),
		currentOps:      0,
		opsLock:         &sync.Mutex{},
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
