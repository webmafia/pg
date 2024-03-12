package pg

import (
	"testing"

	"github.com/webmafia/fast"
)

func BenchmarkTableEncode(b *testing.B) {
	t := Table("foobar")
	buf := fast.NewStringBuffer(512)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.EncodeString(buf)
		buf.Reset()
	}
}
