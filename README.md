[![Go Report Card](https://goreportcard.com/badge/github.com/qbeon/webwire-go)](https://goreportcard.com/report/github.com/qbeon/webwire-go)

# WebWire for Golang
WebWire is an asynchronous duplex messaging library built on top of [WebSockets](https://developer.mozilla.org/de/docs/WebSockets).
The [webwire-go](https://github.com/qbeon/webwire-go) library provides both a client and a server implementation for the Go programming language.

## Features
WebWire provides a compact set of useful features that are not available and/or relatively hard to implement on raw WebSockets.
- **Request-Reply** - Unlike with raw websockets clients can initiate multiple simultaneous requests and receive replies asynchronously.
- **Server-side Signals** - The server can send signals to individual connected clients.
- **Client-side Signals** - Individual clients can send signals to the server.
- **Authentication & Sessions** - Both the client and server can initiate authentication and create sessions. The state of the session is automagically synchronized between the client and the server. WebWire doesn't enfore any kind of authentication technique.
- **Thread Safety** - It's safe to send messages (both on the client and on the server) from multiple goroutines simultaneously.
- **Hooks** - Various hooks provide the ability to asynchronously on different kinds of events and control the servers and clients behavior.
