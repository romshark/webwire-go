package client

func extractMessageIdentifier(message []byte) (arr [8]byte) {
	copy(arr[:], message[1:9])
	return arr
}
