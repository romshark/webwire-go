package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/qbeon/webwire-go"
)

func promptCreds() (username, password string) {
	reader := bufio.NewReader(os.Stdin)

	// Prompt username
	fmt.Print("  Username: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("Failed reading username prompt: %s", err))
	}
	username = input[:len(input)-1]

	// Prompt password
	fmt.Print("  Password: ")
	input, err = reader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("Failed reading password prompt: %s", err))
	}
	password = input[:len(input)-1]

	return username, password
}

// Start runs the main loop blocking the calling goroutine
func (clt *ChatroomClient) Start() {
	defer clt.connection.Close()

	reader := bufio.NewReader(os.Stdin)
MAINLOOP:
	for {
		input, _ := reader.ReadString('\n')

		// Remove new-line character
		input = input[:len(input)-1]

		if len(input) < 1 {
			continue
		}
		switch input {
		case ":x":
			fmt.Println("Closing connection...")
			break MAINLOOP
		case ":login":
			username, password := promptCreds()
			clt.Authenticate(username, password)
		case ":logout":
			// Check if even authenticated
			if clt.connection.Session().Key == "" {
				fmt.Println("Not authenticated, no need to logout")
				break
			}
			// Try to close the session
			if err := clt.connection.CloseSession(); err != nil {
				log.Printf("WARNING: Session destruction failed: %s", err)
			}
			fmt.Println("Logged out, you're anonymous now")
		case ":disconnect":
			clt.connection.Close()
		case ":connect":
			if err := clt.connection.Connect(); err != nil {
				fmt.Printf("Error while connecting: %s\n", err)
			}
		default:
			// Send the message and await server reply
			// for the message to be considered posted
			go func() {
				if _, err := clt.connection.Request(
					context.Background(),
					"msg",
					webwire.NewPayload(webwire.EncodingBinary, []byte(input)),
				); err != nil {
					log.Printf("WARNING: Couldn't send message: %s", err)
				}
			}()
		}
	}
}
