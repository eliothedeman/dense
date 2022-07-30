package dense

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilderKeys(t *testing.T) {
	keys := []string{"hell", "hello"}
	b := NewBuilder[string, string]()
	for _, k := range keys {
		b.Add(k, k)
	}
	m := b.Build()

	assert.Equal(t, m.MustGet("hell"), "hell")
	assert.Equal(t, m.MustGet("hello"), "hello")
}

// func TestBuildSortOrder(t *testing.T) {
// 	b := NewBuilder[int, int]()
// 	for k := 0; k < 12; k++ {
// 		b.Add(k, k)
// 	}
// 	assert.Equal(t, fn.Map(fn.IterSlice(b.b.pairs), func(p pair[[]byte, int]) int {
// 		return p.b
// 	}).Collect(), []int{0, 1, 10, 11, 2, 3, 4, 5, 6, 7, 8, 9})
// }
