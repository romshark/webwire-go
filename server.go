package webwire

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const protocolVersion = "1.3"

// Server represents a headless WebWire server instance,
// where headless means there's no HTTP server that's hosting it
type Server struct {
	impl              ServerImplementation
	sessionManager    SessionManager
	sessionKeyGen     SessionKeyGenerator
	sessionInfoParser SessionInfoParser

	// State
	shutdown        bool
	shutdownRdy     chan bool
	currentOps      uint32
	opsLock         sync.Mutex
	clientsLock     *sync.Mutex
	clients         []*Client
	sessionsEnabled bool
	sessionRegistry *sessionRegistry

	// Internals
	connUpgrader ConnUpgrader
	warnLog      *log.Logger
	errorLog     *log.Logger
}

// NewServer creates a new WebWire server instance
func NewServer(implementation ServerImplementation, opts ServerOptions) *Server {
	if implementation == nil {
		panic(fmt.Errorf("A headed webwire server requires a server implementation, got nil"))
	}

	opts.SetDefaults()

	srv := Server{
		impl:              implementation,
		sessionManager:    opts.SessionManager,
		sessionKeyGen:     opts.SessionKeyGenerator,
		sessionInfoParser: opts.SessionInfoParser,

		// State
		shutdown:        false,
		shutdownRdy:     make(chan bool),
		currentOps:      0,
		opsLock:         sync.Mutex{},
		clients:         make([]*Client, 0),
		clientsLock:     &sync.Mutex{},
		sessionsEnabled: opts.SessionsEnabled,
		sessionRegistry: newSessionRegistry(opts.MaxSessionConnections),

		// Internals
		connUpgrader: newConnUpgrader(),
		warnLog:      opts.WarnLog,
		errorLog:     opts.ErrorLog,
	}

	return &srv
}

// handleSessionRestore handles session restoration (by session key) requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionRestore(clt *Client, msg *Message) error {
	if !srv.sessionsEnabled {
		srv.failMsg(clt, msg, SessionsDisabledErr{})
		return nil
	}

	key := string(msg.Payload.Data)

	if srv.sessionRegistry.maxConns > 0 &&
		srv.sessionRegistry.SessionConnections(key)+1 > srv.sessionRegistry.maxConns {
		srv.failMsg(clt, msg, MaxSessConnsReachedErr{})
		return nil
	}

	//session, err := srv.sessionManager.OnSessionLookup(key)
	exists, creation, info, err := srv.sessionManager.OnSessionLookup(key)
	if err != nil {
		srv.failMsg(clt, msg, nil)
		return fmt.Errorf("CRITICAL: Session search handler failed: %s", err)
	}
	if !exists {
		srv.failMsg(clt, msg, SessNotFoundErr{})
		return nil
	}

	encodedSessionObj := JSONEncodedSession{
		Key:      key,
		Creation: creation,
		Info:     info,
	}

	// JSON encode the session
	encodedSession, err := json.Marshal(&encodedSessionObj)
	if err != nil {
		// TODO: return internal server error and log it
		srv.failMsg(clt, msg, nil)
		return fmt.Errorf(
			"Couldn't encode session object (%v): %s",
			encodedSessionObj,
			err,
		)
	}

	// parse attached session info
	var parsedSessInfo SessionInfo
	if info != nil && srv.sessionInfoParser != nil {
		parsedSessInfo = srv.sessionInfoParser(info)
	}

	clt.setSession(&Session{
		Key:      key,
		Creation: creation,
		Info:     parsedSessInfo,
	})
	if okay := srv.sessionRegistry.register(clt); !okay {
		panic(fmt.Errorf("The number of concurrent session connections was " +
			"unexpectedly exceeded",
		))
	}

	srv.fulfillMsg(clt, msg, Payload{
		Encoding: EncodingUtf8,
		Data:     encodedSession,
	})

	return nil
}

