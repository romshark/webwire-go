package main

import (
	"bufio"
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

func mainLoop() {
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
			*username, *password = promptCreds()
			authenticate()
		case ":logout":
			// Check if even authenticated
			if client.Session().Key == "" {
				fmt.Println("Not authenticated, no need to logout")
				break
			}
			// Try to close the session
			if err := client.CloseSession(); err != nil {
				log.Printf("WARNING: Session destruction failed: %s", err)
			}
			fmt.Println("Logged out, you're anonymous now")
		default:
			// Send the message and await server reply for the message to be considered posted
			go func() {
				if _, err := client.Request("msg", webwire.Payload{
					Data: []byte(input),
				}); err != nil {
					log.Printf("WARNING: Couldn't send message: %s", err)
				}
			}()
		}
	}
}
