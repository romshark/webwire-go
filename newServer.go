package webwire

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/qbeon/webwire-go/message"
)

// NewServer creates a new webwire server instance
func NewServer(
	implementation ServerImplementation,
	opts ServerOptions,
	transport Transport,
) (instance Server, err error) {
	if implementation == nil {
		return nil, errors.New("missing server implementation")
	}

	if transport == nil {
		return nil, errors.New("missing transport layer implementation")
	}

	if err := opts.Prepare(); err != nil {
		return nil, err
	}

	sessionsEnabled := false
	if opts.Sessions == Enabled {
		sessionsEnabled = true
	}

	// Prepare the configuration push message for the webwire accept handshake
	configMsg, err := message.NewAcceptConfMessage(
		message.ServerConfiguration{
			MajorProtocolVersion: 2,
			MinorProtocolVersion: 0,
			ReadTimeout:          opts.ReadTimeout,
			MessageBufferSize:    opts.MessageBufferSize,
			SubProtocolName:      opts.SubProtocolName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"couldn't initialize server configuration-push message: %s",
			err,
		)
	}

	// Initialize the webwire server
	srv := &server{
		transport:         transport,
		impl:              implementation,
		sessionManager:    opts.SessionManager,
		sessionKeyGen:     opts.SessionKeyGenerator,
		sessionInfoParser: opts.SessionInfoParser,
		addr:              url.URL{},
		options:           opts,
		configMsg:         configMsg,
		shutdown:          false,
		shutdownRdy:       make(chan bool),
		currentOps:        0,
		opsLock:           &sync.Mutex{},
		connections:       make([]*connection, 0),
		connectionsLock:   &sync.Mutex{},
		sessionsEnabled:   sessionsEnabled,
		messagePool:       message.NewSyncPool(opts.MessageBufferSize, 1024),
		warnLog:           opts.WarnLog,
		errorLog:          opts.ErrorLog,
	}

	srv.sessionRegistry = newSessionRegistry(
		opts.MaxSessionConnections,
		func(sessionKey string) {
			if err := srv.sessionManager.OnSessionClosed(
				sessionKey,
			); err != nil {
				srv.errorLog.Printf("session registry ")
			}
		},
	)

	// Initialize the transport layer
	if err := transport.Initialize(
		opts,
		srv.isShuttingDown,
		srv.handleConnection,
	); err != nil {
		return nil, fmt.Errorf("couldn't initialize transport layer: %s", err)
	}

	return srv, nil
}
