package client

import (
	"log"
	"os"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

// Options represents the options used during the creation a new client instance
type Options struct {
	// SessionInfoParser defines the optional session info parser function
	SessionInfoParser webwire.SessionInfoParser

	// DefaultRequestTimeout defines the default request timeout duration
	// used by client.Request and client.RestoreSession
	DefaultRequestTimeout time.Duration

	// Autoconnect defines whether the autoconnect feature is to be enabled.
	//
	// If autoconnect is enabled then client.Request, client.TimedRequest and
	// client.RestoreSession won't immediately return a disconnected error
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
}

// SetDefaults sets default values for undefined required options
func (opts *Options) SetDefaults() {
	if opts.SessionInfoParser == nil {
		opts.SessionInfoParser = webwire.GenericSessionInfoParser
	}

	if opts.DefaultRequestTimeout < 1 {
		opts.DefaultRequestTimeout = 60 * time.Second
	}

	if opts.Autoconnect == webwire.OptionUnset {
		opts.Autoconnect = webwire.Enabled
	}

	if opts.ReconnectionInterval < 1 {
		opts.ReconnectionInterval = 2 * time.Second
	}

	// Create default loggers to std-out/err when no loggers are specified
	if opts.WarnLog == nil {
		opts.WarnLog = log.New(
			os.Stdout,
			"WEBWIRE_CLT_WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}
	if opts.ErrorLog == nil {
		opts.ErrorLog = log.New(
			os.Stderr,
			"WEBWIRE_CLT_ERR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}
}
