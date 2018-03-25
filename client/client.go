package client

import (
	"sync/atomic"

	webwire "github.com/qbeon/webwire-go"
	reqman "github.com/qbeon/webwire-go/requestManager"

	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const supportedProtocolVersion = "1.2"

// Status represents the status of a client instance
type Status = int32

const (
	// StatDisabled represents a client instance that has been manually closed, thus disabled
	StatDisabled Status = 0

	// StatDisconnected represents a temporarily disconnected client instance
	StatDisconnected Status = 1

	// StatConnected represents a connected client instance
	StatConnected Status = 2
)

// Client represents an instance of one of the servers clients
type Client struct {
	serverAddr        string
	status            Status
	defaultReqTimeout time.Duration
	reconnInterval    time.Duration
	autoconnect       bool
	hooks             Hooks

	sessionLock sync.RWMutex
	session     *webwire.Session

	// The API lock synchronizes concurrent access to the public client interface.
	// Request, TimedRequest and Signal methods are locked with a shared lock
	// because performing multiple requests and/or signals simultaneously is fine.
	// The Connect, RestoreSession, CloseSession and Close methods are locked exclusively
	// because they should temporarily block any other interaction with this client instance.
	apiLock sync.RWMutex

	// backReconn is a dam that's flushed when the client establishes a connection.
	backReconn *dam
	// connecting prevents multiple autoconnection attempts from spawning
	// superfluous multiple goroutines each polling the server
	connecting bool
	// connectingLock protects the connecting flag from concurrent access
	connectingLock sync.RWMutex

	connectLock sync.Mutex
	connLock    sync.Mutex
	conn        *websocket.Conn

	requestManager reqman.RequestManager

	// Loggers
	warningLog *log.Logger
	errorLog   *log.Logger
}

// NewClient creates a new client instance.
func NewClient(serverAddress string, opts Options) *Client {
	// Prepare configuration
	opts.SetDefaults()

	autoconnect := true
	if opts.Autoconnect == OptDisabled {
		autoconnect = false
	}

	// Initialize new client
	newClt := &Client{
		serverAddress,
		StatDisconnected,
		opts.DefaultRequestTimeout,
		opts.ReconnectionInterval,
		autoconnect,
		opts.Hooks,

		sync.RWMutex{},
		nil,

		sync.RWMutex{},
		newDam(),
		false,
		sync.RWMutex{},
		sync.Mutex{},
		sync.Mutex{},
		nil,

		reqman.NewRequestManager(),

		log.New(
			opts.WarnLog,
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		log.New(
			opts.ErrorLog,
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		),
	}

	if autoconnect {
		// Asynchronously connect to the server immediately after initialization.
		// Call in another goroutine to not block the contructor function caller.
		// Set timeout to zero, try indefinitely until connected.
		go newClt.tryAutoconnect(0)
	}

	return newClt
}

// Status returns the current client status
// which is either disabled, disconnected or connected.
// The client is considered disabled when it was manually closed through client.Close,
// while disconnected is considered a temporary connection loss.
// A disabled client won't autoconnect until enabled again.
func (clt *Client) Status() Status {
	return atomic.LoadInt32(&clt.status)
}

// Connect connects the client to the configured server and
// returns an error in case of a connection failure.
// Automatically tries to restore the previous session
func (clt *Client) Connect() error {
	return clt.connect()
}

// Request sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the request failed for some reason
func (clt *Client) Request(
	name string,
	payload webwire.Payload,
) (webwire.Payload, error) {
	clt.apiLock.RLock()
	defer clt.apiLock.RUnlock()

	if err := clt.tryAutoconnect(clt.defaultReqTimeout); err != nil {
		return webwire.Payload{}, err
	}

	reqType := webwire.MsgRequestBinary
	switch payload.Encoding {
	case webwire.EncodingUtf8:
		reqType = webwire.MsgRequestUtf8
	case webwire.EncodingUtf16:
		reqType = webwire.MsgRequestUtf16
	}
	return clt.sendRequest(reqType, name, payload, clt.defaultReqTimeout)
}

// TimedRequest sends a request containing the given payload to the server
// and asynchronously returns the servers reply
// blocking the calling goroutine.
// Returns an error if the given timeout was exceeded awaiting the response
// or another failure occurred
func (clt *Client) TimedRequest(
	name string,
	payload webwire.Payload,
	timeout time.Duration,
) (webwire.Payload, error) {
	clt.apiLock.RLock()
	defer clt.apiLock.RUnlock()

	if err := clt.tryAutoconnect(timeout); err != nil {
		return webwire.Payload{}, err
	}

	reqType := webwire.MsgRequestBinary
	switch payload.Encoding {
	case webwire.EncodingUtf8:
		reqType = webwire.MsgRequestUtf8
	case webwire.EncodingUtf16:
		reqType = webwire.MsgRequestUtf16
	}
	return clt.sendRequest(reqType, name, payload, timeout)
}

// Signal sends a signal containing the given payload to the server
func (clt *Client) Signal(name string, payload webwire.Payload) error {
	clt.apiLock.RLock()
	defer clt.apiLock.RUnlock()

	if err := clt.connect(); err != nil {
		return err
	}

	msgBytes := webwire.NewSignalMessage(name, payload)

	clt.connLock.Lock()
	defer clt.connLock.Unlock()
	return clt.conn.WriteMessage(websocket.BinaryMessage, msgBytes)
}

// Session returns information about the current session
func (clt *Client) Session() webwire.Session {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil {
		return webwire.Session{}
	}
	return *clt.session
}

// SessionInfo returns the value of a session info field identified by the given key
// in the form of an empty interface that could be casted to either a string, bool, float64 number
// a map[string]interface{} object or an []interface{} array according to JSON data types.
// Returns nil if either there's no session or if the given field doesn't exist.
func (clt *Client) SessionInfo(key string) interface{} {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil || clt.session.Info == nil {
		return nil
	}
	if value, exists := clt.session.Info[key]; exists {
		return value
	}
	return nil
}

// PendingRequests returns the number of currently pending requests
func (clt *Client) PendingRequests() int {
	return clt.requestManager.PendingRequests()
}

// RestoreSession tries to restore the previously opened session.
// Fails if a session is currently already active
func (clt *Client) RestoreSession(sessionKey []byte) error {
	clt.apiLock.Lock()
	defer clt.apiLock.Unlock()

	clt.sessionLock.RLock()
	if clt.session != nil {
		clt.sessionLock.RUnlock()
		return fmt.Errorf("Can't restore session if another one is already active")
	}
	clt.sessionLock.RUnlock()

	if err := clt.tryAutoconnect(clt.defaultReqTimeout); err != nil {
		return err
	}

	restoredSession, err := clt.requestSessionRestoration(sessionKey)
	if err != nil {
		return err
	}
	clt.sessionLock.Lock()
	clt.session = restoredSession
	clt.sessionLock.Unlock()

	return nil
}

// CloseSession closes the currently active session
// and synchronizes the closure to the server if connected.
// If the client is not connected then the synchronization is skipped.
// Does nothing if there's no active session
func (clt *Client) CloseSession() error {
	clt.apiLock.Lock()
	defer clt.apiLock.Unlock()

	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		return nil
	}
	clt.sessionLock.RUnlock()

	// Synchronize session closure to the server if connected
	if atomic.LoadInt32(&clt.status) == StatConnected {
		if _, err := clt.sendNamelessRequest(
			webwire.MsgCloseSession,
			webwire.Payload{},
			clt.defaultReqTimeout,
		); err != nil {
			return err
		}
	}

	// Reset session locally after destroying it on the server
	clt.sessionLock.Lock()
	clt.session = nil
	clt.sessionLock.Unlock()

	return nil
}

// Close gracefully closes the connection and disables the client.
// A disabled client won't autoconnect until enabled again.
func (clt *Client) Close() {
	clt.apiLock.Lock()
	defer clt.apiLock.Unlock()
	clt.close()
}
