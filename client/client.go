package client

import (
	"sync/atomic"

	webwire "github.com/qbeon/webwire-go"
	reqman "github.com/qbeon/webwire-go/requestManager"

	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const supportedProtocolVersion = "1.1"

// Client represents an instance of one of the servers clients
type Client struct {
	serverAddr     string
	isConnected    int32
	defaultTimeout time.Duration
	hooks          Hooks

	sessionLock sync.RWMutex
	session     *webwire.Session

	// Operation lock synchronizes concurrent access to:
	// Connect, RestoreSession, CloseSession and Close
	opLock   sync.Mutex
	connLock sync.Mutex
	conn     *websocket.Conn

	requestManager reqman.RequestManager

	// Loggers
	warningLog *log.Logger
	errorLog   *log.Logger
}

// NewClient creates a new client instance.
// TODO: return pointer to client instance, because copying the client object doesn't make sense
func NewClient(serverAddress string, opts Options) Client {
	opts.SetDefaults()

	return Client{
		serverAddress,
		0,
		opts.DefaultRequestTimeout,
		opts.Hooks,

		sync.RWMutex{},
		nil,

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
}

// IsConnected returns true if the client is connected to the server, otherwise false is returned
func (clt *Client) IsConnected() bool {
	return atomic.LoadInt32(&clt.isConnected) > 0
}

// Connect connects the client to the configured server and
// returns an error in case of a connection failure.
// Automatically tries to restore the previous session
func (clt *Client) Connect() (err error) {
	clt.opLock.Lock()
	defer clt.opLock.Unlock()

	if atomic.LoadInt32(&clt.isConnected) > 0 {
		return nil
	}

	if err := clt.verifyProtocolVersion(); err != nil {
		return err
	}

	connURL := url.URL{Scheme: "ws", Host: clt.serverAddr, Path: "/"}

	clt.connLock.Lock()
	clt.conn, _, err = websocket.DefaultDialer.Dial(connURL.String(), nil)
	if err != nil {
		// TODO: return typed error ConnectionFailure
		return fmt.Errorf("Could not connect: %s", err)
	}
	clt.connLock.Unlock()

	// Setup reader thread
	go func() {
		defer clt.close()
		for {
			_, message, err := clt.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					// Error while reading message
					clt.errorLog.Print("Failed reading message:", err)
					break
				} else {
					// Shutdown client due to clean disconnection
					break
				}
			}
			// Try to handle the message
			if err = clt.handleMessage(message); err != nil {
				clt.warningLog.Print("Failed handling message:", err)
			}
		}
	}()

	atomic.StoreInt32(&clt.isConnected, 1)

	// Read the current sessions key if there is any
	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		return nil
	}
	sessionKey := clt.session.Key
	clt.sessionLock.RUnlock()

	// Try to restore session if necessary
	restoredSession, err := clt.requestSessionRestoration([]byte(sessionKey))
	if err != nil {
		// Just log a warning and still return nil, even if session restoration failed,
		// because we only care about the connection establishment in this method
		clt.warningLog.Printf("Couldn't restore session on reconnection: %s", err)

		// Reset the session
		clt.sessionLock.Lock()
		clt.session = nil
		clt.sessionLock.Unlock()
		return nil
	}

	clt.sessionLock.Lock()
	clt.session = restoredSession
	clt.sessionLock.Unlock()
	return nil
}

// Request sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the request failed for some reason
func (clt *Client) Request(
	name string,
	payload webwire.Payload,
) (webwire.Payload, *webwire.Error) {
	reqType := webwire.MsgRequestBinary
	switch payload.Encoding {
	case webwire.EncodingUtf8:
		reqType = webwire.MsgRequestUtf8
	case webwire.EncodingUtf16:
		reqType = webwire.MsgRequestUtf16
	}
	return clt.sendRequest(reqType, name, payload, clt.defaultTimeout)
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
	if atomic.LoadInt32(&clt.isConnected) < 1 {
		return fmt.Errorf("Trying to send a signal on a disconnected client")
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
	if clt.session == nil {
		return nil
	}
	infoMap, converted := clt.session.Info.(map[string]interface{})
	if !converted {
		return nil
	}
	if value, exists := infoMap[key]; exists {
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
	clt.opLock.Lock()
	defer clt.opLock.Unlock()

	clt.sessionLock.RLock()
	if clt.session != nil {
		clt.sessionLock.RUnlock()
		return fmt.Errorf("Can't restore session if another one is already active")
	}
	clt.sessionLock.RUnlock()

	restoredSession, err := clt.requestSessionRestoration(sessionKey)
	if err != nil {
		// TODO: check for error types
		return fmt.Errorf("Session restoration request failed: %s", err)
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
	clt.opLock.Lock()
	defer clt.opLock.Unlock()

	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		return nil
	}
	clt.sessionLock.RUnlock()

	// Synchronize session closure to the server if connected
	if atomic.LoadInt32(&clt.isConnected) > 0 {
		if _, err := clt.sendNamelessRequest(
			webwire.MsgCloseSession,
			webwire.Payload{},
			clt.defaultTimeout,
		); err != nil {
			return fmt.Errorf("Session destruction request failed: %s", err)
		}
	}

	// Reset session locally after destroying it on the server
	clt.sessionLock.Lock()
	clt.session = nil
	clt.sessionLock.Unlock()

	return nil
}

// Close gracefully closes the connection.
// Does nothing if the client isn't connected
func (clt *Client) Close() {
	clt.opLock.Lock()
	defer clt.opLock.Unlock()
	clt.close()
}
