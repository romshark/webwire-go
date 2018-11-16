package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func humanReadableDataSize(numBytes float64) (float64, string) {
	if math.IsInf(numBytes, 1) || math.IsNaN(numBytes) {
		return 0, "NaN"
	}

	switch {
	case numBytes >= float64(1024*1024*1024*1024*1024):
		return numBytes / float64((1024 * 1024 * 1024 * 1024 * 1024)), "PiB"
	case numBytes >= float64(1024*1024*1024*1024):
		return numBytes / float64((1024 * 1024 * 1024 * 1024)), "TiB"
	case numBytes >= float64(1024*1024*1024):
		return numBytes / float64((1024 * 1024 * 1024)), "GiB"
	case numBytes >= float64(1024*1024):
		return numBytes / float64((1024 * 1024)), "MiB"
	case numBytes >= float64(1024):
		return numBytes / float64(1024), "KiB"
	}
	return numBytes, "Bytes"
}

func random(min, max int64) int64 {
	if min == max {
		return min
	}
	rand.Seed(time.Now().Unix())
	return rand.Int63n(max-min) + min
}

func parseCliArgs() {
	// Parse command line arguments
	flag.Parse()
	err := false

	// Validate server address
	if len(*argServerAddr) < 1 {
		err = true
		log.Printf(
			"INVALID ARGS: invalid server address ('%s')",
			*argServerAddr,
		)
	}

	// Validate the number of concurrent clients
	if *argNumClients < 1 {
		err = true
		log.Printf("INVALID ARGS: number of concurrent clients cannot be zero")
	}

	// Validate request timeout
	if *argReqTimeout < 1 {
		err = true
		log.Printf(
			"INVALID ARGS: request timeout cannot be smaller 1 millisecond",
		)
	}

	// Validate request interval range
	if *argMinReqInterval > *argMaxReqInterval {
		err = true
		log.Printf(
			"INVALID ARGS: min request interval (%d) grater max parameter (%d)",
			*argMinReqInterval,
			*argMaxReqInterval,
		)
	}

	if err {
		os.Exit(1)
	}
}

func randomReqIntervalSleep(min, max uint) time.Duration {
	if min == max {
		return time.Duration(min) * time.Millisecond
	}

	return time.Duration(random(
		int64(min), int64(max),
	)) * time.Millisecond
}
