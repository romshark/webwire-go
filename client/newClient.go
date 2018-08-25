package client

import (
	"context"
	"fmt"
	"sync"

	webwire "github.com/qbeon/webwire-go"
	reqman "github.com/qbeon/webwire-go/requestManager"
)

// NewClient creates a new client instance.
// The new client will immediately begin connecting if autoconnect is enabled
func NewClient(
	serverAddress string,
	implementation Implementation,
	opts Options,
) Client {
	if implementation == nil {
		panic(fmt.Errorf(
			"A webwire client requires a client implementation, got nil",
		))
	}

	// Prepare configuration
	opts.SetDefaults()

	// Enable autoconnect by default
	autoconnect := autoconnectStatus(autoconnectEnabled)
	if opts.Autoconnect == webwire.Disabled {
		autoconnect = autoconnectDisabled
	}

	// Initialize new client
	newClt := &client{
		serverAddr:        serverAddress,
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
		conn:              webwire.NewSocket(),
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

	return newClt
}
