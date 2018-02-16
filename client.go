package webwire

import (
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/url"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
)

func extractMessageId(message []byte) (arr [32]byte) {
	copy(arr[:], message[1:33])
	return arr
}

type Client struct {
	serverAddr string
	defaultTimeout time.Duration
	conn *websocket.Conn
	reqRegister map[[32]byte] chan []byte
	sess Session
}

// NewClient creates a new disconnected client instance. 
func NewClient(serverAddr string, defaultTimeout time.Duration) Client {
	return Client {
		serverAddr,
		defaultTimeout,
		nil,
		make(map[[32]byte] chan []byte, 0),
		Session {},
	}
}

func (clt *Client) setSession(sessionKey []byte) {
	clt.sess.key = string(sessionKey)
	clt.sess.creationDate = time.Now()
}

func (clt *Client) onRequest(payload []byte) ([]byte, error) {
	// TODO: implement real server-request handling
	// instead of current ping-pong
	return payload, nil
}

func (clt *Client) onSignal(payload []byte) error {
	// TODO: implement real server-signal handling
	return nil
}

func (clt *Client) handleRequest(message []byte) error {
	reqId := extractMessageId(message)
	// Handle server request
	result, err := clt.onRequest(message[33:])
	var msg bytes.Buffer
	if err != nil {
		msg.WriteRune(ERROR_RESP)
		msg.Write(reqId[:])
		msg.WriteString(err.Error())
	} else {
		msg.WriteRune(RESPONSE)
		msg.Write(reqId[:])
		msg.Write(result)
	}
	if err = clt.conn.WriteMessage(websocket.TextMessage, msg.Bytes());
	err != nil {
		// TODO: return typed error TransmissionFailure
		return fmt.Errorf("Couldn't send message")
	}
	return nil
}

func (clt *Client) handleSignal(message []byte) error {
	if err := clt.onSignal(message[33:]); err != nil {
		return fmt.Errorf("Signal handler failed: %s", err)
	}
	return nil
}

func (clt *Client) handleFailure(message []byte) error {
	return nil
}

func (clt *Client) handleResponse(message []byte) error {
	reqId := extractMessageId(message)

	if response, exists := clt.reqRegister[reqId]; exists {
		// Fulfill response
		response <- message[33:]
		delete(clt.reqRegister, reqId)
	}

	return nil
}

func (clt *Client) handleMessage(message []byte) error {
	if len(message) < 1 {
		return nil
	}
	switch (message[0:1][0]) {
	case RESPONSE: return clt.handleResponse(message)
	case ERROR_RESP: return clt.handleFailure(message)
	case SIGNAL: return clt.handleSignal(message)
	case REQUEST: return clt.handleRequest(message)
	default: fmt.Printf("Strange message type received: '%c'\n", message[0:1][0])
	}
	return nil
}

// Connect connects the client to the configured server and
// returns an error in case of a connection failure
func (clt *Client) Connect() (err error) {
	if clt.conn != nil {
		return nil
	}
	connUrl := url.URL {Scheme: "ws", Host: clt.serverAddr, Path: "/"}
	clt.conn, _, err = websocket.DefaultDialer.Dial(connUrl.String(), nil)
	if err != nil {
		// TODO: return typed error ConnectionFailure
		return fmt.Errorf("Could not connect: %s", err)
	}

	// Setup reader thread
	// TODO: kill reader thread on connection closure
	go func() {
		defer clt.Close()
		for {
			_, message, err := clt.conn.ReadMessage()
			if err != nil {
				fmt.Println("Failed reading message:", err)
				return
			}
			if err = clt.handleMessage(message); err != nil {
				fmt.Println("Failed handling message:", err)
				return
			}
		}
	}()

	return nil
}

