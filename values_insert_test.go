package pg

import (
	"context"
	"fmt"
)

// func BenchmarkInsert(b *testing.B) {
// 	db := NewDB(nil)

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		db.Insert(context.Background(), Identifier("mjau"),
// 			"foo", "bar",
// 			"foo", "bar",
// 			"foo", "bar",
// 			"foo", "bar",
// 			"foo", "bar",
// 			"foo", "bar",
// 		)
// 	}
// }

func ExampleInsert() {
	db := NewDB(nil)

	vals := db.AcquireValues()
	defer db.ReleaseValues(vals)

	vals.
		Value("foo", "bar").
		Value("baz", 123)

	_, err := db.InsertValues(context.Background(), Identifier("foobar"), vals)

	fmt.Printf("%#v\n", err)

	// Output: Mjau
}
