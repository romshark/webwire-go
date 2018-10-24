package client

import (
	"context"
	"crypto/tls"
	"sync/atomic"

	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	webwire "github.com/qbeon/webwire-go"
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
	reqman "github.com/qbeon/webwire-go/requestManager"
)

const supportedProtocolVersion = "1.5"

// Status represents the status of a client instance
type Status = int32

const (
	// Disabled represents a permanent connection loss
	Disabled Status = 0

	// Disconnected represents a temporarily connection loss
	Disconnected

	// Connected represents a normal connection
	Connected Status = 2
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
	serverAddr        url.URL
	tlsConfig         *tls.Config
	impl              Implementation
	sessionInfoParser webwire.SessionInfoParser
	status            Status
	defaultReqTimeout time.Duration
	reconnInterval    time.Duration
	autoconnect       autoconnectStatus

	sessionLock sync.RWMutex
	session     *webwire.Session

	// The API lock synchronizes concurrent access
	// to the public client interface.
	// Request, and Signal methods are locked with a shared lock
	// because performing multiple requests and/or signals simultaneously
	// is fine.
	// The Connect, RestoreSession, CloseSession and Close methods
	// are locked exclusively because they should temporarily block
	// any other interaction with this client instance.
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
	conn          webwire.Socket
	readerClosing chan bool
	endpointMeta  endpointMeta

	heartbeat      heartbeat
	requestManager reqman.RequestManager

	// Loggers
	warningLog *log.Logger
	errorLog   *log.Logger
}

// Status returns the current client status
// which is either disabled, disconnected or connected.
// The client is considered disabled when it was manually closed
// through client.Close, while disconnected is considered
// a temporary connection loss.
// A disabled client won't autoconnect until enabled again.
func (clt *client) Status() Status {
	return atomic.LoadInt32(&clt.status)
}

// Connect connects the client to the configured server and
// returns an error in case of a connection failure.
// Automatically tries to restore the previous session.
// Enables autoconnect if it was disabled
func (clt *client) Connect() error {
	if atomic.LoadInt32(&clt.autoconnect) == autoconnectDeactivated {
		atomic.StoreInt32(&clt.autoconnect, autoconnectEnabled)
	}
	return clt.connect()
}

// Request sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the request failed for some reason
func (clt *client) Request(
	ctx context.Context,
	name string,
	payload webwire.Payload,
) (webwire.Payload, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	clt.apiLock.RLock()
	defer clt.apiLock.RUnlock()

	if err := clt.tryAutoconnect(ctx, clt.defaultReqTimeout); err != nil {
		return nil, err
	}

	return clt.sendRequest(
		ctx,
		scanPayloadEncoding(payload),
		name,
		payload,
		clt.defaultReqTimeout,
	)
}

// Signal sends a signal containing the given payload to the server
func (clt *client) Signal(name string, payload webwire.Payload) error {
	clt.apiLock.RLock()
	defer clt.apiLock.RUnlock()

	if err := clt.tryAutoconnect(
		context.Background(),
		clt.defaultReqTimeout,
	); err != nil {
		return err
	}

	// Require either a name or a payload or both
	if len(name) < 1 && (payload == nil || len(payload.Data()) < 1) {
		return webwire.NewProtocolErr(
			fmt.Errorf("Invalid request, request message requires " +
				"either a name, a payload or both but is missing both",
			),
		)
	}

	// Initialize payload encoding & data
	var encoding webwire.PayloadEncoding
	var data []byte
	if payload != nil {
		encoding = payload.Encoding()
		data = payload.Data()
	}

	if err := clt.conn.Write(msg.NewSignalMessage(
		name,
		encoding,
		data,
	)); err != nil {
		return err
	}

	clt.heartbeat.reset()

	return nil
}

// Session returns an exact copy of the session object or nil if there's no
// session currently assigned to this client
func (clt *client) Session() *webwire.Session {
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

// RestoreSession tries to restore the previously opened session.
// Fails if a session is currently already active
func (clt *client) RestoreSession(sessionKey []byte) error {
	clt.apiLock.Lock()
	defer clt.apiLock.Unlock()

	clt.sessionLock.RLock()
	if clt.session != nil {
		clt.sessionLock.RUnlock()
		return fmt.Errorf(
			"Can't restore session if another one is already active",
		)
	}
	clt.sessionLock.RUnlock()

	if err := clt.tryAutoconnect(
		context.Background(),
		clt.defaultReqTimeout,
	); err != nil {
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

// CloseSession disables the currently active session
// and acknowledges the server if connected.
// The session will be destroyed if this is it's last connection remaining.
// If the client is not connected then the synchronization is skipped.
// Does nothing if there's no active session
func (clt *client) CloseSession() error {
	clt.apiLock.Lock()
	defer clt.apiLock.Unlock()

	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		return nil
	}
	clt.sessionLock.RUnlock()

	// Synchronize session closure to the server if connected
	if atomic.LoadInt32(&clt.status) == Connected {
		if _, err := clt.sendNamelessRequest(
			context.Background(),
			msg.MsgCloseSession,
			pld.Payload{},
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
func (clt *client) Close() {
	clt.apiLock.Lock()
	defer clt.apiLock.Unlock()

	// Disable autoconnect and set status to disabled
	if atomic.LoadInt32(&clt.autoconnect) != autoconnectDisabled {
		atomic.StoreInt32(&clt.autoconnect, autoconnectDeactivated)
	}

	if atomic.LoadInt32(&clt.status) != Connected {
		atomic.StoreInt32(&clt.status, Disabled)
		return
	}
	atomic.StoreInt32(&clt.status, Disabled)

	if err := clt.conn.Close(); err != nil {
		clt.errorLog.Printf("Failed closing connection: %s", err)
	}

	// Wait for the reader goroutine to die before returning
	<-clt.readerClosing
}
