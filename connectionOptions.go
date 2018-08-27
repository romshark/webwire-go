package webwire

const (
	// UnlimitedConcurrency represents an option for an unlimited
	// number of concurrent handlers of a connection
	UnlimitedConcurrency uint = 0
)

// connectionOptions represents an implementation
// of the ConnectionOptions interface
type connectionOptions struct {
	accept           bool
	concurrencyLimit uint
}

// Accept implements the ConnectionOptions interface
func (conopts *connectionOptions) Accept() bool {
	return conopts.accept
}

// ConcurrencyLimit implements the ConnectionOptions interface
func (conopts *connectionOptions) ConcurrencyLimit() uint {
	return conopts.concurrencyLimit
}

// AcceptConnection accepts an incoming connection using the given configuration
func AcceptConnection(concurrencyLimit uint) ConnectionOptions {
	return &connectionOptions{
		accept:           true,
		concurrencyLimit: concurrencyLimit,
	}
}

// RefuseConnection refuses an incoming connection using the given configuration
func RefuseConnection(reason string) ConnectionOptions {
	return &connectionOptions{
		accept: false,
	}
}
