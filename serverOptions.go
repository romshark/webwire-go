package webwire

import (
	"fmt"
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
	WarnLog               *log.Logger
	ErrorLog              *log.Logger
	ReadTimeout           time.Duration

	// MessageBufferSize defines the size of the message buffer
	MessageBufferSize uint32
}

// Prepare verifies the specified options and sets the default values to
// unspecified options
func (op *ServerOptions) Prepare() error {
	// Enable sessions by default
	if op.Sessions == OptionUnset {
		op.Sessions = Enabled
	}

	if op.Sessions == Enabled && op.SessionManager == nil {
		// Force the default session manager
		// to use the default session directory
		op.SessionManager = NewDefaultSessionManager("")
	}

	if op.Sessions == Enabled && op.SessionKeyGenerator == nil {
		op.SessionKeyGenerator = NewDefaultSessionKeyGenerator()
	}

	if op.SessionInfoParser == nil {
		op.SessionInfoParser = GenericSessionInfoParser
	}

	if op.ReadTimeout < 1*time.Second {
		op.ReadTimeout = 60 * time.Second
	}

	// Create default loggers to std-out/err when no loggers are specified
	if op.WarnLog == nil {
		op.WarnLog = log.New(
			os.Stdout,
			"WWR_WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}
	if op.ErrorLog == nil {
		op.ErrorLog = log.New(
			os.Stderr,
			"WWR_ERR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}

	const minMsgBufferSize = 32

	// Verify the message buffer size
	if op.MessageBufferSize == 0 {
		op.MessageBufferSize = minMsgBufferSize
	} else if op.MessageBufferSize < minMsgBufferSize {
		return fmt.Errorf(
			"message buffer size too small: %d bytes (min: %d bytes)",
			op.MessageBufferSize,
			minMsgBufferSize,
		)
	}

	return nil
}
