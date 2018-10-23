package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/qbeon/webwire-go/examples/chatroom/shared"

	wwrclt "github.com/qbeon/webwire-go/client"
)

// ChatroomClient implements the wwrclt.Implementation interface
type ChatroomClient struct {
	connection wwrclt.Client
}

// NewChatroomClient constructs and returns a new chatroom client instance
func NewChatroomClient(serverAddr url.URL) (*ChatroomClient, error) {
	newChatroomClient := &ChatroomClient{}

	// Initialize connection
	connection, err := wwrclt.NewClient(
		serverAddr,
		newChatroomClient,
		wwrclt.Options{
			DefaultRequestTimeout: 10 * time.Second,
			// Default timeout for timed requests
			ReconnectionInterval: 2 * time.Second,

			// Session info parser function must override the default one
			// for the session info object to be typed as shared.SessionInfo
			SessionInfoParser: shared.SessionInfoParser,

			// Custom loggers
			WarnLog: log.New(
				os.Stdout,
				"WARN: ",
				log.Ldate|log.Ltime|log.Lshortfile,
			),
			ErrorLog: log.New(
				os.Stderr,
				"ERR: ",
				log.Ldate|log.Ltime|log.Lshortfile,
			),
		},
		/*
			--------------------------------------------------------------
			WARNING! NEVER DISABLE CERTIFICATE VERIFICATION IN PRODUCTION!
			--------------------------------------------------------------
			InsecureSkipVerify is enabled for demonstration purposes only
			to allow the use of a self-signed localhost SSL certificate.
			Enabling this option in production is dangerous and irresponsible.
			Alternatively, you can install the "wwrexampleCA.pem" root
			certificate to make your system accept the self-signed "server.crt"
			certificate for "localhost" and disable InsecureSkipVerify.
		*/
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)
	if err != nil {
		return nil, err
	}

	newChatroomClient.connection = connection

	return newChatroomClient, nil
}

var serverAddr = flag.String("addr", "localhost:9090", "server address")
var password = flag.String("pass", "", "password")
var username = flag.String("name", "", "username")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Create a new chatroom client instance including its connection
	serverAddr := url.URL{
		Scheme: "https",
		Host:   *serverAddr,
		Path:   "/",
	}

	fmt.Printf("Connecting to %s\n", serverAddr.String())
	chatroomClient, err := NewChatroomClient(serverAddr)
	if err != nil {
		panic(err)
	}

	// Authenticate if credentials are already provided from the CLI
	if *username != "" && *password != "" {
		chatroomClient.Authenticate(*username, *password)
	}

	// Start the main loop
	chatroomClient.Start()
}
