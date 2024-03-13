package pg

import (
	"fmt"
	"testing"

	"github.com/webmafia/fast"
)

func Example_encodeQuery() {
	buf := fast.NewStringBuffer(256)
	var queryArgs []any
	// q := Query("SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s")
	// q.EncodeQuery(buf, []any{Table("trudeluttan"), 123, Table("mjau"), 456}, &queryArgs)
	// encodeQuery(buf, "SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s", []any{Table("trudeluttan"), 123, Table("mjau"), 456}, &queryArgs)

	fmt.Println(buf.String(), queryArgs)

	// Output: Mjau
}

func Benchmark_encodeQuery(b *testing.B) {
	buf := fast.NewStringBuffer(256)
	queryArgs := make([]any, 0, 5)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encodeQuery(buf, "SELECT * FROM %T WHERE foo = %d", []any{Table("trudeluttan"), 123}, &queryArgs)
		buf.Reset()
		queryArgs = queryArgs[:0]
	}
}
