package webwire

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const protocolVersion = "1.0"

// Hooks represents all callback hook functions
type Hooks struct {
	// OnOptions is an optional hook.
	// It's invoked when the websocket endpoint is examined by the client
	// using the HTTP OPTION method.
	OnOptions func(resp http.ResponseWriter)

	// OnClientConnected is an optional hook.
	// It's invoked when a new client establishes a connection to the server
	OnClientConnected func(client *Client)

	// OnClientDisconnected is an optional hook.
	// It's invoked when a client closes the connection to the server
	OnClientDisconnected func(client *Client)

	// OnSignal is a required hook.
	// It's invoked when the webwire server receives a signal from the client
	OnSignal func(ctx context.Context)

	// OnRequest is an optional hook.
	// It's invoked when the webwire server receives a request from the client.
	// It must return either a response payload or an error
	OnRequest func(ctx context.Context) (response []byte, err *Error)

	// OnSessionCreated is a required hook for sessions to be supported.
	// It's invoked right after the synchronisation of the new session to the remote client.
	// The WebWire server isn't responsible for permanently storing the sessions it creates,
	// it's up to the user to save the given session in this hook either to a database,
	// a filesystem or any other kind of persistent or volatile storage
	// for OnSessionLookup to later be able to restore it by the session key.
	// If OnSessionCreated fails returning an error then the failure is logged
	// but the session isn't destroyed and remains active.
	// The only consequence of OnSessionCreation failing is that the server won't be able
	// to restore the session after the client is disconnected
	OnSessionCreated func(client *Client) error

	// OnSessionLookup is a required hook for sessions to be supported.
	// It's invoked when the server is looking for a specific session given its key.
	// The user is responsible for returning the exact copy of the session object
	// associated with the given key for sessions to be restorable.
	// If OnSessionLookup fails returning an error then the failure is logged
	OnSessionLookup func(key string) (*Session, error)

	// OnSessionClosed is a required hook for sessions to be supported.
	// It's invoked when the active session of the given client
	// is closed (thus destroyed) either by the server or the client himself.
	// The user is responsible for removing the current session of the given client
	// from its storage for the session to be actually and properly destroyed.
	// If OnSessionClosed fails returning an error then the failure is logged
	OnSessionClosed func(client *Client) error
}

// SetDefaults sets undefined required hooks
func (hooks *Hooks) SetDefaults() {
	if hooks.OnClientConnected == nil {
		hooks.OnClientConnected = func(_ *Client) {}
	}

	if hooks.OnClientDisconnected == nil {
		hooks.OnClientDisconnected = func(_ *Client) {}
	}

	if hooks.OnSignal == nil {
		hooks.OnSignal = func(_ context.Context) {}
	}

	if hooks.OnRequest == nil {
		hooks.OnRequest = func(_ context.Context) (response []byte, err *Error) {
			return nil, &Error{
				"NOT_IMPLEMENTED",
				fmt.Sprintf("Request handling is not implemented " +
					" on this server instance",
				),
			}
		}
	}

	if hooks.OnOptions == nil {
		hooks.OnOptions = func(resp http.ResponseWriter) {}
	}
}

// Server represents the actual
type Server struct {
	hooks Hooks

	// Dynamic methods
	launch func() error

	// State
	Addr            string
	clientsLock     *sync.Mutex
	clients         []*Client
	sessionsEnabled bool
	SessionRegistry sessionRegistry

	// Internals
	httpServer *http.Server
	upgrader   websocket.Upgrader
	warnLog    *log.Logger
	errorLog   *log.Logger
}

