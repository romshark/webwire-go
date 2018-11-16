package client

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	reqman "github.com/qbeon/webwire-go/requestManager"
)

// NewClient creates a new client instance.
// The new client will immediately begin connecting if autoconnect is enabled
func NewClient(
	serverAddress url.URL,
	implementation Implementation,
	options Options,
	transport webwire.ClientTransport,
) (Client, error) {
	if implementation == nil {
		return nil, fmt.Errorf("missing client implementation")
	}

	if transport == nil {
		return nil, fmt.Errorf("missing client transport layer implementation")
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

	// Initialize socket
	conn, err := transport.NewSocket(options.DialingTimeout)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize socket: %s", err)
	}

	dialingTimer := time.NewTimer(0)
	<-dialingTimer.C

	// Initialize new client
	newClt := &client{
		serverAddr:     serverAddress,
		options:        options,
		impl:           implementation,
		dialingTimer:   dialingTimer,
		autoconnect:    autoconnect,
		statusLock:     &sync.Mutex{},
		status:         StatusDisconnected,
		sessionLock:    sync.RWMutex{},
		session:        nil,
		apiLock:        sync.RWMutex{},
		backReconn:     newDam(),
		connecting:     false,
		connectingLock: sync.RWMutex{},
		connectLock:    sync.Mutex{},
		conn:           conn,
		readerClosing:  make(chan bool, 1),
		heartbeat:      newHeartbeat(conn, options.ErrorLog),
		requestManager: reqman.NewRequestManager(),
		messagePool:    message.NewSyncPool(options.MessageBufferSize, 1024),
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
