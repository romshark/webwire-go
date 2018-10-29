package client

// verifyProtocolVersion returns true if the given version of the webwire
// protocol is acceptable for this client
func verifyProtocolVersion(major, minor byte) bool {
	if major != 2 {
		return false
	}
	return true
}
