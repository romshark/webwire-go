package webwire

import (
	"log"
	"os"
)

// ServerOptions represents the options used during the creation of a new WebWire server instance
type ServerOptions struct {
	SessionsEnabled       bool
	SessionManager        SessionManager
	SessionKeyGenerator   SessionKeyGenerator
	SessionInfoParser     SessionInfoParser
	MaxSessionConnections uint
	WarnLog               *log.Logger
	ErrorLog              *log.Logger
}

// SetDefaults sets the defaults for undefined required values
func (srvOpt *ServerOptions) SetDefaults() {
	if srvOpt.SessionsEnabled && srvOpt.SessionManager == nil {
		// Force the default session manager to use the default session directory
		srvOpt.SessionManager = NewDefaultSessionManager("")
	}

	if srvOpt.SessionsEnabled && srvOpt.SessionKeyGenerator == nil {
		srvOpt.SessionKeyGenerator = NewDefaultSessionKeyGenerator()
	}

	if srvOpt.SessionInfoParser == nil {
		srvOpt.SessionInfoParser = GenericSessionInfoParser
	}

	// Create default loggers to std-out/err when no loggers are specified
	if srvOpt.WarnLog == nil {
		srvOpt.WarnLog = log.New(
			os.Stdout,
			"WEBWIRE_WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}
	if srvOpt.ErrorLog == nil {
		srvOpt.ErrorLog = log.New(
			os.Stderr,
			"WEBWIRE_ERR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
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
