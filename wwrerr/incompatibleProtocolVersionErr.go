package wwrerr

import "fmt"

// IncompatibleProtocolVersionErr represents a connection error indicating that
// the server requires an incompatible version of the protocol and can't
// therefore be connected to
type IncompatibleProtocolVersionErr struct {
	RequiredVersion  string
	SupportedVersion string
}

// Error implements the error interface
func (err IncompatibleProtocolVersionErr) Error() string {
	return fmt.Sprintf(
		"unsupported protocol version: %s (supported: %s)",
		err.RequiredVersion,
		err.SupportedVersion,
	)
}
