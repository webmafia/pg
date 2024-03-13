package pg

import (
	"fmt"
	"testing"

	"github.com/webmafia/fast"
)

func Example_encodeQuery() {
	buf := fast.NewStringBuffer(256)
	encodeQuery(buf, "SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s", []any{Table("trudeluttan"), 123, Table("mjau"), 456})

	fmt.Printf("%#v", buf.String())

	// Output: Mjau
}

func Benchmark_encodeQuery(b *testing.B) {
	buf := fast.NewStringBuffer(256)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		encodeQuery(buf, "SELECT * FROM %T WHERE foo = %d", []any{Table("trudeluttan"), 123})
	}
}
