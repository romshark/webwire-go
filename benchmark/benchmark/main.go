package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"runtime/pprof"
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

// Max benchmark duration
var argBenchDur = flag.Uint("dur", 0, "benchmark duration in seconds")

var argTracePath = flag.String("trace", "", "path to trace output file")
var argMemProfPath = flag.String("memprof", "", "path to memory profile file")
var argCPUProfPath = flag.String("cpuprof", "", "path to CPU profile file")

func writeMemoryProfile() {
	if *argMemProfPath == "" {
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
		panic(err)
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
			panic(err)
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
		if err := file.Close(); err != nil {
			panic(err)
		}
	}
}

func main() {
	parseCliArgs()

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

	defer startTracer()()
	defer startCPUProfile()()

	bench.Start()
	bench.PrintStats()

	writeMemoryProfile()
}
