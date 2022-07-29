package dense

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesFromKey(t *testing.T) {
	x := NewNumericTrie[int, int]()
	for i := 0; i < 82; i++ {
		x.insert(i, i)

	}
	assert.Equal(t, x.MustGet(81), 81)
}

func TestSingleByteKey(t *testing.T) {
	b := NewBuilder[byte, int]()
	for i := 0; i < 255; i++ {
		b.Add(byte(i), i)
	}
	x := b.Build()
	assert.Equal(t, 1, x.MustGet(1))
}

func BenchmarkinsertInt(b *testing.B) {
	b.Run("trie", func(b *testing.B) {
		x := NewNumericTrie[int, int]()
		for i := 0; i < b.N; i++ {
			x.insert(i, i)
		}
	})

	b.Run("hash", func(b *testing.B) {
		x := make(map[int]int)
		for i := 0; i < b.N; i++ {
			x[i] = i
		}
	})
}

func BenchmarkGetInt(b *testing.B) {
	b.Run("trie", func(b *testing.B) {
		x := NewBuilder[int, int]()
		for i := 0; i < b.N; i++ {
			x.Add(i, i)
		}
		y := x.Build()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = y.MustGet(i)
		}
	})
	b.Run("hash", func(b *testing.B) {
		x := make(map[int]int)
		for i := 0; i < b.N; i++ {
			x[i] = i
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = x[i]
		}
	})
}
