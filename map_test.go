package dense

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkMapInsert(b *testing.B) {

	b.Run("native_empty", func(b *testing.B) {
		m := make(map[int]int)

		for i := 0; i < b.N; i++ {
			m[i] = i
		}

	})
	b.Run("dense_empty", func(b *testing.B) {
		m := NewMap[int, int](0)
		for i := 0; i < b.N; i++ {
			m.Insert(i, i)
		}

	})
	b.Run("native_prealloc", func(b *testing.B) {
		m := make(map[int]int, 10000)

		for i := 0; i < b.N; i++ {
			m[i] = i
		}

	})
	b.Run("dense_prealloc", func(b *testing.B) {
		m := NewMap[int, int](1000)
		for i := 0; i < b.N; i++ {
			m.Insert(i, i)
		}

	})
}

func TestIterMap(t *testing.T) {
	m := NewMap[int, int](0)
	for i := 0; i < 1000; i++ {
		m.Insert(i, i)
	}
	called := 0
	m.Each(func(key, val int) {

		called++
		assert.Equal(t, m.MustGet(key), val)
	})
	assert.Equal(t, 1000, called)
}
