package dense

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmarkinsert(b *testing.B) {
	var buff bytes.Buffer
	for _, size := range []int{1, 2, 4, 8} {
		size = 1024 << size
		b.Run(fmt.Sprintf("trie_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x := newTrie[int]()
				keySize := 0
				for y := 0; y < size; y++ {
					fmt.Fprint(&buff, y)
					keySize += buff.Len()
					x.insert(bvFromBytes(buff.Bytes()), y)
					buff.Reset()
				}
				if i == 0 {
					b.Logf("trie: %d key: %d ratio: %f", x.sizeBytes(), keySize, float64(x.sizeBytes())/float64(keySize))
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
	x := newTrie[int]()
	x.insert(bitvec("10"), 10)
	x.insert(bitvec("11"), 11)
	assert.Equal(t, x.MustGet(bitvec("10")), 10)
	assert.Equal(t, x.MustGet(bitvec("11")), 11)
	assert.Panics(t, func() {
		x.MustGet(bitvec("nope"))
	})
}

func TestInsert(t *testing.T) {
	x := newTrie[int]()
	x.insert(bitvec("4"), 10)
}