// handleSessionClosure handles session destruction requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSessionClosure(clt *Client, msg *Message) error {
	if !srv.sessionsEnabled {
		srv.failMsg(clt, msg, SessionsDisabledErr{})
		return nil
	}

	if !clt.HasSession() {
		// Send confirmation even though no session was closed
		srv.fulfillMsg(clt, msg, Payload{})
		return nil
	}

	srv.deregisterSession(clt)

	// Synchronize session destruction to the client
	if err := clt.notifySessionClosed(); err != nil {
		srv.failMsg(clt, msg, nil)
		return fmt.Errorf("CRITICAL: Internal server error, "+
			"couldn't notify client about the session destruction: %s",
			err,
		)
	}

	// Reset the session on the client agent
	clt.setSession(nil)

	// Send confirmation
	srv.fulfillMsg(clt, msg, Payload{})

	return nil
}

// handleSignal handles incoming signals
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleSignal(clt *Client, msg *Message) {
	srv.opsLock.Lock()
	// Ignore incoming signals during shutdown
	if srv.shutdown {
		srv.opsLock.Unlock()
		return
	}
	srv.currentOps++
	srv.opsLock.Unlock()

	srv.impl.OnSignal(
		context.Background(),
		clt,
		msg,
	)

	// Mark signal as done and shutdown the server if scheduled and no ops are left
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()
}

// handleRequest handles incoming requests
// and returns an error if the ongoing connection cannot be proceeded
func (srv *Server) handleRequest(clt *Client, msg *Message) {
	srv.opsLock.Lock()
	// Reject incoming requests during shutdown, return special shutdown error
	if srv.shutdown {
		srv.opsLock.Unlock()
		srv.failMsgShutdown(clt, msg)
		return
	}
	srv.currentOps++
	srv.opsLock.Unlock()

	replyPayload, returnedErr := srv.impl.OnRequest(
		context.Background(),
		clt,
		msg,
	)
	switch returnedErr.(type) {
	case nil:
		srv.fulfillMsg(clt, msg, replyPayload)
	case ReqErr:
		srv.failMsg(clt, msg, returnedErr)
	case *ReqErr:
		srv.failMsg(clt, msg, returnedErr)
	default:
		srv.errorLog.Printf("Internal error during request handling: %s", returnedErr)
		srv.failMsg(clt, msg, returnedErr)
	}

	// Mark request as done and shutdown the server if scheduled and no ops are left
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()
}

// handleMetadata handles endpoint metadata requests
func (srv *Server) handleMetadata(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(resp).Encode(struct {
		ProtocolVersion string `json:"protocol-version"`
	}{
		protocolVersion,
	})
}

// handleMessage handles incoming messages
func (srv *Server) handleMessage(clt *Client, msg *Message) error {
	switch msg.msgType {
	case MsgSignalBinary:
		fallthrough
	case MsgSignalUtf8:
		fallthrough
	case MsgSignalUtf16:
		srv.handleSignal(clt, msg)

	case MsgRequestBinary:
		fallthrough
	case MsgRequestUtf8:
		fallthrough
	case MsgRequestUtf16:
		srv.handleRequest(clt, msg)

	case MsgRestoreSession:
		return srv.handleSessionRestore(clt, msg)
	case MsgCloseSession:
		return srv.handleSessionClosure(clt, msg)
	}
	return nil
}

// fulfillMsg filfills the message sending the reply
func (srv *Server) fulfillMsg(clt *Client, msg *Message, reply Payload) {
	// Send reply
	if err := clt.conn.Write(
		NewReplyMessage(msg.id, reply),
	); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsg fails the message returning an error reply
func (srv *Server) failMsg(clt *Client, msg *Message, reqErr error) {
	var replyMsg []byte
	switch err := reqErr.(type) {
	case ReqErr:
		replyMsg = NewErrorReplyMessage(msg.id, err.Code, err.Message)
	case *ReqErr:
		replyMsg = NewErrorReplyMessage(msg.id, err.Code, err.Message)
	case MaxSessConnsReachedErr:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgMaxSessConnsReached,
			msg.id,
		)
	case SessNotFoundErr:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgSessionNotFound,
			msg.id,
		)
	case SessionsDisabledErr:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgSessionsDisabled,
			msg.id,
		)
	default:
		replyMsg = NewSpecialRequestReplyMessage(
			MsgInternalError,
			msg.id,
		)
	}

	// Send request failure notification
	if err := clt.conn.Write(replyMsg); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// failMsgShutdown sends request failure reply due to current server shutdown
