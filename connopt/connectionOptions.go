package connopt

// ConnectionAcceptance defines whether a connection is to be accepted
type ConnectionAcceptance byte

const (
	// Accept instructs the server to accept the incoming connection
	Accept ConnectionAcceptance = iota

	// Refuse instructs the server to refuse the incoming connection
	Refuse
)

// ConnectionOptions represents the options applied to an individual connection
// during accept
type ConnectionOptions struct {
	// Connection refuses the incoming connection when explicitly set to
	// wwr.Refuse. It's set to wwr.Accept by default.
	Connection ConnectionAcceptance

	// ConcurrencyLimit defines the maximum number of operations to be processed
	// concurrently for this particular client connection. If ConcurrencyLimit
	// is 0 (which it is by default) then the number of concurrent operations
	// for this particular connection will be limited to 1. Anything below 0
	// will lift the limitation entirely while everything above 0 will set the
	// limit to the specified number of handlers
	ConcurrencyLimit int
}
