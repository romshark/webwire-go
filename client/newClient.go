package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"sync"

	webwire "github.com/qbeon/webwire-go"
	reqman "github.com/qbeon/webwire-go/requestManager"
)

// NewClient creates a new client instance.
// The new client will immediately begin connecting if autoconnect is enabled
func NewClient(
	serverAddress url.URL,
	implementation Implementation,
	opts Options,
	tlsConfig *tls.Config,
) (Client, error) {
	if implementation == nil {
		return nil, fmt.Errorf(
			"webwire client requires a client implementation, got nil",
		)
	}

	// Prepare server address
	if serverAddress.Scheme != "https" {
		serverAddress.Scheme = "http"
	}

	// Prepare configuration
	opts.SetDefaults()

	// Enable autoconnect by default
	autoconnect := autoconnectStatus(autoconnectEnabled)
	if opts.Autoconnect == webwire.Disabled {
		autoconnect = autoconnectDisabled
	}

	if tlsConfig != nil {
		tlsConfig = tlsConfig.Clone()
	}

	// Initialize new client
	newClt := &client{
		serverAddr:        serverAddress,
		tlsConfig:         tlsConfig,
		impl:              implementation,
		sessionInfoParser: opts.SessionInfoParser,
		status:            Disconnected,
		defaultReqTimeout: opts.DefaultRequestTimeout,
		reconnInterval:    opts.ReconnectionInterval,
		autoconnect:       autoconnect,
		sessionLock:       sync.RWMutex{},
		session:           nil,
		apiLock:           sync.RWMutex{},
		backReconn:        newDam(),
		connecting:        false,
		connectingLock:    sync.RWMutex{},
		connectLock:       sync.Mutex{},
		conn:              webwire.NewSocket(tlsConfig),
		readerClosing:     make(chan bool, 1),
		requestManager:    reqman.NewRequestManager(),
		warningLog:        opts.WarnLog,
		errorLog:          opts.ErrorLog,
	}

	if autoconnect == autoconnectEnabled {
		// Asynchronously connect to the server
		// immediately after initialization.
		// Call in another goroutine to prevent blocking
		// the constructor function caller.
		// Set timeout to zero, try indefinitely until connected
		go newClt.tryAutoconnect(context.Background(), 0)
	}

	return newClt, nil
}
