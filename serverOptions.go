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
	ReadTimeout           time.Duration
	WarnLog               *log.Logger
	ErrorLog              *log.Logger
	MaxConnsPerIP         uint

	// ReadBufferSize defines the size of the inbound buffer in bytes
	ReadBufferSize uint32

	// WriteBufferSize defines the size of the outbound buffer in bytes
	WriteBufferSize uint32

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
			"WEBWIRE_WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}
	if op.ErrorLog == nil {
		op.ErrorLog = log.New(
			os.Stderr,
			"WEBWIRE_ERR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}

	const minBufferSize = 1024

	// Verify buffer sizes
	if op.ReadBufferSize == 0 {
		op.ReadBufferSize = minBufferSize
	} else if op.ReadBufferSize < minBufferSize {
		return fmt.Errorf(
			"read buffer size too small: %d bytes (min: %d bytes)",
			op.ReadBufferSize,
			minBufferSize,
		)
	}

	if op.WriteBufferSize == 0 {
		op.WriteBufferSize = minBufferSize
	} else if op.WriteBufferSize < minBufferSize {
		return fmt.Errorf(
			"write buffer size too small: %d bytes (min: %d bytes)",
			op.WriteBufferSize,
			minBufferSize,
		)
	}

	// Verify the message buffer size
	if op.MessageBufferSize == 0 {
		op.MessageBufferSize = op.ReadBufferSize
	} else if op.MessageBufferSize < op.ReadBufferSize {
		return fmt.Errorf(
			"message buffer size too small: %d bytes (min: %d bytes)",
			op.MessageBufferSize,
			op.ReadBufferSize,
		)
	}

	return nil
}
