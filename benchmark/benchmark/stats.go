package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Stats represents benchmark statistics
type Stats struct {
	options BenchmarkOptions

	lock         sync.Mutex
	totalReqTime time.Duration
	maxReqTime   time.Duration
	minReqTime   time.Duration

	totalReqInterval time.Duration
	maxReqInterval   time.Duration
	minReqInterval   time.Duration

	totalReqsPerformed uint64
	totalReqTimeouts   uint64
	totalBytesSent     uint64
	totalBytesReceived uint64
}

// NewStats creates a new statistics recorder
func NewStats(options BenchmarkOptions) *Stats {
	return &Stats{
		options:            options,
		lock:               sync.Mutex{},
		totalReqTime:       time.Duration(0),
		maxReqTime:         time.Duration(0),
		minReqTime:         time.Duration(0),
		totalReqInterval:   time.Duration(0),
		maxReqInterval:     time.Duration(0),
		minReqInterval:     time.Duration(0),
		totalReqsPerformed: uint64(0),
		totalReqTimeouts:   uint64(0),
		totalBytesSent:     uint64(0),
		totalBytesReceived: uint64(0),
	}
}

// Print prints computed statistics to standard output
func (st *Stats) Print() {
	reqPerSec := float64(st.totalReqsPerformed) / float64(st.options.Duration)
	bytesPerSec := float64(uint64(st.totalBytesSent) /
		uint64(st.options.Duration))

	var avgReqTime time.Duration
	var avgReqInterval time.Duration
	var avgPayloadSize uint64
	var timeoutRate float64
	if st.totalReqsPerformed > 0 {
		avgReqTime = st.totalReqTime / time.Duration(st.totalReqsPerformed)
		avgReqInterval = st.totalReqInterval /
			time.Duration(st.totalReqsPerformed)
		avgPayloadSize = st.totalBytesSent / st.totalReqsPerformed
		timeoutRate = float64(st.totalReqTimeouts) /
			float64(st.totalReqsPerformed) * 100
	}

	if math.IsInf(timeoutRate, 1) {
		timeoutRate = 0.0
	}

	fmt.Println(" ")
	log.Printf("  Benchmark finished (%ds)\n", st.options.Duration)
	fmt.Println(" ")
	fmt.Printf("  Max Concurrent Connections:  %d\n", st.options.MaxClients)
	fmt.Printf("  Requests performed:          %d\n", st.totalReqsPerformed)
	fmt.Printf("  Requests timed out:          %d (%.2f%%)\n",
		st.totalReqTimeouts,
		timeoutRate,
	)
	fmt.Println(" ")

	dataSent, sentUnits := humanReadableDataSize(float64(st.totalBytesSent))
	dataReceived, receivedUnits := humanReadableDataSize(
		float64(st.totalBytesReceived),
	)
	avgPayloadSizeNum, avgPayloadSizeUnits := humanReadableDataSize(
		float64(avgPayloadSize),
	)

	fmt.Printf("  Data sent:                   %.2f %s (%d bytes)\n",
		dataSent,
		sentUnits,
		st.totalBytesSent,
	)
	fmt.Printf("  Data received:               %.2f %s (%d bytes)\n",
		dataReceived,
		receivedUnits,
		st.totalBytesReceived,
	)
	fmt.Printf("  Avg payload size:            %.2f %s\n",
		avgPayloadSizeNum,
		avgPayloadSizeUnits,
	)
	fmt.Println(" ")

	numPerSec, units := humanReadableDataSize(bytesPerSec)

	fmt.Printf("  Avg req itv:                 %s\n", avgReqInterval)
	fmt.Printf("  Max req itv:                 %s\n", st.maxReqInterval)
	fmt.Printf("  Min req itv:                 %s\n", st.minReqInterval)
	fmt.Println(" ")

	fmt.Printf("  Avg req time:                %s\n", avgReqTime)
	fmt.Printf("  Max req time:                %s\n", st.maxReqTime)
	fmt.Printf("  Min req time:                %s\n", st.minReqTime)
	fmt.Println(" ")

	fmt.Printf("  Req/s:                       %.0f\n", reqPerSec)
	fmt.Printf("  Bytes/s:                     %.0f\n", bytesPerSec)
	fmt.Printf("  Throughput:                  %.2f %s/s\n", numPerSec, units)
}

// RecordRequest records request statistics
func (st *Stats) RecordRequest(
	interval,
	reqTime time.Duration,
	requestSize,
	replySize int,
) {
	st.lock.Lock()

	st.totalReqsPerformed++
	st.totalBytesSent += uint64(requestSize)
	st.totalBytesReceived += uint64(replySize)

	st.totalReqTime += reqTime
	if st.minReqTime == 0 || reqTime < st.minReqTime {
		st.minReqTime = reqTime
	}
	if reqTime > st.maxReqTime {
		st.maxReqTime = reqTime
	}

	st.totalReqInterval += interval
	if st.minReqInterval == 0 || interval < st.minReqInterval {
		st.minReqInterval = interval
	}
	if interval > st.maxReqInterval {
		st.maxReqInterval = interval
	}

	st.lock.Unlock()
}

// RecordTimedoutRequest records a timed out request
func (st *Stats) RecordTimedoutRequest() {
	atomic.AddUint64(&st.totalReqTimeouts, 1)
}
