package client

import (
	"sync/atomic"
	"time"

	"net/url"
	"sync"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	reqman "github.com/qbeon/webwire-go/requestManager"
)

// Status represents the status of a client instance
type Status = int32

const (
	// StatusDisabled represents a permanent connection loss
	StatusDisabled Status = 0

	// StatusDisconnected represents a temporarily connection loss
	StatusDisconnected Status = 1

	// StatusConnected represents a normal connection
	StatusConnected Status = 2
)

// autoconnectStatus represents the activation of auto-reconnection
type autoconnectStatus = int32

const (
	// autoconnectDisabled represents permanently disabled auto-reconnection
	autoconnectDisabled = 0

	// autoconnectDeactivated represents deactivated auto-reconnection
	autoconnectDeactivated = 1

	// autoconnectEnabled represents activated auto-reconnection
	autoconnectEnabled = 2
)

// client represents an instance of one of the servers clients
type client struct {
	serverAddr   url.URL
	options      Options
	impl         Implementation
	dialingTimer *time.Timer
	autoconnect  autoconnectStatus
	statusLock   *sync.Mutex
	status       Status

	sessionLock sync.RWMutex
	session     *webwire.Session

	// The API lock synchronizes concurrent access to the public client
	// interface. Request, and Signal methods are locked with a shared lock
	// because performing multiple requests and/or signals simultaneously is
	// fine. The Connect, RestoreSession, CloseSession and Close methods are
	// locked exclusively because they should temporarily block any other
	// interaction with this client instance.
	apiLock sync.RWMutex

	// backReconn is a dam that's flushed
	// when the client establishes a connection.
	backReconn *dam
	// connecting prevents multiple autoconnection attempts from spawning
	// superfluous multiple goroutines each polling the server
	connecting bool
	// connectingLock protects the connecting flag from concurrent access
	connectingLock sync.RWMutex

	connectLock   sync.Mutex
	conn          webwire.ClientSocket
	readerClosing chan bool

	heartbeat      heartbeat
	requestManager reqman.RequestManager
	messagePool    message.Pool
}

// Status returns the current client status
// which is either disabled, disconnected or connected.
// The client is considered disabled when it was manually closed
// through client.Close, while disconnected is considered
// a temporary connection loss.
// A disabled client won't autoconnect until enabled again.
func (clt *client) Status() Status {
	clt.statusLock.Lock()
	status := clt.status
	clt.statusLock.Unlock()
	return status
}

// Connect connects the client to the configured server and
// returns an error in case of a connection failure.
// Automatically tries to restore the previous session.
// Enables autoconnect if it was disabled
func (clt *client) Connect() error {
	atomic.CompareAndSwapInt32(
		&clt.autoconnect,
		autoconnectDeactivated,
		autoconnectEnabled,
	)

	return clt.connect()
}

// Session returns an exact copy of the session object or nil if there's no
// session currently assigned to this client
func (clt *client) Session() *webwire.Session {
	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		return nil
	}
	clone := &webwire.Session{
		Key:      clt.session.Key,
		Creation: clt.session.Creation,
	}
	if clt.session.Info != nil {
		clone.Info = clt.session.Info.Copy()
	}
	clt.sessionLock.RUnlock()
	return clone
}

// SessionInfo returns a copy of the session info field value
// in the form of an empty interface to be casted to either concrete type
func (clt *client) SessionInfo(fieldName string) interface{} {
	clt.sessionLock.RLock()
	if clt.session == nil || clt.session.Info == nil {
		clt.sessionLock.RUnlock()
		return nil
	}
	val := clt.session.Info.Value(fieldName)
	clt.sessionLock.RUnlock()
	return val
}

// PendingRequests returns the number of currently pending requests
func (clt *client) PendingRequests() int {
	return clt.requestManager.PendingRequests()
}

// Close gracefully closes the connection and disables the client.
// A disabled client won't autoconnect until enabled again.
func (clt *client) Close() {
	// Apply exclusive lock
	clt.apiLock.Lock()

	clt.close()

	// Disable autoconnect and set status to disabled
	atomic.CompareAndSwapInt32(
		&clt.autoconnect,
		autoconnectEnabled,
		autoconnectDeactivated,
	)

	clt.apiLock.Unlock()
}

// close closes the connection if it's open
func (clt *client) close() {
	// Apply exclusive lock
	clt.statusLock.Lock()

	if clt.status != StatusConnected {
		clt.status = StatusDisabled
		clt.statusLock.Unlock()
		return
	}
	clt.status = StatusDisabled
	clt.statusLock.Unlock()

	if err := clt.conn.Close(); err != nil {
		clt.options.ErrorLog.Printf("Failed closing connection: %s", err)
	}

	// Wait for the reader goroutine to die before returning
	<-clt.readerClosing
}

// setStatus atomically sets the status
func (clt *client) setStatus(newStatus Status) {
	clt.statusLock.Lock()
	clt.status = newStatus
	clt.statusLock.Unlock()
}
