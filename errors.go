package webwire

// ReqErrSrvShutdown represents a request error type indicating that the request cannot be
// processed due to the server currently being shut down
type ReqErrSrvShutdown struct{}

func (err ReqErrSrvShutdown) Error() string {
	return "Server is currently being shut down and won't process the request"
}

// ReqErr represents an error returned in case of a request that couldn't be processed
type ReqErr struct {
	Code    string `json:"c"`
	Message string `json:"m,omitempty"`
}

func (err ReqErr) Error() string {
	return err.Message
}
