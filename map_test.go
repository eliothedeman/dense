package dense

import "testing"

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
