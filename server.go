package webwire

import (
	"io"
	"log"
	"fmt"
	"net"
	"time"
	"sync"
	"context"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
)

const protocolVersion = "1.0"

// OnOptions is an optional hook.
// It's invoked when the websocket endpoint is examined by the client
// using the HTTP OPTION method.
type OnOptions func(resp http.ResponseWriter)

// OnClientConnected is an optional hook.
// It's invoked when a new client establishes a connection to the server
type OnClientConnected func(client *Client)

// OnClientDisconnected is an optional hook.
// It's invoked when a client closes the connection to the server
type OnClientDisconnected func(client *Client)

// OnSignal is a required hook.
// It's invoked when the webwire server receives a signal from the client
type OnSignal func(ctx context.Context)

// OnRequest is an optional hook.
// It's invoked when the webwire server receives a request from the client.
// It must return either a response payload or an error
type OnRequest func(ctx context.Context) (response []byte, err *Error)

// OnSessionCreated is a required hook for sessions to be supported.
// It's invoked right after the creation of a new session.
// The WebWire server isn't responsible for storing the sessions it creates,
// the user must save the given session in this hook either to a database,
// a filesystem or any other kind of persistent or volatile storage
// for onSessionLookup to later be able to restore it by the session key
type OnSessionCreated func(client *Client) (error)

// OnSessionLookup is a required hook for sessions to be supported.
// It's invoked when the server is looking for a specific session given its key.
// The user is responsible for returning the exact copy of the session object
// associated with the given key for sessions to be restorable
type OnSessionLookup func(key string) (*Session, error)

// OnSessionClosed is a required hook for sessions to be supported.
// It's invoked when the active session of the given client
// is closed (thus destroyed) either by the server or the client himself.
// The user is responsible for removing the current session of the given client
// from its storage for the session to be actually and properly destroyed.
type OnSessionClosed func(client *Client) (error)

// Server represents the actual 
type Server struct {
	// Configuration
	onClientConnected OnClientConnected
	onClientDisconnected OnClientDisconnected
	onSignal OnSignal
	onRequest OnRequest
	OnSessionCreated OnSessionCreated
	onSessionLookup OnSessionLookup
	onSessionClosed OnSessionClosed
	onOptions OnOptions

	// Dynamic methods
	launch func() error
	
	// State
	Addr string
	clientsLock *sync.Mutex
	clients []*Client
	sessionsEnabled bool
	activeSessions map[string]*Client

	// Internals
	httpServer *http.Server
	upgrader websocket.Upgrader
	warnLog *log.Logger
	errorLog *log.Logger
}

