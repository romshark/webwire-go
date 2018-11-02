package client

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/msgbuf"
	reqman "github.com/qbeon/webwire-go/requestManager"
)

// NewClient creates a new client instance.
// The new client will immediately begin connecting if autoconnect is enabled
func NewClient(
	serverAddress url.URL,
	implementation Implementation,
	options Options,
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
	if err := options.Prepare(); err != nil {
		return nil, err
	}

	// Enable autoconnect by default
	autoconnect := autoconnectStatus(autoconnectEnabled)
	if options.Autoconnect == webwire.Disabled {
		autoconnect = autoconnectDisabled
	}

	conn := webwire.NewFasthttpSocket(
		options.TLSConfig,
		options.DialingTimeout,
	)

	// Initialize new client
	newClt := &client{
		serverAddr:        serverAddress,
		options:           options,
		impl:              implementation,
		status:            Disconnected,
		autoconnect:       autoconnect,
		sessionLock:       sync.RWMutex{},
		session:           nil,
		apiLock:           sync.RWMutex{},
		backReconn:        newDam(),
		connecting:        false,
		connectingLock:    sync.RWMutex{},
		connectLock:       sync.Mutex{},
		conn:              conn,
		readerClosing:     make(chan bool, 1),
		heartbeat:         newHeartbeat(conn, options.ErrorLog),
		requestManager:    reqman.NewRequestManager(),
		messageBufferPool: msgbuf.NewSyncPool(options.MessageBufferSize, 0),
	}

	if autoconnect == autoconnectEnabled {
		// Asynchronously connect to the server
		// immediately after initialization.
		// Call in another goroutine to prevent blocking
		// the constructor function caller.
		// Set timeout to zero, try indefinitely until connected
		go newClt.tryAutoconnect(context.Background(), false)
	}

	return newClt, nil
}
