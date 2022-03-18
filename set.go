package dense

import (
	"hash/maphash"
	"math/big"

	"github.com/eliothedeman/fn"
)

type unit struct{}

type Set[T comparable] struct {
	m Map[T, unit]
}

func NewSet[T comparable](size int) *Set[T] {
	s := &Set[T]{m: Map[T, unit]{
		hasher: maphash.Hash{},
		data:   make([]Pair[T, unit], size),
		isSet:  big.Int{},
	}}
	s.m.hasher.Reset()
	if size == 0 {
		s.m.grow()

	}

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

func (s *Set[T]) Iter() *fn.Iter[T] {
	return fn.Map(s.m.Iter(), func(t Pair[T, unit]) T {
		return t.Key
	})
}

func (s *Set[T]) Clear() {
	s.m.Clear()
}

func (s *Set[T]) Delete(k T) {
	s.m.Delete(k)
}
