<!-- HEADER -->
<h1 align="center">
	<br>
	<a href="https://github.com/qbeon/webwire-go"><img src="https://cdn.rawgit.com/qbeon/webwire-go/c7c2c74e/docs/img/webwire_logo.svg" alt="WebWire" width="256"></a>
	<br>
	<br>
	WebWire for <a href="https://golang.org/">Go</a>
	<br>
	<sub>An asynchronous duplex messaging library</sub>
</h1>
<p align="center">
	<a href="https://travis-ci.org/qbeon/webwire-go">
		<img src="https://travis-ci.org/qbeon/webwire-go.svg?branch=master" alt="Travis CI: build status">
	</a>
	<a href="https://coveralls.io/github/qbeon/webwire-go?branch=master">
		<img src="https://coveralls.io/repos/github/qbeon/webwire-go/badge.svg?branch=master" alt="Coveralls: Test Coverage">
	</a>
	<a href="https://goreportcard.com/report/github.com/qbeon/webwire-go">
		<img src="https://goreportcard.com/badge/github.com/qbeon/webwire-go" alt="GoReportCard">
	</a>
	<a href="https://codebeat.co/projects/github-com-qbeon-webwire-go-master">
		<img src="https://codebeat.co/badges/809181da-797c-4cdd-bb23-d0324935f3b0" alt="CodeBeat: Status">
	</a>
	<a href="https://codeclimate.com/github/qbeon/webwire-go/maintainability">
		<img src="https://api.codeclimate.com/v1/badges/243a45cacec7d850c64d/maintainability" alt="CodeClimate: Maintainability">
	</a>
	<br>
	<a href="https://opensource.org/licenses/MIT">
		<img src="https://img.shields.io/badge/License-MIT-green.svg" alt="Licence: MIT">
	</a>
	<a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fqbeon%2Fwebwire-go?ref=badge_shield" alt="FOSSA Status">
		<img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fqbeon%2Fwebwire-go.svg?type=shield"/>
	</a>
	<a href="https://godoc.org/github.com/qbeon/webwire-go">
		<img src="https://godoc.org/github.com/qbeon/webwire-go?status.svg" alt="GoDoc">
	</a>
</p>
<p align="center">
	<a href="https://opencollective.com/webwire">
		<img src="https://opencollective.com/webwire/tiers/backer.svg?avatarHeight=64" alt="OpenCollective">
	</a>
</p>
<br>

