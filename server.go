package webwire

import (
	"io"
	"log"
	"fmt"
	"context"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type OnCORS func(resp http.ResponseWriter)

// OnSignal is an optional callback.
// It's invoked when the webwire server receives a signal from the client
type OnSignal func(data []byte, session *Session)

// OnRequest is an optional callback.
// It's invoked when the webwire server receives a request from the client.
// It must return either a response payload or an error.
type OnRequest func(ctx context.Context) (response []byte, err *Error)

// OnSessionCreation is an optional callback.
// It's invoked when a connected client attempts to authenticate itself by the given credentials.
// The user code must return whether or not the authentication was successful.
// Arbitrary data also can also be attached to the session
// by returning anything but nil in the second return value.
type OnSessionCreation func(login, password string) (bool, interface {}, *Error)

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
	onSignal OnSignal
	onRequest OnRequest
	onSessionCreation OnSessionCreation
	onSaveSession OnSaveSession
	onFindSession OnFindSession
	onSessionClosure OnSessionClosure
	onCORS OnCORS
	
	// Internal state
	sessionsEnabled bool
	activeSessions map[string]bool

	upgrader websocket.Upgrader
	warnLog *log.Logger
	errorLog *log.Logger
}

func NewServer(
	onSignal OnSignal,
	onRequest OnRequest,
	onSessionCreation OnSessionCreation,
	onSaveSession OnSaveSession,
	onFindSession OnFindSession,
	onSessionClosure OnSessionClosure,
	onCORS OnCORS,
	warningLogWriter io.Writer,
	errorLogWriter io.Writer,
) Server {
	if onSignal == nil {
		onSignal = func(data []byte, session *Session) {}
	}

	if onRequest == nil {
		onRequest = func(_ context.Context) (response []byte, err *Error) {
			return nil, &Error {
				"NOT_IMPLEMENTED",
				fmt.Sprintf("Request handling is not implemented on this server instance"),
			}
		}
	}

	if onSessionCreation == nil {
		onSessionCreation = func(login, password string) (bool, interface {}, *Error) {
			return false, nil, &Error {
				"NOT_IMPLEMENTED",
				fmt.Sprintf("Session creation is not implemented on this server instance"),
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

	return Server {
		onSignal,
		onRequest,
		onSessionCreation,
		onSaveSession,
		onFindSession,
		onSessionClosure,
		onCORS,
		sessionsEnabled,
		make(map[string]bool),
		websocket.Upgrader {
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
		log.New(
			warningLogWriter,
			"WARNING: ",
			log.Ldate | log.Ltime | log.Lshortfile,
		),
		log.New(
			errorLogWriter,
			"ERROR: ",
			log.Ldate | log.Ltime | log.Lshortfile,
		),
	}
}

// handleSessionRestore handles session restoration (by session key) requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionRestore(
	msg *Message,
	currentSession **Session,
) (error) {
	key := string(msg.payload)
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
	srv.activeSessions[session.key] = true

	// Send confirmation response
	msg.fulfill(nil)

	return nil
}

// handleSessionCreation handles session creation requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionCreation(
	msg *Message,
	currentSession **Session,
) error {
	if (*currentSession) != nil {
		msg.fulfill([]byte((*currentSession).key))
		return nil
	}

	if !srv.sessionsEnabled {
		msg.fail(Error {
			"NOT_IMPLEMENTED",
			"Sessions are not implemented on this server instance",
		})
		return nil
	}

	// Decode credentials
	var credentials struct {
		Login string `json:"l"`
		Password string `json:"p"`
	}
	if err := json.Unmarshal(msg.payload, &credentials); err != nil {
		return fmt.Errorf("CRITICAL: Failed unmarshalling credentials: %s", err)
	}

	// Check client authentication
	success, sessionInfo, sessCrtErr := srv.onSessionCreation(
		credentials.Login,
		credentials.Password,
	)
	if sessCrtErr != nil {
		msg.fail(*sessCrtErr)
		return nil
	}
	if !success {
		msg.fail(Error {
			"INVALID_CRED",
			"Authentication failed (invalid credentials)",
		})
		return nil
	}

	// Save session
	session := NewSession(
		UNKNOWN,
		// TODO: set proper user agent string
		"user agent",
		sessionInfo,
	)
	err := srv.onSaveSession(&session)
	if err != nil {
		return fmt.Errorf("CRITICAL: Session saving handler failed: %s", err)
	}
	(*currentSession) = &session
	srv.activeSessions[session.key] = true

	// Send confirmation response including the session key
	msg.fulfill([]byte(session.key))

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

	key := (*currentSession).key
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
	srv.onSignal(msg.payload, currentSession)
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

	for {
		// Await message
		wsMsgType, message, err := conn.ReadMessage()
		if err != nil {
			if session != nil && (
				websocket.IsCloseError(err) ||
				websocket.IsUnexpectedCloseError(err)) {
				// Mark session as inactive
				delete(srv.activeSessions, session.key)
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
		msg.session = session
		msg.fulfill = func(response []byte) {
			// Send response
			header := append([]byte("p"), *msg.id...)
			err := conn.WriteMessage(
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
			err = conn.WriteMessage(
				websocket.TextMessage,
				append(header, encoded...),
			)
			if err != nil {
				srv.errorLog.Println("Writing failed:", err)
			}
		}

		switch msg.msgType {
		case SIGNAL: err = srv.handleSignal(&msg, session)
		case REQUEST: srv.handleRequest(&msg, session)
		case RESPONSE: err = srv.handleResponse(&msg, session)
		case ERROR_RESP: err = srv.handleErrorResponse(&msg, session)
		case SESS_CREATION: err = srv.handleSessionCreation(&msg, &session)
		case SESS_RESTORE: err = srv.handleSessionRestore(&msg, &session)
		case SESS_CLOSURE: err = srv.handleSessionClosure(&msg, &session)
		}

		if err != nil {
			srv.errorLog.Printf("CRITICAL FAILURE: %s", err)
			break
		}
	}
}