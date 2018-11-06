package wwrerr

import "fmt"

// TransmissionErr represents a connection error indicating a failed
// transmission
type TransmissionErr struct {
	Cause error
}

// Error implements the error interface
func (err TransmissionErr) Error() string {
	return fmt.Sprintf("message transmission failed: %s", err.Cause)
}
