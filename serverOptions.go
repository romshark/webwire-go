package webwire

import (
	"log"
	"os"
	"time"
)

// OptionValue represents the setting value of an option
type OptionValue = int32

const (
	// OptionUnset represents the default unset value
	OptionUnset OptionValue = iota

	// Disabled disables an option
	Disabled

	// Enabled enables an option
	Enabled
)

// ServerOptions represents the options
// used during the creation of a new WebWire server instance
type ServerOptions struct {
	Host                  string
	Sessions              OptionValue
	SessionManager        SessionManager
	SessionKeyGenerator   SessionKeyGenerator
	SessionInfoParser     SessionInfoParser
	MaxSessionConnections uint
	ReadTimeout           time.Duration
	WarnLog               *log.Logger
	ErrorLog              *log.Logger
}

// SetDefaults sets the defaults for undefined required values
func (srvOpt *ServerOptions) SetDefaults() {
	// Enable sessions by default
	if srvOpt.Sessions == OptionUnset {
		srvOpt.Sessions = Enabled
	}

	if srvOpt.Sessions == Enabled && srvOpt.SessionManager == nil {
		// Force the default session manager
		// to use the default session directory
		srvOpt.SessionManager = NewDefaultSessionManager("")
	}

	if srvOpt.Sessions == Enabled && srvOpt.SessionKeyGenerator == nil {
		srvOpt.SessionKeyGenerator = NewDefaultSessionKeyGenerator()
	}

	if srvOpt.SessionInfoParser == nil {
		srvOpt.SessionInfoParser = GenericSessionInfoParser
	}

	if srvOpt.ReadTimeout < 1*time.Second {
		srvOpt.ReadTimeout = 60 * time.Second
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
