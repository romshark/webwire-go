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
	// Hooks define the callback hook functions provided by the user to define behavior
	// on certain events
	Hooks Hooks

	// DefaultRequestTimeout defines the default request timeout duration used in client.Request
	DefaultRequestTimeout time.Duration

	// ReconnectionInterval defines the interval at which autoconnect should poll for a connection.
	// If undefined then the default value of 2 seconds is applied
	ReconnectionInterval time.Duration

	// If autoconnect is enabled, client.Request, client.TimedRequest and client.RestoreSession
	// won't immediately return a disconnected error if there's no active connection to the server,
	// instead they will automatically try to reestablish the connection
	// before the timeout is triggered and a timeout error is returned.
	// Autoconnect is enabled by default
	Autoconnect OptionToggle
	WarnLog     io.Writer
	ErrorLog    io.Writer
}

// SetDefaults sets default values for undefined required options
func (opts *Options) SetDefaults() {
	opts.Hooks.SetDefaults()

	if opts.Autoconnect == OptUnset {
		opts.Autoconnect = OptEnabled
	}

	if opts.DefaultRequestTimeout < 1 {
		opts.DefaultRequestTimeout = 60 * time.Second
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
