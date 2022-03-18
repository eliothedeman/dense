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
	if size == 0 {
		s.grow()
	}
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

func (m *Map[K, V]) find(key K) int {
	h := hashT(&m.hasher, key)
	l := uint64(len(m.data))
	mod := h % l
	for i := mod; i < l; i++ {
		if m.isSet.Bit(int(i)) == 1 {
			if m.data[i].Key == key {
				return int(i)
			}
		} else {
			break
		}
	}
	return -1
}

func (m *Map[K, V]) Get(key K) (val V, hasval bool) {
	idx := m.find(key)
	if idx < 0 {
		return
	}
	val = m.data[idx].Val
	hasval = true
	return
}

func (m *Map[K, V]) GetRef(key K) (val *V, hasval bool) {

	idx := m.find(key)
	if idx < 0 {
		return
	}
	val = &m.data[idx].Val
	hasval = true
	return
}

func (m *Map[K, V]) MustGet(key K) (out V) {
	out = *m.MustGetRef(key)
	return
}
func (m *Map[K, V]) MustGetRef(key K) (out *V) {

	var hasVal bool
	out, hasVal = m.GetRef(key)
	if hasVal {
		return
	}
	log.Panicf("%v not in map[%p]", key, m)
	return
}

func (m *Map[K, V]) Contains(key K) bool {
	return m.find(key) >= 0
}

func (m *Map[K, V]) Len() int {
	return m.numElements
}

func (m *Map[K, V]) Iter() *fn.Iter[Pair[K, V]] {
	i := 0
	return fn.NewIter(func() (out fn.Option[Pair[K, V]]) {
		l := len(m.data)
		if i >= l {
			return fn.None[Pair[K, V]]()
		}
		for i < l {
			if m.isSet.Bit(i) == 1 {
				out = fn.Some(m.data[i])
				i++
				return
			}
			i++
		}
		return fn.None[Pair[K, V]]()
	})
}

func (m *Map[K, V]) Clear() {
	m.isSet.SetInt64(0)
}

func (m *Map[K, V]) Delete(k K) {
	idx := m.find(k)
	// not in the map
	if idx < 0 {
		return
	}
	var blankV V
	var blankK K
	m.data[idx] = Pair[K, V]{blankK, blankV}
	m.isSet.SetBit(&m.isSet, idx, 0)
}
