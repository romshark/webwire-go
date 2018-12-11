package message

import pld "github.com/qbeon/webwire-go/payload"

// CalcMsgLenSignal returns the size of a signal message with the given name and
// payload
func CalcMsgLenSignal(
	name []byte,
	encoding pld.Encoding,
	payload []byte,
) int {
	if encoding == pld.Utf16 && len(name)%2 != 0 {
		return 3 + len(name) + len(payload)
	}
	return 2 + len(name) + len(payload)
}

// CalcMsgLenRequest returns the size of a request message with the given name
// and payload
func CalcMsgLenRequest(
	name []byte,
	encoding pld.Encoding,
	payload []byte,
) int {
	if encoding == pld.Utf16 && len(name)%2 != 0 {
		return 11 + len(name) + len(payload)
	}
	return 10 + len(name) + len(payload)
}
