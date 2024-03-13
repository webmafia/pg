package pg

import (
	"context"
	"fmt"
)

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
