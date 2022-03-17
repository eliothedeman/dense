package dense

import (
	"hash/maphash"
	"math/big"
)

type unit struct{}

type Set[T comparable] struct {
	m Map[T, unit]
}

func NewSet[T comparable](size int) *Set[T] {
	s := &Set[T]{m: Map[T, unit]{
		hasher: maphash.Hash{},
		data:   make([]pair[T, unit], size),
		isSet:  big.Int{},
	}}
	s.m.hasher.Reset()

	return s
}

func (s *Set[T]) Insert(val T) {
	s.m.Insert(val, unit{})
}

func (s *Set[T]) Contains(val T) bool {
	return s.m.Contains(val)
}

func (s *Set[T]) Len() int {
	return s.m.numElements
}
