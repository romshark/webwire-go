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
var argEnableTrace = flag.Bool("trace", false, "enable runtime trace")
var argEnableHTTPS = flag.Bool("https", false, "enable TLS")
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

func main() {
	// Parse command line arguments
	flag.Parse()

	// Start tracer if enabled
	if *argEnableTrace {
		f, err := os.Create("server-trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	server, err := newServer(settings{
		HostAddress:        *argHostAddr,
		HTTPSEnabled:       *argEnableHTTPS,
		CertFilePath:       *argCertFilePath,
		PrivateKeyFilePath: *argPrivateKeyFilePath,
		ReadTimeout:        time.Duration(*argReadTimeout) * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s", server.Address())

	go listenOsSignals(server)
	go userInterface()

	// Launch echo server
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}

	profileFile, err := os.Create("./benchmark.profile")
	if err != nil {
		panic(err)
	}
	defer profileFile.Close()
	pprof.WriteHeapProfile(profileFile)
}
