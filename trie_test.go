package dense

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkInsert(b *testing.B) {
	var buff bytes.Buffer
	for size := range []int{1, 2, 4, 8, 16, 32} {
		size = 1024 << size
		b.Run(fmt.Sprintf("trie_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x := NewTrie[int]()
				for y := 0; y < size; y++ {
					fmt.Fprint(&buff, y)
					x.Insert(buff.Bytes(), y)
					buff.Reset()
				}
			}
		})
		b.Run(fmt.Sprintf("hash_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x := make(map[string]int)
				for y := 0; y < size; y++ {
					fmt.Fprint(&buff, y)
					x[unsafeAsString(buff.Bytes())] = y
					buff.Reset()
				}
			}
		})

	}
}

func TestInsertGet(t *testing.T) {
	x := NewTrie[int]()
	x.Insert([]byte("10"), 10)
	x.Insert([]byte("11"), 11)
	assert.Equal(t, x.MustGet([]byte("10")), 10)
	assert.Equal(t, x.MustGet([]byte("11")), 11)
	assert.Panics(t, func() {
		x.MustGet([]byte("nope"))
	})
}

func TestInsert(t *testing.T) {
	x := NewTrie[int]()
	x.Insert([]byte("4"), 10)
}