func (srv *Server) failMsgShutdown(clt *Client, msg *Message) {
	if err := clt.conn.Write(NewSpecialRequestReplyMessage(
		MsgReplyShutdown,
		msg.id,
	)); err != nil {
		srv.errorLog.Println("Writing failed:", err)
	}
}

// ServeHTTP will make the server listen for incoming HTTP requests
// eventually trying to upgrade them to WebSocket connections
func (srv *Server) ServeHTTP(
	resp http.ResponseWriter,
	req *http.Request,
) {
	// Reject incoming connections during shutdown, pretend the server is temporarily unavailable
	srv.opsLock.Lock()
	if srv.shutdown {
		srv.opsLock.Unlock()
		http.Error(resp, "Server shutting down", http.StatusServiceUnavailable)
		return
	}
	srv.opsLock.Unlock()

	switch req.Method {
	case "OPTIONS":
		srv.impl.OnOptions(resp)
		return
	case "WEBWIRE":
		srv.handleMetadata(resp)
		return
	}

	if !srv.impl.BeforeUpgrade(resp, req) {
		return
	}

	// Establish connection
	conn, err := srv.connUpgrader.Upgrade(resp, req)
	if err != nil {
		srv.errorLog.Print("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Register connected client
	newClient := newClientAgent(conn, req.Header.Get("User-Agent"), srv)

	srv.clientsLock.Lock()
	srv.clients = append(srv.clients, newClient)
	srv.clientsLock.Unlock()

	// Call hook on successful connection
	srv.impl.OnClientConnected(newClient)

	for {
		// Await message
		message, err := conn.Read()
		if err != nil {
			if newClient.HasSession() {
				// Decrement number of connections for this clients session
				srv.sessionRegistry.deregister(newClient)
			}

			if err.IsAbnormalCloseErr() {
				srv.warnLog.Printf("Abnormal closure error: %s", err)
			}

			newClient.unlink()
			srv.impl.OnClientDisconnected(newClient)
			return
		}

		// Parse message
		var msgObject Message
		if err := msgObject.Parse(message); err != nil {
			srv.errorLog.Println("Failed parsing message:", err)
			break
		}

		// Handle message
		if err := srv.handleMessage(newClient, &msgObject); err != nil {
			srv.errorLog.Printf("CRITICAL FAILURE: %s", err)
			break
		}
	}
}

func (srv *Server) deregisterSession(clt *Client) {
	srv.sessionRegistry.deregister(clt)
	if err := srv.sessionManager.OnSessionClosed(clt); err != nil {
		srv.errorLog.Printf("OnSessionClosed hook failed: %s", err)
	}
}

// Shutdown appoints a server shutdown and blocks the calling goroutine until the server
// is gracefully stopped awaiting all currently processed signal and request handlers to return.
// During the shutdown incoming connections are rejected with 503 service unavailable.
// Incoming requests are rejected with an error while incoming signals are just ignored
func (srv *Server) Shutdown() {
	srv.opsLock.Lock()
	srv.shutdown = true
	// Don't block if there's no currently processed operations
	if srv.currentOps < 1 {
		return
	}
	srv.opsLock.Unlock()
	<-srv.shutdownRdy
}

// SessionRegistry returns the public interface of the servers session registry
func (srv *Server) SessionRegistry() SessionRegistry {
	return srv.sessionRegistry
}
