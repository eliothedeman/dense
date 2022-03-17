package dense

import (
	"hash/maphash"
	"log"
	"math/big"
)

type pair[K, V any] struct {
	key K
	val V
}

type Map[K comparable, V any] struct {
	hasher      maphash.Hash
	data        []pair[K, V]
	isSet       big.Int
	numElements int
}

func NewMap[K comparable, V any](size int) *Map[K, V] {
	s := &Map[K, V]{
		hasher: maphash.Hash{},
		data:   make([]pair[K, V], size),
		isSet:  big.Int{},
	}
	s.hasher.Reset()
	return s
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
			m.data[i] = pair[K, V]{key, val}
			return
		}
		// already in the map
		if m.data[i].key == key {
			return
		}
	}

	if !foundSlot {
		next := NewMap[K, V](len(m.data) * 2)
		next.Insert(key, val)
		for i := range m.data {
			kv := &m.data[i]
			next.Insert(kv.key, kv.val)
		}
		*m = *next
	}
}

func (m *Map[K, V]) Get(key K) (val V, hasval bool) {

	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 1 {
			if m.data[i].key == key {
				return m.data[i].val, true
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
			if m.data[i].key == key {
				out = m.data[i].val
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
			if m.data[i].key == key {
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
