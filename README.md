# WebWire for Golang
WebWire is an asynchronous duplex messaging library built on top of [WebSockets](https://developer.mozilla.org/de/docs/WebSockets).
The [webwire-go](https://github.com/qbeon/webwire-go) library provides both a client and a server implementation for the Go programming language.

## Features
WebWire provides a compact set of useful features that are not available and/or relatively hard to implement on raw WebSockets.
- **Request-Reply**, unlike with raw websockets clients can initiate multiple requests and get receive replies.
- **Server-side Signals**, the WebWire can send signals to individual connected clients.
- **Client-side Signals**, the individual clients can send signals to the server.
- **Authentication & Sessions**, both the client and server can initiate authentications and create sessions.
- **Thread Safety**, it's safe to send messages to clients from multiple goroutines simultaneously.
- **Hooks**, multiple hooks provide the ability to asynchronously react on different kinds of events.
