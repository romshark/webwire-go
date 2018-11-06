package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"runtime/trace"
	"syscall"
	"time"

	wwr "github.com/qbeon/webwire-go"
)

func listenOsSignals(server wwr.Server) {
	// Listen for OS signals and shutdown server in case of demanded termination
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-osSignals
	log.Printf("Termination (%s), shutting down...", sig)
	if err := server.Shutdown(); err != nil {
		log.Printf("Error during server shutdown: %s", err)
	}
	log.Println("Server gracefully terminated")
}

var argHostAddr = flag.String("addr", "localhost:8081", "server host address")
var argTransport = flag.String("transport", "http", "http / https")
var argCertFilePath = flag.String(
	"tls_cert",
	"./server.crt",
	"path to the TLS certificate file",
)
var argPrivateKeyFilePath = flag.String(
	"tls_key",
	"./server.key",
	"path to the TLS private-key file",
)
var argReadTimeout = flag.Uint64("rtimeo", 10, "read timeout in seconds")
var argTracePath = flag.String("trace", "", "path to trace output file")
var argMemProfPath = flag.String("memprof", "", "path to memory profile file")
var argCPUProfPath = flag.String("cpuprof", "", "path to CPU profile file")

func writeMemoryProfile() {
	if *argCPUProfPath == "" {
		return
	}

	file, err := os.Create(*argMemProfPath)
	if err != nil {
		panic(err)
	}
	if err := pprof.WriteHeapProfile(file); err != nil {
		log.Fatal(err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

func startCPUProfile() func() {
	if *argCPUProfPath == "" {
		return func() {}
	}
	file, err := os.Create(*argCPUProfPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := pprof.StartCPUProfile(file); err != nil {
		log.Fatal(err)
	}
	return func() {
		pprof.StopCPUProfile()
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func startTracer() func() {
	// Start tracer if enabled
	if *argTracePath == "" {
		return func() {}
	}
	file, err := os.Create(*argTracePath)
	if err != nil {
		panic(err)
	}

	if err := trace.Start(file); err != nil {
		panic(err)
	}
	return func() {
		trace.Stop()
		file.Close()
	}
}

func main() {
	// Parse command line arguments
	flag.Parse()

	defer startCPUProfile()()
	defer startTracer()()

	server, err := newServer(settings{
		HostAddress:        *argHostAddr,
		Transport:          *argTransport,
		CertFilePath:       *argCertFilePath,
		PrivateKeyFilePath: *argPrivateKeyFilePath,
		ReadTimeout:        time.Duration(*argReadTimeout) * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	addr := server.Address()
	log.Printf("Listening on %s", addr.String())

	go listenOsSignals(server)
	go userInterface()

	// Launch echo server
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}

	writeMemoryProfile()
}
