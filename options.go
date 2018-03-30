package webwire

import (
	"io"
	"os"
)

// ServerOptions represents the options used during the creation of a new WebWire server instance
type ServerOptions struct {
	Hooks                 Hooks
	SessionsEnabled       bool
	SessionManager        SessionManager
	MaxSessionConnections uint
	WarnLog               io.Writer
	ErrorLog              io.Writer
}

// SetDefaults sets the defaults for undefined required values
func (srvOpt *ServerOptions) SetDefaults() {
	srvOpt.Hooks.SetDefaults()

	if srvOpt.SessionManager == nil {
		// Force the default session manager to use the default session directory
		srvOpt.SessionManager = NewDefaultSessionManager("")
	}

	if srvOpt.WarnLog == nil {
		srvOpt.WarnLog = os.Stdout
	}

	if srvOpt.ErrorLog == nil {
		srvOpt.ErrorLog = os.Stderr
	}
}

// HeadedServerOptions represents the options used during the creation of
// a new headed WebWire server instance
type HeadedServerOptions struct {
	ServerAddress string
	ServerOptions ServerOptions
}

// SetDefaults sets default values to undefined options
func (opts *HeadedServerOptions) SetDefaults() {
	opts.ServerOptions.SetDefaults()
}
