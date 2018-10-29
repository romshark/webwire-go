package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"runtime/trace"
)

// Server address
var argServerAddr = flag.String("addr", ":8081", "server address")

// Number of concurrent clients
var argNumClients = flag.Uint("clients", 10, "number of concurrent clients")

// Request timeout
var argReqTimeout = flag.Uint(
	"req-timeo",
	10000,
	"default request timeout in milliseconds",
)

// Min/Max request interval
var argMinReqInterval = flag.Uint(
	"min-req-itv",
	0,
	"min interval between each request in milliseconds",
)
var argMaxReqInterval = flag.Uint(
	"max-req-itv",
	0,
	"max interval between each request in milliseconds",
)

// Min/Max payload size
var argMinPayloadSize = flag.Uint64(
	"min-pld-sz",
	1024,
	"request payload size in bytes",
)
var argMaxPayloadSize = flag.Uint64(
	"max-pld-sz",
	1024,
	"request payload size in bytes",
)

// Max benchmark duration
var argBenchDur = flag.Uint("dur", 0, "benchmark duration in seconds")

var argEnableTrace = flag.Bool("trace", false, "enable runtime trace")

func main() {
	parseCliArgs()

	// Start tracer if enabled
	if *argEnableTrace {
		f, err := os.Create("bench-trace.out")
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

	// Load options
	optionsRaw, err := ioutil.ReadFile("options.json")
	if err != nil {
		log.Fatalf("couldn't read benchmarks options: %s", err)
	}
	var options BenchmarkOptions
	if err := json.Unmarshal(optionsRaw, &options); err != nil {
		log.Fatalf("couldn't unmarshal options: %s", err)
	}

	if *argBenchDur != 0 {
		options.Duration = *argBenchDur
	}

	// Initialize and start benchmark
	bench := NewBenchmark(
		url.URL{
			Host: *argServerAddr,
		},
		options,
	)
	bench.Start()
	bench.PrintStats()
}
