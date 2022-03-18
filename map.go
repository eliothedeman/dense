package dense

import (
	"log"

	"golang.org/x/exp/constraints"
)

type Pair[K, V any] struct {
	Key K
	Val V
}

type Map[K comparable, V any] struct {
	lookup map[K]uint32
	data   []V
}

func (m *Map[K, V]) ensureInit() {
	if len(m.data) == 0 {
		m.data = make([]V, 2)
	}
	if m.data == nil {
		m.lookup = make(map[K]uint32, len(m.data))
	}
}

func NewMap[K comparable, V any](size int) *Map[K, V] {
	m := &Map[K, V]{
		data:   make([]V, max(size, 2)),
		lookup: make(map[K]uint32, max(size, 2)),
	}
	return m
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

func (m *Map[K, V]) Insert(key K, val V) {
	m.ensureInit()
	l := len(m.data)
	m.data = append(m.data, val)
	m.lookup[key] = uint32(l)
}

func (m *Map[K, V]) find(key K) int {
	m.ensureInit()
	idx, ok := m.lookup[key]
	if !ok {
		return -1
	}

	return int(idx)
}

func (m *Map[K, V]) Get(key K) (val V, hasval bool) {
	idx := m.find(key)
	if idx < 0 {
		return
	}
	val = m.data[idx]
	hasval = true
	return
}

func (m *Map[K, V]) GetRef(key K) (val *V, hasval bool) {

	idx := m.find(key)
	if idx < 0 {
		return
	}
	val = &m.data[idx]
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
	return len(m.lookup)
}

func (m *Map[K, V]) Each(f func(key K, val V)) {
	for k, v := range m.lookup {
		f(k, m.data[v])
	}
}

func (m *Map[K, V]) EachRef(f func(key K, valRef *V)) {
	for k, v := range m.lookup {
		f(k, &m.data[v])
	}
}

func (m *Map[K, V]) Clear() {
	m.lookup = nil
	m.ensureInit()
}

func (m *Map[K, V]) Delete(k K) {
	old, ok := m.lookup[k]
	if !ok {
		return
	}
	delete(m.lookup, k)
	var blankV V
	m.data[old] = blankV
}
