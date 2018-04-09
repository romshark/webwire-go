package client

import (
	"sync/atomic"

	webwire "github.com/qbeon/webwire-go"
	reqman "github.com/qbeon/webwire-go/requestManager"

	"fmt"
	"log"
	"sync"
	"time"
)

const supportedProtocolVersion = "1.3"

// Status represents the status of a client instance
type Status = int32

const (
	// StatDisconnected represents a temporarily disconnected client instance
	StatDisconnected Status = 0

	// StatConnected represents a connected client instance
	StatConnected Status = 1
)

// Autoconnect represents the activation of automatic reconnection
type Autoconnect = int32

const (
	// AutoconnectDisabled represents deactivated automatic reconnection
	AutoconnectDisabled = 0

	// AutoconnectEnabled represents activated automatic reconnection
	AutoconnectEnabled = 1
)

// Client represents an instance of one of the servers clients
type Client struct {
	serverAddr        string
	impl              Implementation
	sessionInfoParser webwire.SessionInfoParser
	status            Status
	defaultReqTimeout time.Duration
	reconnInterval    time.Duration
	autoconnect       int32

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

	connectLock   sync.Mutex
	conn          webwire.Socket
	readerClosing chan bool

	requestManager reqman.RequestManager

	// Loggers
	warningLog *log.Logger
	errorLog   *log.Logger
}

// NewClient creates a new client instance.
func NewClient(
	serverAddress string,
	implementation Implementation,
	opts Options,
) *Client {
	if implementation == nil {
		panic(fmt.Errorf(
			"A webwire client requires a client implementation, got nil",
		))
	}

	// Prepare configuration
	opts.SetDefaults()

	autoconnect := int32(1)
	if opts.Autoconnect == OptDisabled {
		autoconnect = int32(0)
	}

	// Initialize new client
	newClt := &Client{
		serverAddress,
		implementation,
		opts.SessionInfoParser,
		StatDisconnected,
		opts.DefaultRequestTimeout,
		opts.ReconnectionInterval,
		autoconnect,

		sync.RWMutex{},
		nil,

		sync.RWMutex{},
		newDam(),
		false,
		sync.RWMutex{},
		sync.Mutex{},
		newSocket(nil),
		make(chan bool, 1),

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

	if autoconnect == AutoconnectEnabled {
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
// Automatically tries to restore the previous session.
// Enables autoconnect if it was disabled
func (clt *Client) Connect() error {
	atomic.StoreInt32(&clt.autoconnect, AutoconnectEnabled)
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

	return clt.conn.Write(msgBytes)
}

// Session returns an exact copy of the session object or nil if there's no
// session currently assigned to this client
func (clt *Client) Session() *webwire.Session {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil {
		return nil
	}
	clone := &webwire.Session{
		Key:      clt.session.Key,
		Creation: clt.session.Creation,
	}
	if clt.session.Info != nil {
		clone.Info = clt.session.Info.Copy()
	}
	return clone
}

// SessionInfo returns a copy of the session info field value
// in the form of an empty interface to be casted to either concrete type
func (clt *Client) SessionInfo(fieldName string) interface{} {
	clt.sessionLock.RLock()
	defer clt.sessionLock.RUnlock()
	if clt.session == nil || clt.session.Info == nil {
		return nil
	}
	return clt.session.Info.Value(fieldName)
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

	atomic.StoreInt32(&clt.autoconnect, AutoconnectDisabled)

	if atomic.LoadInt32(&clt.status) != StatConnected {
		return
	}
	if err := clt.conn.Close(); err != nil {
		clt.errorLog.Printf("Failed closing connection: %s", err)
	}

	// Wait for the reader goroutine to die before returning
	<-clt.readerClosing
}
