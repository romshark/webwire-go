package client

import (
	"sync/atomic"

	webwire "github.com/qbeon/webwire-go"
	reqman "github.com/qbeon/webwire-go/requestManager"

	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const supportedProtocolVersion = "1.0"

// Client represents an instance of one of the servers clients
type Client struct {
	serverAddr     string
	isConnected    int32
	defaultTimeout time.Duration
	hooks          Hooks
	session        *webwire.Session

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

// NewClient creates a new disconnected client instance.
func NewClient(
	serverAddr string,
	hooks Hooks,
	defaultTimeout time.Duration,
	warningLogWriter io.Writer,
	errorLogWriter io.Writer,
) Client {
	hooks.SetDefaults()

	return Client{
		serverAddr,
		0,
		defaultTimeout,
		hooks,
		nil,

		sync.Mutex{},
		sync.Mutex{},
		nil,

		reqman.NewRequestManager(),

		log.New(
			warningLogWriter,
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		log.New(
			errorLogWriter,
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

	// Try to restore session if necessary
	if clt.session == nil {
		return nil
	}
	restoredSession, err := clt.requestSessionRestoration([]byte(clt.session.Key))
	if err != nil {
		// Just log a warning and still return nil, even if session restoration failed,
		// because we only care about the connection establishment in this method
		clt.warningLog.Printf("Couldn't restore session on reconnection: %s", err)

		// Reset the session
		clt.session = nil
	}
	clt.session = restoredSession

	return nil
}

// Request sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the request failed for some reason
func (clt *Client) Request(payload []byte) ([]byte, *webwire.Error) {
	return clt.sendRequest(webwire.MsgRequest, payload, clt.defaultTimeout)
}

// TimedRequest sends a request containing the given payload to the server
// and asynchronously returns the servers reply
// blocking the calling goroutine.
// Returns an error if the given timeout was exceeded awaiting the response
// or another failure occurred
func (clt *Client) TimedRequest(
	payload []byte,
	timeout time.Duration,
) ([]byte, *webwire.Error) {
	return clt.sendRequest(webwire.MsgRequest, payload, timeout)
}

// Signal sends a signal containing the given payload to the server
func (clt *Client) Signal(payload []byte) error {
	if atomic.LoadInt32(&clt.isConnected) < 1 {
		return fmt.Errorf("Trying to send a signal on a disconnected client")
	}

	var msg bytes.Buffer
	msg.WriteRune(webwire.MsgSignal)
	msg.Write(payload)
	clt.connLock.Lock()
	defer clt.connLock.Unlock()
	return clt.conn.WriteMessage(websocket.TextMessage, msg.Bytes())
}

// Session returns information about the current session
func (clt *Client) Session() webwire.Session {
	clt.opLock.Lock()
	defer clt.opLock.Unlock()

	if clt.session == nil {
		return webwire.Session{}
	}
	return *clt.session
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

	if clt.session != nil {
		return fmt.Errorf("Can't restore session if another one is already active")
	}

	restoredSession, err := clt.requestSessionRestoration(sessionKey)
	if err != nil {
		// TODO: check for error types
		return fmt.Errorf("Session restoration request failed: %s", err)
	}
	clt.session = restoredSession

	return nil
}

// CloseSession closes the currently active session
// and synchronizes the closure to the server if connected.
// If the client is not connected then the synchronization is skipped.
// Does nothing if there's no active session
func (clt *Client) CloseSession() error {
	clt.opLock.Lock()
	defer clt.opLock.Unlock()

	if clt.session == nil {
		return nil
	}

	// Synchronize session closure to the server if connected
	if atomic.LoadInt32(&clt.isConnected) > 0 {
		if _, err := clt.sendRequest(
			webwire.MsgCloseSession,
			nil,
			clt.defaultTimeout,
		); err != nil {
			return fmt.Errorf("Session destruction request failed: %s", err)
		}
	}

	// Reset session locally after destroying it on the server
	clt.session = nil

	return nil
}

// Close gracefully closes the connection.
// Does nothing if the client isn't connected
func (clt *Client) Close() {
	clt.opLock.Lock()
	defer clt.opLock.Unlock()
	clt.close()
}