// NewServer creates a new WebWire server instance.
func NewServer(
	addr string,
	onClientConnected OnClientConnected,
	onClientDisconnected OnClientDisconnected,
	onSignal OnSignal,
	onRequest OnRequest,
	OnSessionCreated OnSessionCreated,
	onSessionLookup OnSessionLookup,
	onSessionClosed OnSessionClosed,
	onOptions OnOptions,
	warningLogWriter io.Writer,
	errorLogWriter io.Writer,
) (*Server, error) {
	if onClientConnected == nil {
		onClientConnected = func(_ *Client) {}
	}

	if onClientDisconnected == nil {
		onClientDisconnected = func(_ *Client) {}
	}

	if onSignal == nil {
		onSignal = func(_ context.Context) {}
	}

	if onRequest == nil {
		onRequest = func(_ context.Context) (response []byte, err *Error) {
			return nil, &Error {
				"NOT_IMPLEMENTED",
				fmt.Sprintf("Request handling is not implemented " +
					" on this server instance",
				),
			}
		}
	}

	if onOptions == nil {
		onOptions = func(resp http.ResponseWriter)	{}
	}

	sessionsEnabled := false
	if OnSessionCreated != nil &&
		onSessionLookup != nil &&
		onSessionClosed != nil {
		sessionsEnabled = true
	}

	srv := Server {
		// Configuration
		onClientConnected: onClientConnected,
		onClientDisconnected: onClientDisconnected,
		onSignal: onSignal,
		onRequest: onRequest,
		OnSessionCreated: OnSessionCreated,
		onSessionLookup: onSessionLookup,
		onSessionClosed: onSessionClosed,
		onOptions: onOptions,

		// State
		clients: make([]*Client, 0),
		clientsLock: &sync.Mutex {},
		sessionsEnabled: sessionsEnabled,
		activeSessions: make(map[string]*Client),

		// Internals
		warnLog: log.New(
			warningLogWriter,
			"WARNING: ",
			log.Ldate | log.Ltime | log.Lshortfile,
		),
		errorLog: log.New(
			errorLogWriter,
			"ERROR: ",
			log.Ldate | log.Ltime | log.Lshortfile,
		),
	}

	// Initialize websocket
	srv.upgrader = websocket.Upgrader {
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	// Initialize HTTP server
	srv.httpServer = &http.Server {
		Addr: addr,
		Handler: &srv,
	}

	// Determine final address
	addr = srv.httpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	// Initialize TCP/IP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up TCP/IP listener: %s", err)
	}

	srv.launch = func() error {
		// Launch server
		err = srv.httpServer.Serve(
			tcpKeepAliveListener{listener.(*net.TCPListener)},
		)
		if err != nil {
			return fmt.Errorf("HTTP Server failure: %s", err)
		}
		return nil
	}

	// Remember HTTP server address
	srv.Addr = listener.Addr().String()

	return &srv, nil
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// handleSessionRestore handles session restoration (by session key) requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionRestore(msg *Message) (error) {
	key := string(msg.Payload)
	if _, exists := srv.activeSessions[key]; exists {
		msg.fail(Error {
			"SESSION_ACTIVE",
			fmt.Sprintf(
				"The session identified by key: '%s' is already active",
				key,
			),
		})
		return nil
	}

	session, err := srv.onSessionLookup(key)
	if err != nil {
		msg.fail(Error {
			"INTERNAL_ERROR",
			fmt.Sprintf("Session restoration request not could have been fulfilled"),
		})
		return fmt.Errorf(
			"CRITICAL: Session search handler failed: %s", err,
		)
	}
	if session == nil {
		msg.fail(Error {
			"SESSION_NOT_FOUND",
			fmt.Sprintf("No session associated with key: '%s'", key),
		})
		return nil
	}
	msg.Client.Session = session
	srv.activeSessions[session.Key] = msg.Client

	// Send confirmation response
	msg.fulfill(nil)

	return nil
}

// handleSessionClosure handles session destruction requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionClosure(msg *Message) error {
	if msg.Client.Session == nil {
		// Send confirmation even though no session was closed
		msg.fulfill(nil)
		return nil
	}

	if err := srv.deregisterSession(msg.Client); err != nil {
		msg.fail(Error {
			"INTERNAL_ERROR",
			fmt.Sprintf("Session destruction request not could have been fulfilled"),
		})
		return fmt.Errorf("CRITICAL: Internal server error, session destruction failed: %s", err)
	}

	// Synchronize session destruction to the client
	if err := msg.Client.notifySessionClosed(); err != nil {
		msg.fail(Error {
			"INTERNAL_ERROR",
			fmt.Sprintf("Session destruction request not could have been fulfilled"),
		})
		return fmt.Errorf("CRITICAL: Internal server error," +
			" couldn't notify client about the session destruction: %s",
			err,
		)
	}

	// Send confirmation
	msg.fulfill(nil)

	return nil
}

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSignal(msg *Message) error {
	srv.onSignal(context.WithValue(context.Background(), MESSAGE, *msg))
	return nil
}

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleRequest(msg *Message) {
	response, err := srv.onRequest(
		context.WithValue(context.Background(), MESSAGE, *msg),
	)
	if err != nil {
		msg.fail(*err)
	}
	msg.fulfill(response)
}

// handleReply handles incoming responses to requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleReply(msg *Message) error {
	// TODO: implement server-side response handling
	return nil
}

// handleErrorResponse handles incoming responses to failed requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleErrorResponse(msg *Message) error {
	// TODO: implement server-side error-response handling
	return nil
}

