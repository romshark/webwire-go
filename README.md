[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/qbeon/webwire-go.svg?branch=master)](https://travis-ci.org/qbeon/webwire-go)
[![Coverage Status](https://coveralls.io/repos/github/qbeon/webwire-go/badge.svg?branch=master)](https://coveralls.io/github/qbeon/webwire-go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/qbeon/webwire-go)](https://goreportcard.com/report/github.com/qbeon/webwire-go)
[![GoDoc](https://godoc.org/github.com/zalando/skipper?status.svg)](https://godoc.org/github.com/qbeon/webwire-go)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/romshark)

# WebWire for Golang
WebWire is a high-level asynchronous duplex messaging library built on top of [WebSockets](https://developer.mozilla.org/de/docs/WebSockets) and an open source binary message protocol with builtin authentication and support for UTF8 and UTF16 encoding.
The [webwire-go](https://github.com/qbeon/webwire-go) library provides both a client and a server implementation for the Go programming language. An official [JavaScript client](https://github.com/qbeon/webwire-js) implementation is also available. WebWire provides a compact set of useful features that are not available and/or cumbersome to implement on raw WebSockets such as Request-Reply, Sessions, Thread-safety etc.

#### Table of Contents
- [Introduction](https://github.com/qbeon/webwire-go#webwire-for-golang)
- [Installation](https://github.com/qbeon/webwire-go#installation)
- [Contribution](https://github.com/qbeon/webwire-go#contribution)
  - [Maintainers](https://github.com/qbeon/webwire-go#maintainers)
- [WebWire Binary Protocol](https://github.com/qbeon/webwire-go#webwire-binary-protocol)
- [Examples](https://github.com/qbeon/webwire-go#examples)
- [Features](https://github.com/qbeon/webwire-go#features)
  - [Request-Reply](https://github.com/qbeon/webwire-go#request-reply)
  - [Client-side Signals](https://github.com/qbeon/webwire-go#client-side-signals)
  - [Server-side Signals](https://github.com/qbeon/webwire-go#server-side-signals)
  - [Namespaces](https://github.com/qbeon/webwire-go#namespaces)
  - [Sessions](https://github.com/qbeon/webwire-go#sessions)
  - [Automatic Session Restoration](https://github.com/qbeon/webwire-go/blob/master/README.md#automatic-session-restoration)
  - [Automatic Connection Maintenance](https://github.com/qbeon/webwire-go/blob/master/README.md#automatic-connection-maintenance)
  - [Thread-Safety](https://github.com/qbeon/webwire-go#thread-safety)
  - [Hooks](https://github.com/qbeon/webwire-go#hooks)
    - [Server-side Hooks](https://github.com/qbeon/webwire-go#server-side-hooks)
    - [Client-side Hooks](https://github.com/qbeon/webwire-go#client-side-hooks)
  - [Graceful Shutdown](https://github.com/qbeon/webwire-go#graceful-shutdown)
  - [Seamless JavaScript Support](https://github.com/qbeon/webwire-go#seamless-javascript-support)
- [Dependencies](https://github.com/qbeon/webwire-go#dependencies)


## Installation
Choose any stable release from [the available release tags](https://github.com/qbeon/webwire-go/releases) and copy the source code into your project's vendor directory: ```$YOURPROJECT/vendor/github.com/qbeon/webwire-go```

## Contribution
Feel free to report bugs and propose new features or changes in the [Issues section](https://github.com/qbeon/webwire-go/issues).

### Maintainers
| Maintainer | Role | Specialization |
| :--- | :--- | :--- |
| **[Roman Sharkov](https://github.com/romshark)** | Core Maintainer | Dev (Go, JavaScript) |
| **[Daniil Trishkin](https://github.com/FromZeus)** | CI Maintainer | DevOps |

## WebWire Binary Protocol
WebWire is built for speed and portability implementing an open source binary protocol.
![Protocol Subset Diagram](https://github.com/qbeon/webwire-go/blob/master/docs/img/wwr_msgproto_diagram.svg)

The first byte defines the [type of the message](https://github.com/qbeon/webwire-go/blob/master/message.go#L71). Requests and replies contain an incremental 8-byte identifier that must be unique in the context of the senders' session. A 0 to 255 bytes long 7-bit ASCII encoded name is contained in the header of a signal or request message.
A header-padding byte is applied in case of UTF16 payload encoding to properly align the payload sequence.
Fraudulent messages are recognized by analyzing the message length, out-of-range memory access attacks are therefore prevented.

## Examples
- **[Echo](https://github.com/qbeon/webwire-go/tree/master/examples/echo)** - Demonstrates a simple request-reply implementation.

- **[Pub-Sub](https://github.com/qbeon/webwire-go/tree/master/examples/pubsub)** - Demonstrantes a simple publisher-subscriber tolopology implementation.

- **[Chat Room](https://github.com/qbeon/webwire-go/tree/master/examples/chatroom)** - Demonstrates advanced use of the library. The corresponding [JavaScript Chat Room Client](https://github.com/qbeon/webwire-js/tree/master/examples/chatroom-client-vue) is implemented with the [Vue.js framework](https://vuejs.org/).

## Features
### Request-Reply
Clients can initiate multiple simultaneous requests and receive replies asynchronously. Requests are multiplexed through the connection similar to HTTP2 pipelining.

```go
// Send a request to the server, will block the goroutine until replied
reply, err := client.Request("", wwr.Payload{Data: []byte("sudo rm -rf /")})
if err != nil {
  // Oh oh, request failed for some reason!
}
reply // Here we go!
 ```

Timed requests will timeout and return an error if the server doesn't manage to reply within the specified time frame.

```go
// Send a request to the server, will block the goroutine for 200ms at max
reply, err := client.TimedRequest("", wwr.Payload{Data: []byte("hurry up!")}, 200*time.Millisecond)
if err != nil {
  // Probably timed out!
}
reply // Just in time!
```

### Client-side Signals
Individual clients can send signals to the server. Signals are one-way messages guaranteed to arrive, though they're not guaranteed to be processed like requests are. In cases such as when the server is being shut down, incoming signals are ignored by the server and dropped while requests will acknowledge the failure.

```go
// Send signal to server
err := client.Signal("eventA", wwr.Payload{Data: []byte("something")})
```

### Server-side Signals
The server also can send signals to individual connected clients.

```go
func onRequest(ctx context.Context) (wwr.Payload, error) {
  msg := ctx.Value(wwr.Msg).(wwr.Message)
  // Send a signal to the client before replying to the request
  msg.Client.Signal("", wwr.Payload{Data: []byte("ping!")})
  return wwr.Payload{}, nil
}
```

### Namespaces
Different kinds of requests and signals can be differentiated using the builtin namespacing feature.

```go
func onRequest(ctx context.Context) (wwr.Payload, error) {
  msg := ctx.Value(wwr.Msg).(wwr.Message)
  switch msg.Name {
  case "auth":
    // Authentication request
  case "query":
    // Query request
  }
  return wwr.Payload{}, nil
}
```
```go
func onSignal(ctx context.Context) {
  msg := ctx.Value(wwr.Msg).(wwr.Message)
  switch msg.Name {
  case "event A":
    // handle event A
  case "event B":
    // handle event B
  }
}
```

### Sessions
Individual connections can get sessions assigned to identify them. The state of the session is automagically synchronized between the client and the server. WebWire doesn't enforce any kind of authentication technique though, it just provides you a way to authenticate a connection. WebWire also doesn't enforce any kind of session storage, it's up to the user to implement any kind of volatile or persistent session storage, be it a database or a simple map.

```go
func onRequest(ctx context.Context) (wwr.Payload, error) {
  msg := ctx.Value(wwr.Msg).(wwr.Message)
  client := msg.Client
  // Verify credentials
  if string(msg.Payload.Data) != "secret:pass" {
    return wwr.Payload{}, wwr.Error {
      Code: "WRONG_CREDENTIALS",
      Message: "Incorrect username or password, try again"
    }
  }
  // Create session, will automatically synchronize to the client
  err := client.CreateSession(/*arbitrary data*/); err != nil {
    return wwr.Payload{}, fmt.Errorf("Couldn't create session for some reason")
  }
  client.Session // return wwr.Payload{}, nil
}
```

### Automatic Session Restoration
The client will automatically try to restore the previously opened session during connection establishment when getting disconnected without explicitly closing the session before.

```go
// Will automatically restore the previous session if there was any
err := client.Connect()
```

The session can also be restored manually given its key assuming the server didn't yet delete it. Session restoration will fail and return an error if the provided key doesn't correspond to any active session on the server or else if there's another active session already assigned to this client.
```go
err := client.RestoreSession([]byte("yoursessionkeygoeshere"))
```

### Automatic Connection Maintenance
The WebWire client maintains the connection fully automatically to guarantee maximum connection uptime. It will automatically reconnect in the background whenever the connection is lost.

The only things to remember are:
- Client API methods such as `client.Request`, `client.TimedRequest` and `client.RestoreSession` will timeout if the server is unavailable for the entire duration of the specified timeout and thus the client fails to reconnect.
- `client.Signal` will immediately return a `DisconnectedErr` error if there's no connection at the time the signal was sent.

This feature is entirely optional and can be disabled at will which will cause `client.Request`, `client.TimedRequest` and `client.RestoreSession` to immediately return a `DisconnectedErr` error when there's no connection at the time the request is made.

### Thread Safety
It's safe to use both the session agents (those that are provided by the server through messages) and the client concurrently from multiple goroutines, the library automatically synchronizes concurrent operations.

### Hooks
Various hooks provide the ability to asynchronously react to different kinds of events and control the behavior of both the client and the server.

#### Server-side Hooks
- OnOptions
- BeforeUpgrade
- OnClientConnected
- OnClientDisconnected
- OnSignal
- OnRequest
- OnSessionKeyGeneration
- OnSessionCreated
- OnSessionLookup
- OnSessionClosed

#### Client-side Hooks
- OnServerSignal
- OnSessionCreated
- OnSessionClosed
- OnDisconnected

### Graceful Shutdown
The server will finish processing all ongoing signals and requests before closing when asked to shut down.
```go
// Will block the calling goroutine until all handlers have finished
server.Shutdown()
```
While the server is shutting down new connections are refused with `503 Service Unavailable` and incoming new requests from connected clients will be rejected with a special error: `RegErrSrvShutdown`. Any incoming signals from connected clients will be ignored during the shutdown.


### Seamless JavaScript Support
The [official JavaScript library](https://github.com/qbeon/webwire-js) enables seamless support for various JavaScript environments providing a fully compliant client implementation supporting the latest feature set of the [webwire-go](https://github.com/qbeon/webwire-go) library.

## Dependencies
This library depends on:
- **[websocket](https://github.com/gorilla/websocket)** version [v1.2.0](https://github.com/gorilla/websocket/releases/tag/v1.2.0) by **[Gorilla web toolkit](https://github.com/gorilla)** - A WebSocket implementation for Go.  
This library is used internally to abstract away the underlying websockets implementation.

----

© 2018 Roman Sharkov <roman.sharkov@qbeon.com>
