package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
)

func userInterface() {
	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmd, err := buf.ReadBytes('\n')
		command := string(cmd[:len(cmd)-1])
		if err != nil {
			log.Print(err)
			continue
		}

		switch command {
		// Manually trigger the garbage collector
		case "gc":
			log.Print("Initiating garbage collection...")
			runtime.GC()
		default:
			log.Printf("Invalid command: '%s'", command)
		}
	}
}
