package pg

import (
	"fmt"
	"testing"

	"github.com/webmafia/fast"
)

func Benchmark_encodeQuery(b *testing.B) {
	buf := fast.NewStringBuffer(256)
	queryArgs := make([]any, 0, 5)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encodeQuery(buf, "SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s", []any{Identifier("foo"), 123, Identifier("bar"), 456}, &queryArgs)
		buf.Reset()
		queryArgs = queryArgs[:0]
	}
}

func Benchmark_encodeQuery2(b *testing.B) {
	buf := fast.NewStringBuffer(256)
	queryArgs := make([]any, 0, 5)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encodeQuery(buf, "SELECT * FROM %T WHERE %T", []any{
			Identifier("foobar"),
			Or(
				Eq("foo", "bar"),
				Eq("baz", "bez"),
			),
		}, &queryArgs)
		buf.Reset()
		queryArgs = queryArgs[:0]
	}
}

func Example_encodeQuery() {
	buf := fast.NewStringBuffer(256)
	queryArgs := make([]any, 0, 5)

	err := encodeQuery(buf, "SELECT * FROM table WHERE foo = %[1]d AND bar != %[1]d AND baz != %d", []any{123, 456}, &queryArgs)

	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String(), queryArgs)

	// Output: SELECT * FROM table WHERE foo = $1 AND bar != $2 AND baz != $3 [123 123 456]
}
