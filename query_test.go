package pg

import (
	"context"
	"fmt"
	"testing"

	"github.com/webmafia/fast"
)

func Example_encodeQuery() {
	buf := fast.NewStringBuffer(256)
	var queryArgs []any
	// q := Query("SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s")
	// q.EncodeQuery(buf, []any{Table("trudeluttan"), 123, Table("mjau"), 456}, &queryArgs)
	encodeQuery(buf, "SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s", []any{Identifier("trudeluttan"), 123, Identifier("mjau"), 456}, &queryArgs)

	fmt.Println(buf.String(), queryArgs)

	// Output: Mjau
}

func Example_encodeQuery2() {
	buf := fast.NewStringBuffer(256)
	var queryArgs []any
	// q := Query("SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s")
	// q.EncodeQuery(buf, []any{Table("trudeluttan"), 123, Table("mjau"), 456}, &queryArgs)
	encodeQuery(buf, "SELECT * FROM %T WHERE %T", []any{
		Identifier("trudeluttan"),
		Or(
			Eq("foo", "bar"),
			Eq("baz", "bez"),
		),
	}, &queryArgs)

	fmt.Println(buf.String(), queryArgs)

	// Output: Mjau
}

func Benchmark_encodeQuery(b *testing.B) {
	buf := fast.NewStringBuffer(256)
	queryArgs := make([]any, 0, 5)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encodeQuery(buf, "SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s", []any{Identifier("trudeluttan"), 123, Identifier("mjau"), 456}, &queryArgs)
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
			Identifier("trudeluttan"),
			Or(
				Eq("foo", "bar"),
				Eq("baz", "bez"),
			),
		}, &queryArgs)
		buf.Reset()
		queryArgs = queryArgs[:0]
	}
}

func BenchmarkQuery(b *testing.B) {
	db := NewDB(nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.Query(context.Background(), "SELECT * FROM %T WHERE foo = %d AND bar = %s AND baz = %s", Identifier("trudeluttan"), 123, Identifier("mjau"), 456)
	}
}
