package client

import (
	webwire "github.com/qbeon/webwire-go"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

const supportedProtocolVersion = "1.0"

// Hooks represents all callback hook functions
type Hooks struct {
	// OnServerSignal is an optional callback.
	// It's invoked when the webwire client receives a signal from the server
	OnServerSignal func([]byte)

	// OnSessionCreated is an optional callback.
	// It's invoked when the webwire client receives a new session
	OnSessionCreated func(*webwire.Session)

	// OnSessionClosed is an optional callback.
	// It's invoked when the clients session was closed
	// either by the server or by himself
	OnSessionClosed func()
}

// SetDefaults sets undefined required hooks
func (hooks *Hooks) SetDefaults() {
	if hooks.OnServerSignal == nil {
		hooks.OnServerSignal = func(_ []byte) {}
	}

	if hooks.OnSessionCreated == nil {
		hooks.OnSessionCreated = func(_ *webwire.Session) {}
	}

	if hooks.OnSessionClosed == nil {
		hooks.OnSessionClosed = func() {}
	}
}

// replyObj is used by the request registry
type replyObj struct {
	Reply []byte
	Error *webwire.Error
}

func extractMessageIdentifier(message []byte) (arr [32]byte) {
	copy(arr[:], message[1:33])
	return arr
}

// Client represents an instance of one of the servers clients
type Client struct {
	serverAddr     string
	defaultTimeout time.Duration
	hooks          Hooks
	conn           *websocket.Conn
	session        *webwire.Session

	reqRegisterLock sync.Mutex
	reqRegister     map[[32]byte]chan replyObj

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
		defaultTimeout,
		hooks,
		nil,
		nil,
		sync.Mutex{},
		make(map[[32]byte]chan replyObj),
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

func (clt *Client) handleSessionCreated(message []byte) {
	// Set new session
	var session webwire.Session

	if err := json.Unmarshal(message, &session); err != nil {
		clt.errorLog.Printf("Failed unmarshalling session object: %s", err)
		return
	}

	clt.session = &session
	clt.hooks.OnSessionCreated(&session)
}

func (clt *Client) handleSessionClosed() {
	// Destroy local session
	clt.session = nil

	clt.hooks.OnSessionClosed()
}

func (clt *Client) handleFailure(message []byte) {
	requestIdentifier := extractMessageIdentifier(message)

	if reply, exists := clt.reqRegister[requestIdentifier]; exists {
		// Decode error
		var err webwire.Error
		if err := json.Unmarshal(message[33:], &err); err != nil {
			clt.errorLog.Printf("Failed unmarshalling error reply: %s", err)
		}

		// Fail request
		reply <- replyObj{
			Reply: nil,
			Error: &err,
		}
		clt.reqRegisterLock.Lock()
		delete(clt.reqRegister, requestIdentifier)
		clt.reqRegisterLock.Unlock()
	}
}

func (clt *Client) handleReply(message []byte) {
	requestIdentifier := extractMessageIdentifier(message)

	if reply, exists := clt.reqRegister[requestIdentifier]; exists {
		// Fulfill request
		reply <- replyObj{
			Reply: message[33:],
			Error: nil,
		}
		clt.reqRegisterLock.Lock()
		delete(clt.reqRegister, requestIdentifier)
		clt.reqRegisterLock.Unlock()
	}
}

func (clt *Client) handleMessage(message []byte) error {
	if len(message) < 1 {
		return nil
	}
	switch message[0:1][0] {
	case webwire.MsgReply:
		clt.handleReply(message)
	case webwire.MsgErrorReply:
		clt.handleFailure(message)
	case webwire.MsgSignal:
		clt.hooks.OnServerSignal(message[1:])
	case webwire.MsgSessionCreated:
		clt.handleSessionCreated(message[1:])
	case webwire.MsgSessionClosed:
		clt.handleSessionClosed()
	default:
		clt.warningLog.Printf(
			"Strange message type received: '%c'\n",
			message[0:1][0],
		)
	}
	return nil
}

// verifyProtocolVersion requests the endpoint metadata
// to verify the server is running a supported protocol version
func (clt *Client) verifyProtocolVersion() error {
	// Initialize HTTP client
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	request, err := http.NewRequest(
		"WEBWIRE", "http://"+clt.serverAddr+"/", nil,
	)
	if err != nil {
		return fmt.Errorf("Couldn't create HTTP metadata request: %s", err)
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Endpoint metadata request failed: %s", err)
	}

	// Read response body
	defer response.Body.Close()
	encodedData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Couldn't read metadata response body: %s", err)
	}

	// Unmarshal response
	var metadata struct {
		ProtocolVersion string `json:"protocol-version"`
	}
	if err := json.Unmarshal(encodedData, &metadata); err != nil {
		return fmt.Errorf(
			"Couldn't parse HTTP metadata response ('%s'): %s",
			string(encodedData),
			err,
		)
	}

	// Verify metadata
	if metadata.ProtocolVersion != supportedProtocolVersion {
		return fmt.Errorf(
			"Unsupported protocol version: %s (%s is supported by this client)",
			metadata.ProtocolVersion,
			supportedProtocolVersion,
		)
	}

	return nil
}

// Connect connects the client to the configured server and
// returns an error in case of a connection failure
func (clt *Client) Connect() (err error) {
	if clt.conn != nil {
		return nil
	}

	if err := clt.verifyProtocolVersion(); err != nil {
		return err
	}

	connURL := url.URL{Scheme: "ws", Host: clt.serverAddr, Path: "/"}
	clt.conn, _, err = websocket.DefaultDialer.Dial(connURL.String(), nil)
	if err != nil {
		// TODO: return typed error ConnectionFailure
		return fmt.Errorf("Could not connect: %s", err)
	}

	// Setup reader thread
	go func() {
		defer clt.Close()
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

	return nil
}

func (clt *Client) sendRequest(
	messageType rune,
	payload []byte,
	timeout time.Duration,
) ([]byte, *webwire.Error) {
	// Connect before attempting to send the request
	if err := clt.Connect(); err != nil {
		return nil, &webwire.Error{
			Message: fmt.Sprintf("Couldn't connect: %s", err),
		}
	}

	id := uuid.NewV4()
	var requestIdentifier [32]byte
	copy(requestIdentifier[:], strings.Replace(id.String(), "-", "", -1))
	var msg bytes.Buffer
	msg.WriteRune(messageType)
	msg.Write(requestIdentifier[:])
	msg.Write(payload)

	timeoutTimer := time.NewTimer(timeout).C
	responseChannel := make(chan replyObj)

	// Register request
	clt.reqRegisterLock.Lock()
	clt.reqRegister[requestIdentifier] = responseChannel
	clt.reqRegisterLock.Unlock()

	// Send request
	err := clt.conn.WriteMessage(websocket.TextMessage, msg.Bytes())
	if err != nil {
		// TODO: return typed error TransmissionFailure
		return nil, &webwire.Error{
			Message: fmt.Sprintf("Couldn't send message: %s", err),
		}
	}

	// Block until request either times out or a response is received
	select {
	case <-timeoutTimer:
		// TODO: return typed TimeoutError
		return nil, &webwire.Error{
			Message: fmt.Sprintf("Request timed out"),
		}
	case reply := <-responseChannel:
		if reply.Error != nil {
			return nil, reply.Error
		}
		return reply.Reply, nil
	}
}

// Request sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the request failed for some reason.
// Attempts to automatically connect to the server
// if no connection has yet been established
func (clt *Client) Request(payload []byte) ([]byte, *webwire.Error) {
	return clt.sendRequest(webwire.MsgRequest, payload, clt.defaultTimeout)
}

// TimedRequest sends a request containing the given payload to the server
// and asynchronously returns the servers reply
// blocking the calling goroutine.
// Returns an error if the given timeout was exceeded awaiting the response
// ar another failure occurred.
// Attempts to automatically connect to the server
// if no connection has yet been established
func (clt *Client) TimedRequest(
	payload []byte,
	timeout time.Duration,
) ([]byte, *webwire.Error) {
	return clt.sendRequest(webwire.MsgRequest, payload, timeout)
}

// Signal sends a signal containing the given payload to the server.
// Attempts to automatically connect to the server
// if no connection has yet been established
func (clt *Client) Signal(payload []byte) error {
	// Connect before attempting to send the signal
	if err := clt.Connect(); err != nil {
		return fmt.Errorf("Couldn't connect to server")
	}

	var msg bytes.Buffer
	msg.WriteRune(webwire.MsgSignal)
	msg.Write(payload)
	return clt.conn.WriteMessage(websocket.TextMessage, msg.Bytes())
}

// Session returns information about the current session
func (clt *Client) Session() webwire.Session {
	if clt.session == nil {
		return webwire.Session{}
	}
	return *clt.session
}

// RestoreSession tries to restore the previously opened session
// Fails if a session is currently already active
// Attempts to automatically connect to the server
// if no connection has yet been established
func (clt *Client) RestoreSession(sessionKey []byte) error {
	// Connect before attempting session restoration
	if err := clt.Connect(); err != nil {
		return fmt.Errorf("Couldn't connect: %s", err)
	}

	if _, err := clt.sendRequest(
		webwire.MsgRestoreSession,
		sessionKey,
		clt.defaultTimeout,
	); err != nil {
		// TODO: check for error types
		return fmt.Errorf("Session restoration request failed: %s", err)
	}

	return nil
}

// CloseSession closes the currently active session.
// Does nothing if there's no active session
func (clt *Client) CloseSession() error {
	if clt.conn == nil {
		return fmt.Errorf("Cannot close a session of a disconnected client")
	}

	if clt.session == nil {
		return nil
	}

	if _, err := clt.sendRequest(
		webwire.MsgCloseSession,
		nil,
		clt.defaultTimeout,
	); err != nil {
		return fmt.Errorf("Session destruction request failed: %s", err)
	}

	// Reset session locally after destroying it on the server
	clt.session = nil

	return nil
}

// Close gracefully closes the connection.
// Does nothing if the client isn't connected
func (clt *Client) Close() {
	if clt.conn == nil {
		return
	}
	clt.conn.Close()
}