// handleMetadata handles endpoint metadata requests
func (srv *Server) handleMetadata(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(struct {
		ProtocolVersion string `json:"protocol-version"`
	} {
		protocolVersion,
	})
}

// ServeHTTP will make the server listen for incoming HTTP requests
// eventually trying to upgrade them to WebSocket connections
func (srv Server) ServeHTTP(
	resp http.ResponseWriter,
	req *http.Request,
) {
	switch(req.Method) {
	case "OPTIONS":
		srv.onOptions(resp)
		return
	case "WEBWIRE":
		srv.handleMetadata(resp)
		return
	}

	// Establish connection
	conn, err := srv.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		srv.errorLog.Print("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Register connected client
	newClient := &Client {
		&srv,
		&sync.Mutex {},
		conn,
		time.Now(),
		nil,
	}

	srv.clientsLock.Lock()
	srv.clients = append(srv.clients, newClient)
	srv.clientsLock.Unlock()

	// Call hook on successful connection
	srv.onClientConnected(newClient)

	for {
		// Await message
		wsMsgType, message, err := conn.ReadMessage()
		if err != nil {
			isClosed := websocket.IsCloseError(err)
			isUnexpectedlyClosed := websocket.IsUnexpectedCloseError(err)
			if newClient.Session != nil && (isClosed || isUnexpectedlyClosed) {
				// Mark session as inactive
				delete(srv.activeSessions, newClient.Session.Key)
			}
			if isClosed {
				srv.onClientDisconnected(newClient)
				break
			} else if isUnexpectedlyClosed {
				srv.onClientDisconnected(newClient)
				break
			}
			srv.warnLog.Println("Reading failed:", err)
			break
		}

		// Parse message
		msg, err := ParseMessage(message)
		if err != nil {
			srv.errorLog.Println("Failed parsing message:", err)
			break
		}

		// Prepare message
		// Reference the client associated with this message
		msg.Client = newClient

		// Create fulfillment lambda
		msg.fulfill = func(response []byte) {
			// Send response
			header := append([]byte("p"), *msg.id...)
			err := newClient.write(
				wsMsgType, append(header, response...),
			)
			if err != nil {
				srv.errorLog.Println("Writing failed:", err)
			}
		}

		// Create failure lambda
		msg.fail = func(errObj Error) {
			encoded, err := json.Marshal(errObj)
			if err != nil {
				encoded = []byte("CRITICAL: could not encode error report")
			}

			// Send request failure notification
			header := append([]byte("e"), *msg.id...)
			err = newClient.write(
				websocket.TextMessage,
				append(header, encoded...),
			)
			if err != nil {
				srv.errorLog.Println("Writing failed:", err)
			}
		}

		// Handle message
		switch msg.msgType {
		case MsgSignal: err = srv.handleSignal(&msg)
		case MsgRequest: srv.handleRequest(&msg)
		case MsgReply: err = srv.handleReply(&msg)
		case MsgErrorReply: err = srv.handleErrorResponse(&msg)
		case MsgRestoreSession: err = srv.handleSessionRestore(&msg)
		case MsgCloseSession: err = srv.handleSessionClosure(&msg)
		}

		if err != nil {
			srv.errorLog.Printf("CRITICAL FAILURE: %s", err)
			break
		}
	}
}

func (srv *Server) registerSession(clt *Client, session *Session) error {
	clt.Session = session
	srv.activeSessions[session.Key] = clt
	err := srv.OnSessionCreated(clt)
	if err != nil {
		return fmt.Errorf("Couldn't save session: %s", err)
	}
	return nil
}

func (srv *Server) deregisterSession(clt *Client) error {
	delete(srv.activeSessions, clt.Session.Key)
	err := srv.onSessionClosed(clt)
	if err != nil {
		return fmt.Errorf("CRITICAL: Session closure handler failed: %s", err)
	}
	return nil
}

// Run will launch the server blocking the calling goroutine
func (srv *Server) Run() error {
	return srv.launch()
}
