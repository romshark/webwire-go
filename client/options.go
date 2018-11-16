package client

import (
	"fmt"
	"log"
	"os"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

// Options represents the options used during the creation a new client instance
type Options struct {
	// SessionInfoParser defines the optional session info parser function
	SessionInfoParser webwire.SessionInfoParser

	// DialingTimeout defines the dialing timeout
	DialingTimeout time.Duration

	// DefaultRequestTimeout defines the default request timeout duration
	// used by client.Request and client.RestoreSession
	DefaultRequestTimeout time.Duration

	// Autoconnect defines whether the autoconnect feature is to be enabled.
	//
	// If autoconnect is enabled then client.Request and client.RestoreSession
	// won't immediately return a disconnected error
	// if there's no active connection to the server,
	// instead they will automatically try to reestablish the connection
	// before the timeout is triggered and a timeout error is returned.
	//
	// Autoconnect is enabled by default
	Autoconnect webwire.OptionValue

	// ReconnectionInterval defines the interval at which autoconnect
	// should retry connection establishment.
	// If undefined then the default value of 2 seconds is applied
	ReconnectionInterval time.Duration

	// WarnLog defines the warn logging output target
	WarnLog *log.Logger

	// ErrorLog defines the error logging output target
	ErrorLog *log.Logger

	// MessageBufferSize defines the size of the inbound message buffer
	MessageBufferSize uint32
}

// Prepare validates the specified options and sets the default values for
// unspecified options
func (op *Options) Prepare() error {
	if op.SessionInfoParser == nil {
		op.SessionInfoParser = webwire.GenericSessionInfoParser
	}

	if op.DefaultRequestTimeout < 1 {
		op.DefaultRequestTimeout = 60 * time.Second
	}

	if op.Autoconnect == webwire.OptionUnset {
		op.Autoconnect = webwire.Enabled
	}

	if op.ReconnectionInterval < 1 {
		op.ReconnectionInterval = 2 * time.Second
	}

	// Create default loggers to std-out/err when no loggers are specified
	if op.WarnLog == nil {
		op.WarnLog = log.New(
			os.Stdout,
			"WEBWIRE_CLT_WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}
	if op.ErrorLog == nil {
		op.ErrorLog = log.New(
			os.Stderr,
			"WEBWIRE_CLT_ERR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}

	// Set default dialing timeout
	if op.DialingTimeout < 1 {
		op.DialingTimeout = 5 * time.Second
	}

	// Verify buffer sizes
	const minBufferSize = 1024

	// Verify the message buffer size
	if op.MessageBufferSize == 0 {
		op.MessageBufferSize = minBufferSize
	} else if op.MessageBufferSize < minBufferSize {
		return fmt.Errorf(
			"message buffer size too small: %d bytes (min: %d bytes)",
			op.MessageBufferSize,
			minBufferSize,
		)
	}

	return nil
}