func (clt *Client) sendRequest(
	messageType rune,
	payload []byte,
	timeout time.Duration,
) ([]byte, error) {
	// Connect before attempting to send the request
	if err := clt.Connect(); err != nil {
		return nil, fmt.Errorf("Couldn't connect: %s", err)
	}

	id := uuid.NewV4()
	var reqId [32]byte
	copy(reqId[:], strings.Replace(id.String(), "-", "", -1))
	var msg bytes.Buffer
	msg.WriteRune(messageType)
	msg.Write(reqId[:])
	msg.Write(payload)

	timeoutTimer := time.NewTimer(timeout).C
	responseChannel := make(chan []byte)

	// Register request
	clt.reqRegister[reqId] = responseChannel

	// Send request
	if err := clt.conn.WriteMessage(websocket.TextMessage, msg.Bytes()); err != nil {
		// TODO: return typed error TransmissionFailure
		return nil, fmt.Errorf("Couldn't send message: %s", err)
	}

	// Block until request either times out or a response is received
	select {
	case <- timeoutTimer:
		// TODO: return typed TimeoutError
		return nil, fmt.Errorf("Request timed out")
	case response := <- responseChannel:
		return response, nil
	}
}

// Request sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the request failed for some reason.
// Attempts to automatically connect to the server
// if no connection has yet been established
func (clt *Client) Request(payload []byte) ([]byte, error) {
	return clt.sendRequest(REQUEST, payload, clt.defaultTimeout)
}

// TimedRequest sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the given timeout was exceeded awaiting the response
// ar another failure occurred.
// Attempts to automatically connect to the server
// if no connection has yet been established
func (clt *Client) TimedRequest(payload []byte, timeout time.Duration) ([]byte, error) {
	return clt.sendRequest(REQUEST, payload, timeout)
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
	msg.WriteRune(SIGNAL)
	msg.Write(payload)
	if err := clt.conn.WriteMessage(websocket.TextMessage, msg.Bytes());
	err != nil {
		return err
	}
	return nil
}

// Authenticate attempts to authenticate and create a new session
// using the given credentials.
// Fails if a session is currently already active.
// Attempts to automatically connect to the server
// if no connection has yet been established
func (clt *Client) Authenticate(login, password string) error {
	if len(clt.sess.key) < 1 {
		// TODO: return typed error SessionActive
		return fmt.Errorf("another session is currently active")
	}
	if len(login) < 1 {
		return fmt.Errorf("missing login parameter")
	}
	if len(password) < 1 {
		return fmt.Errorf("missing password parameter")
	}

	// Connect before attempting authentication
	if err := clt.Connect(); err != nil {
		return fmt.Errorf("Couldn't connect: %s", err)
	}

	// Request session creation
	jsonBuff, err := json.Marshal(struct {
		Login string `json:"l"`
		Password string `json:"p"`
	} {
		login,
		password,
	})
	if err != nil {
		return fmt.Errorf("Couldn't marshal credentials: %s", err)
	}

	sessionKey, err := clt.sendRequest(
		SESS_CREATION, jsonBuff, clt.defaultTimeout,
	)
	clt.setSession(sessionKey)
	if err != nil {
		return fmt.Errorf("Couldn't request session key: %s", err)
	}

	return nil
}

// Session returns information about the current session
func (clt *Client) Session() Session {
	return Session{}
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

	if _, err := clt.sendRequest(SESS_RESTORE, sessionKey, clt.defaultTimeout);
	err != nil {
		// TODO: check for error types
		return fmt.Errorf("Session restoration request failed: %s", err)
	}
	
	return nil
}

// CloseSession closes the currently active session.
// Does nothing if there's no active session
func (clt *Client) CloseSession() error {
	if len(clt.sess.key) < 1 {
		return nil
	}

	// Connect before attempting session restoration
	if err := clt.Connect(); err != nil {
		return fmt.Errorf("Couldn't connect: %s", err)
	}

	if _, err := clt.sendRequest(SESS_CLOSURE, nil, clt.defaultTimeout);
	err != nil {
		return fmt.Errorf("Session closure request failed: %s", err)
	}

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
