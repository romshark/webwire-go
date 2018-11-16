package message

import "testing"

// BenchmarkParseRequestBinary benchmarks parsing of
// binary encoded request messages
func BenchmarkParseRequestBinary(b *testing.B) {
	// Generate a random request message
	// with 1 KiB (binary) payload
	// and a random name
	encoded, _, _, _ := rndRequestMsg(
		MsgRequestBinary,
		1, 255,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseRequestUtf8 benchmarks parsing of
// UTF8 encoded request messages
func BenchmarkParseRequestUtf8(b *testing.B) {
	// Generate a random request message
	// with 1 KiB (UTF8) payload
	// and a random name
	encoded, _, _, _ := rndRequestMsg(
		MsgRequestUtf8,
		1, 255,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseRequestUtf16 benchmarks parsing of
// UTF16 encoded request messages
func BenchmarkParseRequestUtf16(b *testing.B) {
	// Generate a random request message
	// with 1 KiB (UTF16) payload
	// and a random name
	encoded, _, _, _ := rndRequestMsgUtf16(
		1, 255,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseRequestBinaryLarge benchmarks parsing of
// large binary encoded request messages
func BenchmarkParseRequestBinaryLarge(b *testing.B) {
	// Generate a random request message
	// with 1 MiB (binary) payload
	// and a random name of maximum possible length
	encoded, _, _, _ := rndRequestMsg(
		MsgRequestBinary,
		255, 255,
		134217728, 134217728,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseReplyBinary benchmarks parsing of
// binary encoded reply messages
func BenchmarkParseReplyBinary(b *testing.B) {
	// Generate a random reply message
	// with 1 KiB (binary) payload
	// and a random name
	encoded, _, _ := rndReplyMsg(
		MsgReplyBinary,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseReplyUtf8 benchmarks parsing of
// UTF8 encoded reply messages
func BenchmarkParseReplyUtf8(b *testing.B) {
	// Generate a random reply message
	// with 1 KiB (UTF8) payload
	// and a random name
	encoded, _, _ := rndReplyMsg(
		MsgReplyUtf8,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseReplyUtf16 benchmarks parsing of
// UTF16 encoded reply messages
func BenchmarkParseReplyUtf16(b *testing.B) {
	// Generate a random reply message
	// with 1 KiB payload (UTF16)
	// and a random name
	encoded, _, _ := rndReplyMsgUtf16(
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseReplyBinaryLarge benchmarks parsing of
// large binary encoded request messages
func BenchmarkParseReplyBinaryLarge(b *testing.B) {
	// Generate a random reply message
	// with 1 MiB (binary) payload
	// and a random name of maximum possible length
	encoded, _, _ := rndReplyMsg(
		MsgReplyBinary,
		134217728, 134217728,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseSignalBinary benchmarks parsing of
// binary encoded signal messages
func BenchmarkParseSignalBinary(b *testing.B) {
	// Generate a random signal message
	// with 1 KiB (binary) payload
	// and a random name
	encoded, _, _ := rndSignalMsg(
		MsgSignalBinary,
		1, 255,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseSignalUtf8 benchmarks parsing of
// UTF8 encoded signal messages
func BenchmarkParseSignalUtf8(b *testing.B) {
	// Generate a random signal message
	// with 1 KiB (UTF8) payload
	// and a random name
	encoded, _, _ := rndSignalMsg(
		MsgSignalUtf8,
		1, 255,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseSignalUtf16 benchmarks parsing of
// UTF16 encoded signal messages
func BenchmarkParseSignalUtf16(b *testing.B) {
	// Generate a random signal message
	// with 1 KiB (UTF8) payload
	// and a random name
	encoded, _, _ := rndSignalMsgUtf16(
		1, 255,
		1024, 1024,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}

// BenchmarkParseSignalBinaryLarge benchmarks parsing of
// binary encoded signal messages
func BenchmarkParseSignalBinaryLarge(b *testing.B) {
	// Generate a random signal message
	// with 128 MiB (binary) payload
	// and a random name of maximum possible length
	encoded, _, _ := rndSignalMsg(
		MsgSignalBinary,
		255, 255,
		134217728, 134217728,
	)
	msg := NewMessage(uint32(len(encoded)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msg.ReadBytes(encoded); err != nil {
			b.Fatalf("Failed parsing: %s", err)
		}
	}
}