<!-- CONTENT -->
WebWire is a high-performance transport independent asynchronous duplex messaging library and an open source binary message protocol with builtin authentication and support for UTF8 and UTF16 encoding.
The [webwire-go](https://github.com/qbeon/webwire-go) library provides a server implementation for the Go programming language.
<br>

#### Table of Contents
- [Installation](#installation)
	- [Dep](#dep)
	- [Go Get](#go-get)
- [Contribution](#contribution)
	- [Maintainers](#maintainers)
- [WebWire Binary Protocol](#webwire-binary-protocol)
- [Examples](#examples)
- [Features](#features)
	- [Request-Reply](#request-reply)
	- [Client-side Signals](#client-side-signals)
	- [Server-side Signals](#server-side-signals)
	- [Namespaces](#namespaces)
	- [Sessions](#sessions)
	- [Concurrency](#concurrency)
	- [Hooks](#hooks)
		- [Server-side Hooks](#server-side-hooks)
		- [SessionManager Hooks](#sessionmanager-hooks)
		- [Client-side Hooks](#client-side-hooks)
		- [SessionKeyGenerator Hooks](#sessionkeygenerator-hooks)
	- [Graceful Shutdown](#graceful-shutdown)
	- [Multi-Language Support](#multi-language-support)
	- [Security](#security)


## Installation
Choose any stable release from [the available release tags](https://github.com/qbeon/webwire-go/releases) and copy the source code into your project's vendor directory: ```$YOURPROJECT/vendor/github.com/qbeon/webwire-go```. All necessary transitive [dependencies](https://github.com/qbeon/webwire-go#dependencies) are already embedded into the `webwire-go` repository.

### Dep
If you're using [dep](https://github.com/golang/dep), just use [dep ensure](https://golang.github.io/dep/docs/daily-dep.html#adding-a-new-dependency) to add a specific version of webwire-go including all its transitive dependencies to your project: ```dep ensure -add github.com/qbeon/webwire-go@v1.0.0-rc1```. This will remove all embedded transitive dependencies and move them to your projects `vendor` directory.

### Go Get
You can also use [go get](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies): ```go get github.com/qbeon/webwire-go``` but beware that this will fetch the latest commit of the [master branch](https://github.com/qbeon/webwire-go/commits/master) which is currently **not yet** considered a stable release branch. It's therefore recommended to use [dep](https://github.com/qbeon/webwire-go#dep) instead.

## Contribution
Contribution of any kind is always welcome and appreciated, check out our [Contribution Guidelines](https://github.com/qbeon/webwire-go/blob/master/CONTRIBUTING.md) for more information!

### Maintainers
| Maintainer | Role | Specialization |
| :--- | :--- | :--- |
| **[Roman Sharkov](https://github.com/romshark)** | Core Maintainer | Dev (Go, JavaScript) |
| **[Daniil Trishkin](https://github.com/FromZeus)** | CI Maintainer | DevOps |

## WebWire Binary Protocol
WebWire is built for speed and portability implementing an open source [binary protocol](https://github.com/qbeon/webwire-go/blob/master/docs/protocol-sequences.svg).
![Protocol Subset Diagram](https://github.com/qbeon/webwire-go/blob/master/docs/img/wwr_msgproto_diagram.svg)

The first byte defines the [type of the message](https://github.com/qbeon/webwire-go/blob/master/message/message.go#L91). Requests and replies contain an incremental 8-byte identifier that must be unique in the context of the senders' session. A 0 to 255 bytes long 7-bit ASCII encoded name is contained in the header of a signal or request message.
A header-padding byte is applied in case of UTF16 payload encoding to properly align the payload sequence.
Fraudulent messages are recognized by analyzing the message length, out-of-range memory access attacks are therefore prevented.

## Examples
- **[Echo](https://github.com/qbeon/webwire-go-examples/tree/master/echo)** - Demonstrates a simple request-reply implementation using the [Go client](https://github.com/qbeon/webwire-go-client).

- **[Pub-Sub](https://github.com/qbeon/webwire-go-examples/tree/master/pubsub)** - Demonstrantes a simple publisher-subscriber tolopology using the [Go client](https://github.com/qbeon/webwire-go-client).

- **[Chat Room](https://github.com/qbeon/webwire-go-examples/tree/master/chatroom)** - Demonstrates advanced use of the library. The corresponding [JavaScript Chat Room Client](https://github.com/qbeon/webwire-js/tree/master/examples/chatroom-client-vue) is implemented with the [Vue.js framework](https://vuejs.org/).

## Features
### Request-Reply
Clients can initiate multiple simultaneous requests and receive replies asynchronously. Requests are multiplexed through the connection similar to HTTP2 pipelining. The below examples are using the [webwire Go client](https://github.com/qbeon-webwire-go-client).

```go
// Send a request to the server,
// this will block the goroutine until either a reply is received
// or the default timeout triggers (if there is one)
reply, err := client.Request(
	context.Background(), // No cancelation, default timeout
	nil,                  // No name
	wwr.Payload{
		Data: []byte("sudo rm -rf /"), // Binary request payload
	},
)
defer reply.Close() // Close the reply
if err != nil {
	// Oh oh, the request failed for some reason!
}
reply.PayloadUtf8() // Here we go!
 ```

Requests will respect cancelable contexts and deadlines

```go
cancelableCtx, cancel := context.WithCancel(context.Background())
defer cancel()
timedCtx, cancelTimed := context.WithTimeout(cancelableCtx, 1*time.Second)
defer cancelTimed()

// Send a cancelable request to the server with a 1 second deadline
// will block the goroutine for 1 second at max
reply, err := client.Request(timedCtx, nil, wwr.Payload{
	Encoding: wwr.EncodingUtf8,
	Data:     []byte("hurry up!"),
})
defer reply.Close()

// Investigate errors manually...
switch err.(type) {
case wwr.ErrCanceled:
	// Request was prematurely canceled by the sender
case wwr.ErrDeadlineExceeded:
	// Request timed out, server didn't manage to reply
	// within the user-specified context deadline
case wwr.TimeoutErr:
	// Request timed out, server didn't manage to reply
	// within the specified default request timeout duration
case nil:
	// Replied successfully
}

// ... or check for a timeout error the easier way:
if err != nil {
	if wwr.IsErrTimeout(err) {
		// Timed out due to deadline excess or default timeout
	} else {
		// Unexpected error
	}
}

reply // Just in time!
```

### Client-side Signals
Individual clients can send signals to the server. Signals are one-way messages guaranteed to arrive, though they're not guaranteed to be processed like requests are. In cases such as when the server is being shut down, incoming signals are ignored by the server and dropped while requests will acknowledge the failure. The below examples are using the [webwire Go client](https://github.com/qbeon-webwire-go-client).

```go
// Send signal to server
err := client.Signal(
	[]byte("eventA"),
	wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("something"),
	},
)
```

### Server-side Signals
The server also can send signals to individual connected clients.

```go
func OnRequest(
  _ context.Context,
  conn wwr.Connection,
  _ wwr.Message,
) (wwr.Payload, error) {
	// Send a signal to the client before replying to the request
	conn.Signal(
		nil, // No message name
		wwr.Payload{
			Encoding: wwr.EncodingUtf8,
			Data:     []byte("example")),
		},
	)

	// Reply nothing
	return wwr.Payload{}, nil
}
```

### Namespaces
Different kinds of requests and signals can be differentiated using the builtin namespacing feature.

```go
func OnRequest(
	_ context.Context,
	_ wwr.Connection,
	msg wwr.Message,
) (wwr.Payload, error) {
	switch msg.Name() {
	case "auth":
		// Authentication request
		return wwr.Payload{
			Encoding: wwr.EncodingUtf8,
      			Data:     []byte("this is an auth request"),
		}
	case "query":
		// Query request
		return wwr.Payload{
			Encoding: wwr.EncodingUtf8,
			Data:     []byte("this is a query request"),
		}
	}

	// Otherwise return nothing
	return wwr.Payload{}, nil
}
```
```go
func OnSignal(
	_ context.Context,
	_ wwr.Connection,
	msg wwr.Message,
) {
	switch string(msg.Name()) {
	case "event A":
		// handle event A
	case "event B":
		// handle event B
	}
}
```

### Sessions
Individual connections can get sessions assigned to identify them. The state of the session is automagically synchronized between the client and the server. WebWire doesn't enforce any kind of authentication technique though, it just provides a way to authenticate a connection. WebWire also doesn't enforce any kind of session storage, the user could implement a custom session manager implementing the WebWire `SessionManager` interface to use any kind of volatile or persistent session storage, be it a database or a simple in-memory map.

```go
func OnRequest(
	_ context.Context,
	conn wwr.Connection,
	msg wwr.Message,
) (wwr.Payload, error) {
	// Verify credentials
	if string(msg.Payload()) != "secret:pass" {
		return wwr.Payload{}, wwr.ReqErr {
			Code:    "WRONG_CREDENTIALS",
			Message: "Incorrect username or password, try again",
		}
	}
	// Create session (will automatically synchronize to the client)
	err := conn.CreateSession(/*something that implements wwr.SessionInfo*/)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create session for some reason")
	}

	// Complete request, reply nothing
	return wwr.Payload{}, nil
}
```

WebWire provides a basic file-based session manager implementation out of the box used by default when no custom session manager is defined. The default session manager creates a file with a .wwrsess extension for each opened session in the configured directory (which, by default, is the directory of the executable). During the restoration of a session the file is looked up by name using the session key, read and unmarshalled recreating the session object.

### Concurrency
Messages are parsed and handled concurrently in a separate goroutine by default. The total number of concurrently executed handlers can be independently throttled down for each individual connection, which is unlimited by default.

All exported interfaces provided by both the server and the client are thread safe and can thus safely be used concurrently from within multiple goroutines, the library automatically synchronizes all concurrent operations.

### Hooks
Various hooks provide the ability to asynchronously react to different kinds of events and control the behavior of both the client and the server.

#### Server-side Hooks
- OnClientConnected
- OnClientDisconnected
- OnSignal
- OnRequest

#### SessionManager Hooks
- OnSessionCreated
- OnSessionLookup
- OnSessionClosed

#### Client-side Hooks
- OnServerSignal
- OnSessionCreated
- OnSessionClosed
- OnDisconnected

#### SessionKeyGenerator Hooks
- Generate

### Graceful Shutdown
The server will finish processing all ongoing signals and requests before closing when asked to shut down.
```go
// Will block until all handlers have finished
server.Shutdown()
```
While the server is shutting down new connections are refused with `503 Service Unavailable` and incoming new requests from connected clients will be rejected with a special error: `RegErrSrvShutdown`. Any incoming signals from connected clients will be ignored during the shutdown.

Server-side client connections also support graceful shutdown, a connection will be closed when all work on it is done,
while incoming requests and signals are handled similarly to shutting down the server.
```go
// Will block until all work on this connection is done
connection.Close()
```

### Multi-Language Support
The following libraries provide seamless support for various development environments providing fully compliant protocol implementations supporting the latest features.
- **Go (server & client)**: An [official Go client](https://github.com/qbeon/webwire-go-client) implementation is available.
- **JavaScript (client)**: An [official JavaScript library](https://github.com/qbeon/webwire-js) enables seamless support for various JavaScript environments ([93% of web-browsers](https://caniuse.com/#search=websockets) & [Node.js](https://nodejs.org/en/)) providing a fully compliant client implementation (requires a websocket-based transport implementation such as [qbeon/webwire-go-gorilla](https://github.com/qbeon/webwire-go-gorilla) or [qbeon/webwire-go-fasthttp](https://github.com/qbeon/webwire-go-fasthttp)).

### Security
A webwire server can be hosted by a [TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) protected server transport implementation to prevent [man-in-the-middle attacks](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) as well as to verify the identity of the server during connection establishment. Setting up a TLS protected websocket server for example is easy:
```go
// Setup a secure webwire server instance
server, err := wwr.NewServer(
	serverImplementation,
	wwr.ServerOptions{
		Host: "localhost:443",
	},
	// Use a TLS protected transport layer
	&wwrgorilla.Transport{
		TLS: &wwrgorilla.TLS{
			// Provide key and certificate
			CertFilePath:       "path/to/certificate.crt",
			PrivateKeyFilePath: "path/to/private.key",
			// Specify TLS configs
			Config: &tls.Config{
				MinVersion:               tls.VersionTLS12,
				CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				},
			},
		},
	},
)
if err != nil {
	panic(fmt.Errorf("failed setting up wwr server: %s", err))
}
// Launch
if err := server.Run(); err != nil {
	panic(fmt.Errorf("wwr server failed: %s", err))
}
```
The above code example is using the [webwire-go-gorilla](https://github.com/qbeon/webwire-go-gorilla) transport implementation.

----

© 2018 Roman Sharkov <roman.sharkov@qbeon.com>
