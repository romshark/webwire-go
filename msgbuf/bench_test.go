package msgbuf_test

import (
	"testing"

	"github.com/qbeon/webwire-go/msgbuf"
)

var messageBufferSize uint32 = 1024

// BenchmarkFastPoolConc tests the fast pool using 16 concurrently
func BenchmarkFastPoolConc(b *testing.B) {
	var pool msgbuf.Pool = msgbuf.NewFastPool(messageBufferSize, 0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.Get()
	}
}

// BenchmarkSyncPoolConc tests the sync.Pool based pool concurrently
func BenchmarkSyncPoolConc(b *testing.B) {
	var pool msgbuf.Pool = msgbuf.NewSyncPool(messageBufferSize, 0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.Get()
	}
}
