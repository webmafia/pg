package pg

import (
	"fmt"
	"log"
	"testing"

	"github.com/webmafia/fast"
)

func Example_indexQuery() {
	// buf := fast.NewStringBuffer(256)
	q := indexQuery("SELECT * FROM %T WHERE foo = %d")

	fmt.Printf("%#v", q)

	// Output: Mjau
}

func Example_indexQuery_encode() {
	buf := fast.NewStringBuffer(256)
	q := indexQuery("SELECT * FROM %T WHERE foo = %d")

	err := q.encodeQuery(buf, []any{Table("trudeluttan"), 123})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())

	// Output: Mjau
}

func Benchmark_indexQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = indexQuery(`
			SELECT *
			FROM %T
			WHERE foo = %d
		`)
	}
}

func Benchmark_indexQuery_encode(b *testing.B) {
	buf := fast.NewStringBuffer(256)
	q := indexQuery("SELECT * FROM %T WHERE foo = %d")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := q.encodeQuery(buf, []any{Table("trudeluttan"), 123})

		if err != nil {
			b.Fatal(err)
		}
	}
}
