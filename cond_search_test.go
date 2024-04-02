package pg

import (
	"fmt"
	"testing"
)

func ExamplePrefixSearch() {
	fmt.Println(PrefixSearch("  hello      world    "))
	// Output: hello:* world:*
}

func BenchmarkPrefixSearch(b *testing.B) {
	str := "  hello      world    "
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = PrefixSearch(str)
	}
}
