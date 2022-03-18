package dense

import (
	"hash/maphash"
	"log"
	"math/big"

	"github.com/eliothedeman/fn"
	"golang.org/x/exp/constraints"
)

type Pair[K, V any] struct {
	Key K
	Val V
}

type Map[K comparable, V any] struct {
	hasher      maphash.Hash
	data        []Pair[K, V]
	isSet       big.Int
	numElements int
}

func NewMap[K comparable, V any](size int) *Map[K, V] {
	s := &Map[K, V]{
		hasher: maphash.Hash{},
		data:   make([]Pair[K, V], size),
		isSet:  big.Int{},
	}
	s.hasher.Reset()
	s.grow()
	return s
}

func min[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a

	}
	return b
}

func max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func (m *Map[K, V]) grow() {
	oldD := m.data
	oldI := m.isSet
	m.data = make([]Pair[K, V], max(len(m.data)*2, 2))
	m.isSet.SetInt64(0)
	for i := range oldD {
		kv := &m.data[i]
		if oldI.Bit(i) == 1 {
			m.Insert(kv.Key, kv.Val)
		}
	}
}

func (m *Map[K, V]) Insert(key K, val V) {
	m.numElements += 1
	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	foundSlot := false
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 0 {
			m.isSet.SetBit(&m.isSet, int(i), 1)
			foundSlot = true
			m.data[i] = Pair[K, V]{key, val}
			return
		}
		// already in the map
		if m.data[i].Key == key {
			return
		}
	}

	if !foundSlot {
		m.grow()
		m.Insert(key, val)
	}
}

func (m *Map[K, V]) Get(key K) (val V, hasval bool) {

	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 1 {
			if m.data[i].Key == key {
				return m.data[i].Val, true
			}
		} else {
			break
		}
	}
	return
}

func (m *Map[K, V]) GetRef(key K) (val *V, hasval bool) {

	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 1 {
			if m.data[i].Key == key {
				return &m.data[i].Val, true
			}
		} else {
			break
		}
	}
	return
}

func (m *Map[K, V]) MustGet(key K) (out V) {

	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 1 {
			if m.data[i].Key == key {
				out = m.data[i].Val
				return
			}
		} else {
			break
		}
	}
	log.Panicf("%v not in map[%p]", key, m)
	return
}
func (m *Map[K, V]) MustGetRef(key K) (out *V) {

	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 1 {
			if m.data[i].Key == key {
				out = &m.data[i].Val
				return
			}
		} else {
			break
		}
	}
	log.Panicf("%v not in map[%p]", key, m)
	return
}

func (m *Map[K, V]) Contains(key K) bool {
	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 1 {
			if m.data[i].Key == key {
				return true
			}
		} else {
			break
		}
	}
	return false
}

func (m *Map[K, V]) Len() int {
	return m.numElements
}

func (m *Map[K, V]) Iter() *fn.Iter[Pair[K, V]] {
	i := 0
	return fn.NewIter(func() (out fn.Option[Pair[K, V]]) {
		hasNext := false
		for m.isSet.Bit(i) != 1 && i <= len(m.data) {
			hasNext = true
			i++
		}
		if !hasNext {
			return fn.None[Pair[K, V]]()
		}
		out = fn.Some(m.data[i])
		i++
		return
	})
}
