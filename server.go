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

type OnCORS func(resp http.ResponseWriter)

// OnClientConnected is an optional callback.
// It's invoked when a new client successfuly connects to the webwire server
type OnClientConnected func(client *Client)

// OnSignal is an optional callback.
// It's invoked when the webwire server receives a signal from the client
type OnSignal func(ctx context.Context)

// OnRequest is an optional callback.
// It's invoked when the webwire server receives a request from the client.
// It must return either a response payload or an error.
type OnRequest func(ctx context.Context) (response []byte, err *Error)

// OnSaveSession is a required callback invoked after the creation of a new session.
// Because webwire isn't responsible for storing the sessions the library users must
// provide this callback to persist the session to whatever storage they like
type OnSaveSession func(session *Session) (error)

// OnFindSession is a required callback.
// It's invoked when the webwire server is looking for a specific session.
// Because webwire isn't responsible for storing the sessions the library users must
// provide this callback to find a persisted session by the given key
type OnFindSession func(key string) (*Session, error)

// OnSessionClosure is a required callback.
// It's invoked when the webwire server is closing a specific session.
// Because webwire isn't responsible for storing the sessions the library users must
// provide this callback to close and delete the session associated with the given key
type OnSessionClosure func(key string) (error)

type Server struct {
	// Configuration
	onClientConnected OnClientConnected
	onSignal OnSignal
	onRequest OnRequest
	onSaveSession OnSaveSession
	onFindSession OnFindSession
	onSessionClosure OnSessionClosure
	onCORS OnCORS

	// Dynamic methods
	launch func() error
	
	// State
	Addr string
	clientsLock *sync.Mutex
	clients []*Client
	sessionsEnabled bool
	activeSessions map[string]bool

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
	onSignal OnSignal,
	onRequest OnRequest,
	onSaveSession OnSaveSession,
	onFindSession OnFindSession,
	onSessionClosure OnSessionClosure,
	onCORS OnCORS,
	warningLogWriter io.Writer,
	errorLogWriter io.Writer,
) (*Server, error) {
	if onClientConnected == nil {
		onClientConnected = func(_ *Client) {}
	}

	if onSignal == nil {
		onSignal = func(_ context.Context) {}
	}

	if onRequest == nil {
		onRequest = func(_ context.Context) (response []byte, err *Error) {
			return nil, &Error {
				"NOT_IMPLEMENTED",
				fmt.Sprintf("Request handling is not implemented on this server instance"),
			}
		}
	}

	if onCORS == nil {
		onCORS = func(resp http.ResponseWriter)	{}
	}

	sessionsEnabled := false
	if onSaveSession != nil && onFindSession != nil && onSessionClosure != nil {
		sessionsEnabled = true
	}

	srv := Server {
		// Configuration
		onClientConnected: onClientConnected,
		onSignal: onSignal,
		onRequest: onRequest,
		onSaveSession: onSaveSession,
		onFindSession: onFindSession,
		onSessionClosure: onSessionClosure,
		onCORS: onCORS,

		// State
		clients: make([]*Client, 0),
		clientsLock: &sync.Mutex {},
		sessionsEnabled: sessionsEnabled,
		activeSessions: make(map[string]bool),

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
		err = srv.httpServer.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
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
func (srv *Server) handleSessionRestore(
	msg *Message,
	currentSession **Session,
) (error) {
	key := string(msg.Payload)
	if _, exists := srv.activeSessions[key]; exists {
		msg.fail(Error {
			"SESS_ACTIVE",
			fmt.Sprintf(
				"The session identified by key: '%s' is already active",
				key,
			),
		})
		return nil
	}

	session, err := srv.onFindSession(key)
	if err != nil {
		return fmt.Errorf(
			"CRITICAL: Session search handler failed: %s", err,
		)
	}
	if session == nil {
		msg.fail(Error {
			"SESS_NOT_FOUND",
			fmt.Sprintf("No session associated with key: '%s'", key),
		})
		return nil
	}
	(*currentSession) = session
	srv.activeSessions[session.Key] = true

	// Send confirmation response
	msg.fulfill(nil)

	return nil
}

// handleSessionClosure handles session closure signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionClosure(
	msg *Message,
	currentSession **Session,
) error {
	if *currentSession == nil {
		// Send confirmation even though no session was closed
		msg.fulfill(nil)
		return nil
	}

	key := (*currentSession).Key
	err := srv.onSessionClosure(key)
	if err != nil {
		return fmt.Errorf("CRITICAL: Session closure handler failed: %s", err)
	}
	*currentSession = nil
	delete(srv.activeSessions, key)

	// Send confirmation
	msg.fulfill(nil)

	return nil
}

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSignal(
	msg *Message,
	currentSession *Session,
) error {
	srv.onSignal(context.WithValue(context.Background(), MESSAGE, *msg))
	return nil
}

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleRequest(
	msg *Message,
	currentSession *Session,
) {
	response, err := srv.onRequest(
		context.WithValue(context.Background(), MESSAGE, *msg),
	)
	if err != nil {
		msg.fail(*err)
	}
	msg.fulfill(response)
}

// handleResponse handles incoming responses to requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleResponse(
	msg *Message,
	currentSession *Session,
) error {
	// TODO: implement server-side response handling
	return nil
}

// handleErrorResponse handles incoming responses to failed requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleErrorResponse(
	msg *Message,
	currentSession *Session,
) error {
	// TODO: implement server-side error-response handling
	return nil
}

func (srv *Server) CORS(resp http.ResponseWriter) {
	srv.onCORS(resp)
}

func (srv Server) ServeHTTP(
	resp http.ResponseWriter,
	req *http.Request,
) {
	if req.Method == "OPTIONS" {
		srv.CORS(resp)
		return
	}

	var session *Session

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
			if session != nil && (
				websocket.IsCloseError(err) ||
				websocket.IsUnexpectedCloseError(err)) {
				// Mark session as inactive
				delete(srv.activeSessions, session.Key)
				break
			} else if websocket.IsCloseError(err) {
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
		msg.Client = newClient
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

		switch msg.msgType {
		case MsgTyp_SIGNAL: err = srv.handleSignal(&msg, session)
		case MsgTyp_REQUEST: srv.handleRequest(&msg, session)
		case MsgTyp_RESPONSE: err = srv.handleResponse(&msg, session)
		case MsgTyp_ERROR_RESP: err = srv.handleErrorResponse(&msg, session)
		case MsgTyp_SESS_RESTORE: err = srv.handleSessionRestore(&msg, &session)
		case MsgTyp_SESS_CLOSURE: err = srv.handleSessionClosure(&msg, &session)
		}

		if err != nil {
			srv.errorLog.Printf("CRITICAL FAILURE: %s", err)
			break
		}
	}
}

func (srv *Server) registerSession(session *Session) error {
	err := srv.onSaveSession(session)
	if err != nil {
		return fmt.Errorf("Couldn't save session: %s", err)
	}
	srv.activeSessions[session.Key] = true
	return nil
}

func (srv *Server) Run() error {
	return srv.launch()
}

func (srv *Server) ClientsNum() int {
	srv.clientsLock.Lock()
	defer srv.clientsLock.Unlock()
	ln := len(srv.clients)
	return ln
}
