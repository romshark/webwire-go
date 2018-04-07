package main

import (
	"flag"
	"time"

	wwrclt "github.com/qbeon/webwire-go/client"
)

// ChatroomClient implements the wwrclt.Implementation interface
type ChatroomClient struct {
	connection *wwrclt.Client
}

// NewChatroomClient constructs and returns a new chatroom client instance
func NewChatroomClient(serverAddr string) *ChatroomClient {
	newChatroomClient := &ChatroomClient{}

	// Initialize connection
	newChatroomClient.connection = wwrclt.NewClient(
		serverAddr,
		newChatroomClient,
		wwrclt.Options{
			// Address of the webwire server
			// Default timeout for timed requests
			DefaultRequestTimeout: 10 * time.Second,
			ReconnectionInterval:  2 * time.Second,
		},
	)

	return newChatroomClient
}

var serverAddr = flag.String("addr", ":9090", "server address")
var password = flag.String("pass", "", "password")
var username = flag.String("name", "", "username")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Create a new chatroom client instance including its connection
	chatroomClient := NewChatroomClient(*serverAddr)

	// Authenticate if credentials are already provided from the CLI
	if *username != "" && *password != "" {
		chatroomClient.Authenticate(*username, *password)
	}

	// Start the main loop
	chatroomClient.Start()
}
