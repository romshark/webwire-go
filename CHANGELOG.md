# v1.0.0 - RC1

Released on 13th June 2018

This is the first release candidate version of the library. The API and the implementation are considered stable.

## Summary

Since the first beta-release many new requested features as well as stability and usability improvements have been introduces.
- The WebWire binary protocol has reached version 1.4 receiving many improvements to error reporting and bug-fixes.
- Modularity has been improved by moving responsibilities to separate implementable interfaces such as `SessionManager`, `ServerImplementation`, `SessionKeyGenerator`, `SessionInfo`, `SessionInfoParser` and `client.Implementation`.
- The session-key generator is now (optionally) replaceable. ([#10](https://github.com/qbeon/webwire-go/issues/10))
- A new default file-based session manager has been introduced, which can be replaced by any custom session manager implementation at will. ([#8](https://github.com/qbeon/webwire-go/issues/8))
- The server now shuts down gracefully waiting for active handlers to finish before terminating, gracefully rejecting incoming requests and dropping incoming signals. This means that signals are still guaranteed to arrive though not guaranteed to be processed while requests are guaranteed to either be processed or explicitly refused acknowledging the client. ([#9](https://github.com/qbeon/webwire-go/issues/9))
- Client agents are now closed gracefully waiting for all ongoing operations to finish before termination.
- Sessions can now (optionally) have multiple simultaneous connections which is especially useful for browser-based applications running multiple application instances in separate tabs each requiring a separate client. The maximum number of simultaneous connections per session is configurable. ([#4](https://github.com/qbeon/webwire-go/issues/4), [#20](https://github.com/qbeon/webwire-go/issues/20))
- The client implementation has received major improvements and now (optionally) maintains the connection automatically reconnecting when the connection is lost. ([#11](https://github.com/qbeon/webwire-go/issues/11), [#12](https://github.com/qbeon/webwire-go/issues/12))
- Support for [dep](https://github.com/golang/dep) has been added. Nevertheless, webwire remains `go get`'able thanks to vendor-embedded dependencies, which are shipped with the library though stripped when `dep` is used. ([#2](https://github.com/qbeon/webwire-go/issues/2), [216d118](https://github.com/qbeon/webwire-go/commit/216d1186526dc0e3d5b18f07985f8a01deef09df))
- Test coverage is now  monitored by coveralls. ([8582d04](https://github.com/qbeon/webwire-go/commit/8582d04efb176517022978b761b35c6b15973971), [161be37](https://github.com/qbeon/webwire-go/commit/161be373627ba45219527a850a8a231a5eed4761), [fd40a4e](https://github.com/qbeon/webwire-go/commit/fd40a4e0301fcb66e67575629b2a51cd20c1e4c9))
- Requests don't require a payload any longer. A request may now provide either only a name or only a payload or both (but not none of them). ([#22](https://github.com/qbeon/webwire-go/issues/22), [b498cca](https://github.com/qbeon/webwire-go/commit/b498cca0582bc8f197838931561daaa5969609f6))

## Changes

### Server

- The payload type now implements an encoding conversion method to convert it to a UTF8 encoded string. ([a05f330](https://github.com/qbeon/webwire-go/commit/a05f3305387e4711418171cc310cddfa1a977a4f))
- The standard go error interface is now used for errors returned in the request hook. Returning a non-`wwr.Error` error type will now return an empty `internal error` to the client and log it on the server to prevent leaks of sensitive information. ([dfdb7b3](https://github.com/qbeon/webwire-go/commit/dfdb7b3fac89932a5d33d3373f250dcdfa1b1ac4), [#PR 6](https://github.com/qbeon/webwire-go/pull/6))
- The client implementation now provides a new `SessionInfo` method to read specific session info fields. ([6541735](https://github.com/qbeon/webwire-go/commit/65417353dc28f0f1083f5a1642fdebfbeab1a609))
- Multi-Connection sessions have been introduces which allows clients to create multiple connections within a single session. This is especially useful for browsers where a single app could be opened in multiple browser-tabs each requiring one connection. ([144c6ad](https://github.com/qbeon/webwire-go/commit/144c6ad984f885c7939fb30cf3de8fb3dde5c908))
- The graceful shutdown feature has been introduces which prevents active handlers from being interrupted during server shutdown. ([678f860](https://github.com/qbeon/webwire-go/commit/678f8602c1efe3f48c29609bc228c9597474d44a))
- New session creation error types have been added: `ErrMaxSessConnsReached`, `SessNotFoundErr`, `ErrSessionsDisabled`. ([ac73199](https://github.com/qbeon/webwire-go/commit/ac73199a3328ee6445d759771273e7b0fe7b1811), [060c85b](https://github.com/qbeon/webwire-go/commit/060c85ba31d0d647f9cb291e79d11d15774534a6), [37e0a86](https://github.com/qbeon/webwire-go/commit/37e0a86c1ce4f2bcd8a9c137adb09ed40149d909))
- All signal and request handling related hooks have been moved to the new `ServerImplementation` interface, which must be implemented by the library user. ([4020691](https://github.com/qbeon/webwire-go/commit/4020691e6343493b3b24ca6e5ea9210a5b42e28b))
- The generation of the session keys have been moved to the new `SessionKeyGenerator` interface, which, by default, is implemented by the default session key generator. ([f292d1e](https://github.com/qbeon/webwire-go/commit/f292d1ea01d4d825a632a92f3de30b526389ffee))
- All session related hooks have been moved to the new `SessionManager` interface, which, by default, is implemented by the default file-based session manager. ([c5e515b](https://github.com/qbeon/webwire-go/commit/c5e515bbf83795720ab6ff6831617a1bb44bde1a))
- The underlying socket implementation is now abstracted away using the `Socket` interface, which is currently implemented using [gorilla/websocket](https://github.com/gorilla/websocket). ([2daff67](https://github.com/qbeon/webwire-go/commit/2daff6755037fc948806aa04ce0be8130bb7a8d5))
- The client agent and message objects are now directly passed to the `OnRequest` and `OnSignal` hooks instead of passing them through the context. ([393d513](https://github.com/qbeon/webwire-go/commit/393d51322b3eb9c7da67a0d47b802ed1070e6283))
- `log.Logger`is now used for logging instead of `io.Writer` to allow for better customization of logs. ([af872e5](https://github.com/qbeon/webwire-go/commit/af872e50925cf8d27ac987f25224f7a79f7cb616))
- The server is now defined as a `webwire.Server` interface type including new methods such as: `Run`, `ActiveSessionsNum`, `SessionConnectionsNum` and `CloseSession`. A webwire server instance is now created using the new constructor function. ([e1e9259](https://github.com/qbeon/webwire-go/commit/e1e925933de325a1ea31aae07666548019926ca0))
- The `SessionConnections` method will now return a list of client agents instead of just the number of concurrent connections of a session. ([8f789c0](https://github.com/qbeon/webwire-go/commit/8f789c07c2bd09fd6821bfcb62b7f429b364c045))
- Add a new client agent method `Close` for closing an ongoing
connection on the server side. ([e1e9259](https://github.com/qbeon/webwire-go/commit/e1e925933de325a1ea31aae07666548019926ca0))
- The method `ActiveSessions` was renamed to `ActiveSessionsNum`. ([e1e9259](https://github.com/qbeon/webwire-go/commit/e1e925933de325a1ea31aae07666548019926ca0))
- Deferred client agent shutdown has been implemented. The client agent now keeps track of all currently processed tasks and closes only when it's idle. ([d3b8313](https://github.com/qbeon/webwire-go/commit/d3b83133e0fd62318f9500c4b32220c0d7792d30))
- Support for external HTTP(S) servers has been added. The server interface now implements the standard Go HTTP handler interface and provides a new headless WebWire server constructor. ([8de01cc](https://github.com/qbeon/webwire-go/commit/8de01cc60ba89fc51e05c0b8539e912f2a31b306))
- Add a new `LastLookup` field to the session object denoting the last time a session was looked up to prevent active sessions from being garbage collected.
Also update the SessionManager interface and introduce a new
return type `SessionLookupResult` to ease use declaration of the
`OnSessionLookup` hook, which must now return a `SessNotFoundErr` in case the session wasn't found by key. ([cbd6017](https://github.com/qbeon/webwire-go/commit/cbd6017494ae5c423a7e419e2e650bd3d0eb81b3))
- A new protocol error message type has been introduced indicating
protocol violations when parsing a message.
The webwire server will now return a protocol error message if the determined message type requires a reply, otherwise it'll just drop the message. ([2316b59](https://github.com/qbeon/webwire-go/commit/2316b59f6cac7b01913c4d4ffff79e5298553f99))

### Client
- `client.Request` and `client.TimedRequest` will now return a special `ReqTimeoutErr` error type on timeouts. ([70edd63](https://github.com/qbeon/webwire-go/commit/70edd636294422b062448ac3c62239ee9aa8f4ab), [37e0a86](https://github.com/qbeon/webwire-go/commit/37e0a86c1ce4f2bcd8a9c137adb09ed40149d909))
- A new `ErrDisconnected` error type indicating connection loss has been introduced. ([132698f](https://github.com/qbeon/webwire-go/commit/132698f04aa40cfb6ab52b5462897e48283b8fc4))
- The client API now provides a new `IsConnected` method. ([a6ee217](https://github.com/qbeon/webwire-go/commit/a6ee217efdccd12d64614fd8bd8f7f3f7d8875f5))
- The client now maintains the connection and automatically reestablishes it when it's lost, unless explicitly instructed otherwise. During connection downtime outgoing requests are buffered and resent when the connection is reestablished, all calls that fire a request will block the calling goroutine until they either time out or succeed. ([aad55e7](https://github.com/qbeon/webwire-go/commit/aad55e73a0024ac20891b9c4ff6f4dfd5a45f6a4), [a7b0974](https://github.com/qbeon/webwire-go/commit/a7b0974be49b0aed1215b17755b47653b1925696))
- All hooks have been moved to the new `Implementation` interface, which must be implemented by the library user. ([6c4011d](https://github.com/qbeon/webwire-go/commit/6c4011dc41c12e9ad2320c63117c5ce925fb30da))
- A new `SessionInfo` interface and a `SessionInfoParser`
function type have been introduced to make session info objects immutable and thus data race save. ([2b6da64](https://github.com/qbeon/webwire-go/commit/2b6da6470b95dbe2aa7a5d4877d2768b450f3f2f))

###  Bug Fixes
- A data race caused by unsynchronized access to the client session object has been fixed. ([ccd0fa1](https://github.com/qbeon/webwire-go/commit/ccd0fa1aae8d6be28338e1963da270d4c9558a03))
- A protocol bug related to UTF16 encoded messages has been fixed. ([83a802b](https://github.com/qbeon/webwire-go/commit/83a802bbf2955cd37ad6a57db3647aac6b3e9b5e))
- Concurrent access to the client agent has been synchronized preventing data races. ([8ef8b85](https://github.com/qbeon/webwire-go/commit/8ef8b852333697c67b404dc1508d2457b0cfed17))
- `client.Close`will now block the calling goroutine until the socket-reader goroutine finally dies. ([85933a3](https://github.com/qbeon/webwire-go/commit/85933a39d6970a4145d58e75e4b0461c0f797044))
- Fixed security issue in the error-reply parser. ([c1adf42](https://github.com/qbeon/webwire-go/commit/c1adf429c29b4865e600949b2a29e96b21a177de))
- Fixed broken session restoration (attached data won't restore) by using the session info parser
during session deserialization (from a file or a database).
Also the session info parser is now used when receiving the session on the client. ([9b74216](https://github.com/qbeon/webwire-go/commit/9b742162c0eb6716dc7d379b0aaa524eac3e3f0f))
- Prevent sessions from being destroyed when the last active connection is closed. ([8b0dbd5](https://github.com/qbeon/webwire-go/commit/8b0dbd59e77af85cfaa99c9498fa6c13471af204))

### Tests

Both the acceptance and the unit tests received major updates, fixes and improvements. The total test coverage reached 76% and is now monitored by [coveralls.io](https://coveralls.io).

### Examples

All examples have been updated to use the new APIs.

### Documentation

The documentation received major updates, fixes and improvements.

## Special Thanks

Without the contributors this release may have never happened.

Special thanks to:
- [Daniil Trishkin - FromZeus](https://github.com/FromZeus)
- [Roman Rulkov - xrei](https://github.com/xrei)
- [Alexey Palazhchenko - AlekSi](https://github.com/AlekSi)
- [Ruslan Zarifov - Xeizzen](https://github.com/Xeizzen)
- [Daniel Sharkov - DanielSharkov](https://github.com/DanielSharkov)

----

# v1.0.0 - Beta 1

Released on 6th March 2018

This release is the first official beta version of the library. The API is generally stable though minor changes to it may happen during the beta-phase.

## Summary

- Since the last release major improvements to the protocol have been made, the WebWire protocol v1.1 is now completely binary.
- The API of the library has been improved to be more user-friendly and intuitive.
- The implementation of the sessions feature has been completed.
- The request and signal namespacing feature has been implemented.

## Changes

### Bug-Fixes

Several bugs as well as security and stability issues in the library and the test suite regarding thread safety, parser security and state synchronization reliability have been fixed since the last release.

- Automated tests of the library won't cause dead-locks in case of asynchronous misbehavior anymore.
- A stability issue of the synchronization of the session between the client and the server potentially causing desynchronization has been fixed.
- A security issue in the binary parser potentially causing a segmentation fault in case of malicious messages providing incorrect size flags has been fixed.
- Concurrent access to the client and the client agent won't cause data races anymore.
- A bug disabling automatic session synchronization in case of a manual session restoration has been fixed.
- A bug in the clients request registry causing a failed request to succeed after its failure has been fixed.
- The payload of a reply message is now optional and won't cause parser errors on the client.
- Session objects are now marked with the correct encoding (UTF8 instead of UTF16).

### Tests

Several automated tests have been added since the last release. Existing tests were updated to use the new API version and received major security updates.

### Examples

Existing examples have been updated to demonstrate the usage of the new API version.

### Documentation

General improvements to the documentation have been made. A visual representation of a subset of the new binary protocol has been added.

----

# v0.1.0 - Alpha 1
Released in 22nd February 2018

This release is the first official alpha version of the library. The API is yet unstable and may contain several bugs, though essential functionality has already been implemented.

## Changes

### Examples

This release contains 3 examples:

- **Echo** - a basic request-reply topology demonstration.
- **PubSub** - a basic publish-subscribe topology demonstration.
- **Chatroom** - an advanced demonstration of the capabilities of the library including req-rep, pub-sub, authentication and sessions.

### Documentation

Significant code parts are fully documented.
