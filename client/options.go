package client

import (
	"io"
	"os"
	"time"
)

// OptionToggle represents the value of a togglable option
type OptionToggle int

const (
	// OptUnset defines unset togglable options
	OptUnset OptionToggle = iota

	// OptDisabled defines disabled togglable options
	OptDisabled

	// OptEnabled defines enabled togglable options
	OptEnabled
)

// Options represents the options used during the creation a new client instance
type Options struct {
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
	Autoconnect OptionToggle

	// ReconnectionInterval defines the interval at which autoconnect
	// should retry connection establishment.
	// If undefined then the default value of 2 seconds is applied
	ReconnectionInterval time.Duration

	// WarnLog defines the warn logging output target
	WarnLog io.Writer

	// ErrorLog defines the error logging output target
	ErrorLog io.Writer
}

// SetDefaults sets default values for undefined required options
func (opts *Options) SetDefaults() {
	if opts.DefaultRequestTimeout < 1 {
		opts.DefaultRequestTimeout = 60 * time.Second
	}

	if opts.Autoconnect == OptUnset {
		opts.Autoconnect = OptEnabled
	}

	if opts.ReconnectionInterval < 1 {
		opts.ReconnectionInterval = 2 * time.Second
	}

	if opts.WarnLog == nil {
		opts.WarnLog = os.Stdout
	}

	if opts.ErrorLog == nil {
		opts.ErrorLog = os.Stderr
	}
}
