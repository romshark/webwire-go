[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/qbeon/tmdwg-go.svg?branch=master)](https://travis-ci.org/qbeon/tmdwg-go)
[![Coverage Status](https://coveralls.io/repos/github/qbeon/tmdwg-go/badge.svg?branch=master)](https://coveralls.io/github/qbeon/tmdwg-go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/qbeon/tmdwg-go)](https://goreportcard.com/report/github.com/qbeon/tmdwg-go)
[![GoDoc](https://godoc.org/github.com/zalando/skipper?status.svg)](https://godoc.org/github.com/qbeon/tmdwg-go)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/romshark)

# Timed WaitGroup for Go

[tmdwg-go](https://github.com/qbeon/tmdwg-go) provides a **timed wait group** implementation similar to **[sync.WaitGroup](https://golang.org/pkg/sync/#WaitGroup)**.

It's purpose is simple: it blocks all goroutines that called it's `wg.Wait()` method and frees them when:
- either the timeout is reached...
- or the progress is reached

In case the timeout was reached before the progress `wg.Wait()` will return a timeout error, otherwise it'll return `nil`.

The timed wait group is fully thread safe and may safely be used concurrently from within multiple goroutines.