// NewServer creates a new WebWire server instance.
func NewServer(
	addr string,
	hooks Hooks,
	warningLogWriter io.Writer,
	errorLogWriter io.Writer,
) (*Server, error) {
	hooks.SetDefaults()

	sessionsEnabled := false
	if hooks.OnSessionCreated != nil &&
		hooks.OnSessionLookup != nil &&
		hooks.OnSessionClosed != nil {
		sessionsEnabled = true
	}

	srv := Server{
		hooks: hooks,

		// State
		clients:         make([]*Client, 0),
		clientsLock:     &sync.Mutex{},
		sessionsEnabled: sessionsEnabled,
		SessionRegistry: newSessionRegistry(),

		// Internals
		warnLog: log.New(
			warningLogWriter,
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		errorLog: log.New(
			errorLogWriter,
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		),
	}

	// Initialize websocket
	srv.upgrader = websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	// Initialize HTTP server
	srv.httpServer = &http.Server{
		Addr:    addr,
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
func (srv *Server) handleSessionRestore(msg *Message) error {
	if !srv.sessionsEnabled {
		msg.fail(Error{
			"SESSIONS_DISABLED",
			"Sessions are disabled on this server instance",
		})
		return nil
	}

	key := string(msg.Payload)

	sessionExists := srv.SessionRegistry.Exists(key)

	if sessionExists {
		msg.fail(Error{
			"SESSION_ACTIVE",
			fmt.Sprintf(
				"The session identified by key: '%s' is already active",
				key,
			),
		})
		return nil
	}

	session, err := srv.hooks.OnSessionLookup(key)
	if err != nil {
		msg.fail(Error{
			"INTERNAL_ERROR",
			fmt.Sprintf(
				"Session restoration request not could have been fulfilled",
			),
		})
		return fmt.Errorf(
			"CRITICAL: Session search handler failed: %s", err,
		)
	}
	if session == nil {
		msg.fail(Error{
			"SESSION_NOT_FOUND",
			fmt.Sprintf("No session associated with key: '%s'", key),
		})
		return nil
	}

	// JSON encode the session
	encodedSession, err := json.Marshal(session)
	if err != nil {
		msg.fail(Error{
			"INTERNAL_ERROR",
			fmt.Sprintf(
				"Session restoration request not could have been fulfilled",
			),
		})
		return fmt.Errorf("Couldn't encode session object (%v): %s", session, err)
	}

	msg.Client.Session = session
	srv.SessionRegistry.register(msg.Client)

	msg.fulfill(encodedSession)

	return nil
}

// handleSessionClosure handles session destruction requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionClosure(msg *Message) error {
	if !srv.sessionsEnabled {
		msg.fail(Error{
			"SESSIONS_DISABLED",
			"Sessions are disabled on this server instance",
		})
		return nil
	}

	if msg.Client.Session == nil {
		// Send confirmation even though no session was closed
		msg.fulfill(nil)
		return nil
	}

	srv.deregisterSession(msg.Client)

	// Synchronize session destruction to the client
	if err := msg.Client.notifySessionClosed(); err != nil {
		msg.fail(Error{
			"INTERNAL_ERROR",
			fmt.Sprintf(
				"Session destruction request not could have been fulfilled",
			),
		})
		return fmt.Errorf("CRITICAL: Internal server error, "+
			"couldn't notify client about the session destruction: %s",
			err,
		)
	}

	// Reset the session on the client agent
	msg.Client.Session = nil

	// Send confirmation
	msg.fulfill(nil)

	return nil
}

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSignal(msg *Message) {
	srv.hooks.OnSignal(context.WithValue(context.Background(), MESSAGE, *msg))
}

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleRequest(msg *Message) {
	response, err := srv.hooks.OnRequest(
		context.WithValue(context.Background(), MESSAGE, *msg),
	)
	if err != nil {
		msg.fail(*err)
	}
	msg.fulfill(response)
}

// handleMetadata handles endpoint metadata requests
func (srv *Server) handleMetadata(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(struct {
		ProtocolVersion string `json:"protocol-version"`
	}{
		protocolVersion,
	})
}

// handleMessage handles incoming messages
func (srv *Server) handleMessage(msg *Message) error {
	switch msg.msgType {
	case MsgSignal:
		srv.handleSignal(msg)
	case MsgRequest:
		srv.handleRequest(msg)
	case MsgRestoreSession:
		return srv.handleSessionRestore(msg)
	case MsgCloseSession:
		return srv.handleSessionClosure(msg)
	}
	return nil
}

// ServeHTTP will make the server listen for incoming HTTP requests
// eventually trying to upgrade them to WebSocket connections
func (srv *Server) ServeHTTP(
	resp http.ResponseWriter,
	req *http.Request,
) {
	switch req.Method {
	case "OPTIONS":
		srv.hooks.OnOptions(resp)
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
	newClient := &Client{
		srv,
		&sync.Mutex{},
		conn,
		time.Now(),
		nil,
	}

	srv.clientsLock.Lock()
	srv.clients = append(srv.clients, newClient)
	srv.clientsLock.Unlock()

	// Call hook on successful connection
	srv.hooks.OnClientConnected(newClient)

	for {
		// Await message
		_, message, err := conn.ReadMessage()
		if err != nil {
			if newClient.Session != nil {
				// Mark session as inactive
				srv.SessionRegistry.deregister(newClient)
			}

			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				srv.warnLog.Printf("Reading failed: %s", err)
			}

			srv.hooks.OnClientDisconnected(newClient)
			return
		}

		// Parse message
		var msg Message
		if err := msg.Parse(message); err != nil {
			srv.errorLog.Println("Failed parsing message:", err)
			break
		}

		// Prepare message
		// Reference the client associated with this message
		msg.Client = newClient

		msg.createReplyCallback(newClient, srv)
		msg.createFailCallback(newClient, srv)

		// Handle message
		if err := srv.handleMessage(&msg); err != nil {
			srv.errorLog.Printf("CRITICAL FAILURE: %s", err)
			break
		}
	}
}

func (srv *Server) registerSession(clt *Client) {
	srv.SessionRegistry.register(clt)
	// Execute session creation hook
	if err := srv.hooks.OnSessionCreated(clt); err != nil {
		srv.errorLog.Printf("OnSessionCreated hook failed: %s", err)
	}
}

func (srv *Server) deregisterSession(clt *Client) {
	srv.SessionRegistry.deregister(clt)
	if err := srv.hooks.OnSessionClosed(clt); err != nil {
		srv.errorLog.Printf("OnSessionClosed hook failed: %s", err)
	}
}

// Run will launch the server blocking the calling goroutine
func (srv *Server) Run() error {
	return srv.launch()
}
