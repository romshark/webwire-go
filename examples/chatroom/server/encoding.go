package main

// http://play.golang.org/p/fVf7duRtdH

import (
	"bytes"
	"fmt"
	"unicode/utf16"
	"unicode/utf8"

	wwr "github.com/qbeon/webwire-go"
)

// ToUTF8 was borrowed from https://gist.github.com/bradleypeabody/185b1d7ed6c0c2ab6cec
func ToUTF8(payload wwr.Payload) (string, error) {
	data := payload.Data

	if payload.Encoding == wwr.EncodingBinary {
		return "", fmt.Errorf("Cannot convert binary payload to UTF8")
	}

	if payload.Encoding == wwr.EncodingUtf8 {
		return string(data), nil
	}

	// Convert from UTF16
	if len(data)%2 != 0 {
		return "", fmt.Errorf("Unaligned (%d) payload data", len(data))
	}

	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)
	lb := len(data)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(data[i]) + (uint16(data[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}
	return ret.String(), nil
}
