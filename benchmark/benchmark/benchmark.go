package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	wwr "github.com/qbeon/webwire-go"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
)

// Benchmark represents the request benchmark
type Benchmark struct {
	serverAddr    url.URL
	stats         *Stats
	options       BenchmarkOptions
	staticPayload []byte
}

// BenchmarkOptions represents the benchmark options
type BenchmarkOptions struct {
	Duration             uint `json:"duration"`
	MaxClients           uint `json:"max-connections"`
	MaxRequestsPerClient uint `json:"max-reqs-per-connection"`
	MinPayloadSize       uint `json:"min-payload-size"`
	MaxPayloadSize       uint `json:"max-payload-size"`
	RequestTimeout       uint `json:"req-timeo"`
	MinReqInterval       uint `json:"min-req-interval"`
	MaxReqInterval       uint `json:"max-req-interval"`
}

// NewBenchmark creates a new benchmark instance
func NewBenchmark(serverAddr url.URL, options BenchmarkOptions) *Benchmark {
	benchmark := &Benchmark{
		serverAddr: serverAddr,
		stats:      NewStats(options),
		options:    options,
	}

	if options.MinPayloadSize == options.MaxPayloadSize {
		benchmark.staticPayload = make([]byte, options.MinPayloadSize)
	}

	return benchmark
}

// Start starts
func (bc *Benchmark) Start() {
	clients := make([]*Client, bc.options.MaxClients)
	for i := uint(0); i < bc.options.MaxClients; i++ {
		clients[i] = NewClient(
			bc.serverAddr,
			time.Duration(bc.options.RequestTimeout)*time.Millisecond,
			&wwrfasthttp.ClientTransport{},
		)
	}

	log.Printf("All clients (%d) operational", bc.options.MaxClients)

	wg := sync.WaitGroup{}
	wg.Add(len(clients))

	timeoutTriggered := false
	triggerLock := sync.Mutex{}

	time.AfterFunc(time.Duration(bc.options.Duration)*time.Second, func() {
		log.Print("Finishing...")
		triggerLock.Lock()
		timeoutTriggered = true
		triggerLock.Unlock()
	})

	for _, clt := range clients {
		c := clt
		go func() {
			// Initialize the semaphore synchronizing the number of simultaneous
			// requests per client
			requestSlots := semaphore.NewWeighted(int64(
				bc.options.MaxRequestsPerClient,
			))

			for {
				if err := requestSlots.Acquire(
					context.Background(),
					1,
				); err != nil {
					panic(fmt.Errorf("couldn't acquire request slot: %s", err))
				}

				// Check for shutdown
				triggerLock.Lock()
				if timeoutTriggered {
					triggerLock.Unlock()
					break
				}
				triggerLock.Unlock()

				// Wait before sending the next request
				interval := randomReqIntervalSleep(
					bc.options.MinReqInterval,
					bc.options.MaxReqInterval,
				)
				if interval != 0 {
					time.Sleep(interval)
				}

				// Capture time of request beginning
				start := time.Now()

				// Generate payload data
				var payloadData []byte
				if bc.options.MinPayloadSize == bc.options.MaxPayloadSize {
					payloadData = bc.staticPayload
				} else {
					payloadData = make([]byte, random(
						int64(bc.options.MinPayloadSize),
						int64(bc.options.MaxPayloadSize),
					))
				}

				// Send request and await reply
				reply, err := c.Request(wwr.Payload{
					Encoding: wwr.EncodingBinary,
					Data:     payloadData,
				})

				// Compute elapsed time since request start
				elapsed := time.Since(start)

				requestSlots.Release(1)

				// Investigate request error
				switch err := err.(type) {
				case nil:
				case wwr.TimeoutErr:
					bc.stats.RecordTimedoutRequest()
					continue
				default:
					panic(fmt.Errorf(
						"ERROR: Unexpected request error: %s",
						err,
					))
				}

				// Determine the length of the payload and close the reply to
				// release the buffer
				payloadLength := len(reply.Payload())
				reply.Close()

				// Update stats
				bc.stats.RecordRequest(
					interval,
					elapsed,
					len(payloadData),
					payloadLength,
				)
			}

			c.Close()
			wg.Done()
		}()
	}
	wg.Wait()
}

// PrintStats prints the statistics to standard output
func (bc *Benchmark) PrintStats() {
	bc.stats.Print()
}
