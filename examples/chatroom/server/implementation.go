package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/chatroom/shared"
)

// ChatRoomServer implements the webwire.ServerImplementation interface
type ChatRoomServer struct {
	connected map[*wwr.Client]bool
	lock      sync.RWMutex
}

// NewChatRoomServer constructs a new pub-sub webwire server implementation instance
func NewChatRoomServer() *ChatRoomServer {
	return &ChatRoomServer{
		make(map[*wwr.Client]bool),
		sync.RWMutex{},
	}
}

/****************************************************************\
	Message Broadcaster
\****************************************************************/

// broadcastMessage sends a message on behalf of the given user
// to all connected clients
func (srv *ChatRoomServer) broadcastMessage(name string, msg string) {
	// Marshal message
	encoded, err := json.Marshal(shared.ChatMessage{
		User: name,
		Msg:  msg,
	})
	if err != nil {
		panic(fmt.Errorf("Couldn't marshal chat message: %s", err))
	}

	// Send message as a signal to each connected client
	srv.lock.RLock()
	log.Printf("Broadcast message to %d clients", len(srv.connected))
	for client := range srv.connected {
		// Send message as signal
		if err := client.Signal("", wwr.Payload{
			Encoding: wwr.EncodingUtf8,
			Data:     encoded,
		}); err != nil {
			log.Printf(
				"WARNING: failed sending signal to client %s : %s",
				client.Info().RemoteAddr,
				err,
			)
		}
	}
	srv.lock.RUnlock()
}

/****************************************************************\
	Authentication Handler
\****************************************************************/

// onAuth handles incoming authentication requests.
// It parses and verifies the provided credentials and either rejects the authentication
// or confirms it eventually creating a session and returning the session key
func (srv *ChatRoomServer) handleAuth(
	_ context.Context,
	client *wwr.Client,
	message *wwr.Message,
) (wwr.Payload, error) {
	credentialsText, err := message.Payload.Utf8()
	if err != nil {
		return wwr.Payload{}, wwr.ReqErr{
			Code:    "DECODING_FAILURE",
			Message: fmt.Sprintf("Failed decoding message: %s", err),
		}
	}

	log.Printf("Client attempts authentication: %s", client.Info().RemoteAddr)

	// Try to parse credentials
	var credentials shared.AuthenticationCredentials
	if err := json.Unmarshal([]byte(credentialsText), &credentials); err != nil {
		return wwr.Payload{}, fmt.Errorf("Failed parsing credentials: %s", err)
	}

	// Verify username
	password, userExists := userAccounts[credentials.Name]
	if !userExists {
		return wwr.Payload{}, wwr.ReqErr{
			Code:    "INEXISTENT_USER",
			Message: fmt.Sprintf("No such user: '%s'", credentials.Name),
		}
	}

	// Verify password
	if password != credentials.Password {
		return wwr.Payload{}, wwr.ReqErr{
			Code:    "WRONG_PASSWORD",
			Message: "Provided password is wrong",
		}
	}

	// Finally create a new session
	if err := client.CreateSession(&shared.SessionInfo{
		Username: credentials.Name,
	}); err != nil {
		return wwr.Payload{}, fmt.Errorf("Couldn't create session: %s", err)
	}

	log.Printf(
		"Created session for user %s (%s)",
		client.Info().RemoteAddr,
		credentials.Name,
	)

	// Reply to the request, use default binary encoding
	return wwr.Payload{
		Data: []byte(client.SessionKey()),
	}, nil
}

/****************************************************************\
	Message Handler
\****************************************************************/

func (srv *ChatRoomServer) handleMessage(
	_ context.Context,
	client *wwr.Client,
	message *wwr.Message,
) (wwr.Payload, error) {
	msgStr, err := message.Payload.Utf8()
	if err != nil {
		log.Printf(
			"Received invalid message from %s, couldn't convert payload to UTF8: %s",
			client.Info().RemoteAddr,
			err,
		)
		return wwr.Payload{}, nil
	}

	log.Printf(
		"Received message from %s: '%s' (%d, %s)",
		client.Info().RemoteAddr,
		msgStr,
		len(message.Payload.Data),
		message.Payload.Encoding.String(),
	)

	name := "Anonymous"
	// Try to read the name from the session
	if client.HasSession() {
		name = client.SessionInfo("username").(string)
	}

	srv.broadcastMessage(name, msgStr)

	return wwr.Payload{}, nil
}

/****************************************************************\
	Hook implementations
\****************************************************************/

// OnOptions implements the webwire.ServerImplementation interface.
// Sets HTTP access control headers to satisfy CORS
func (srv *ChatRoomServer) OnOptions(resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "WEBWIRE")
}

// OnSignal implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *ChatRoomServer) OnSignal(
	_ context.Context,
	_ *wwr.Client,
	_ *wwr.Message,
) {
}

// BeforeUpgrade implements the webwire.ServerImplementation interface.
// Must return true to ensure incoming connections are accepted
func (srv *ChatRoomServer) BeforeUpgrade(resp http.ResponseWriter, req *http.Request) bool {
	return true
}

// OnRequest implements the webwire.ServerImplementation interface.
// Receives the message and dispatches it to the according handler
func (srv *ChatRoomServer) OnRequest(
	ctx context.Context,
	client *wwr.Client,
	message *wwr.Message,
) (response wwr.Payload, err error) {
	switch message.Name {
	case "auth":
		return srv.handleAuth(ctx, client, message)
	case "msg":
		return srv.handleMessage(ctx, client, message)
	}
	return wwr.Payload{}, wwr.ReqErr{
		Code:    "BAD_REQUEST",
		Message: fmt.Sprintf("Unsupported request name: %s", message.Name),
	}
}

// OnClientConnected implements the webwire.ServerImplementation interface.
// Registers new connected clients
func (srv *ChatRoomServer) OnClientConnected(newClient *wwr.Client) {
	info := newClient.Info()
	log.Printf(
		"New client connected: %s | %s",
		info.RemoteAddr,
		info.UserAgent,
	)
	srv.lock.Lock()
	defer srv.lock.Unlock()
	srv.connected[newClient] = true
}

// OnClientDisconnected implements the webwire.ServerImplementation interface.
// Deregisters gone clients
func (srv *ChatRoomServer) OnClientDisconnected(client *wwr.Client) {
	log.Printf("Client %s disconnected", client.Info().RemoteAddr)
	srv.lock.Lock()
	defer srv.lock.Unlock()
	delete(srv.connected, client)
}
