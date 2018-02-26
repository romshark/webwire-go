package client

func extractMessageIdentifier(message []byte) (arr [32]byte) {
	copy(arr[:], message[1:33])
	return arr
}
