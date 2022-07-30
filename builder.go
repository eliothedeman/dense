package dense

import (
	"bytes"
	"log"
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type builder[T any] struct {
	pairs []pair[bitvec, T]
}

var _ sort.Interface = &builder[any]{}

func (b *builder[T]) Len() int {
	return len(b.pairs)
}

func (b *builder[T]) Less(i, j int) bool {
	return 0 > bytes.Compare(b.pairs[i].a.bytes(), b.pairs[j].a.bytes())
	// jn := b.pairs[j]
	// in := b.pairs[i]
	// for n, x := range in.a {
	// 	if n >= len(jn.a) {
	// 		return false
	// 	}
	// 	if jn.a[n] > x {
	// 		return true
	// 	}
	// }
	// // they are the same
	// return false
}

func (b *builder[T]) Swap(i, j int) {
	b.pairs[i], b.pairs[j] = b.pairs[j], b.pairs[i]
}

func (b *builder[T]) add(key bitvec, val T) *builder[T] {
	b.pairs = append(b.pairs, pair[bitvec, T]{key, val})
	return b
}

func (b *builder[T]) build() *trie[T] {
	sort.Sort(b)
	l := 0
	t := newTrie[T]()
	keyFound := true
	for keyFound {
		keyFound = false
		for i := range b.pairs {
			p := &b.pairs[i]
			if len(p.a) >= l {
				log.Printf("%o", p.a[:l])
				t.createNodesTo(p.a[:l])
				keyFound = true
			}
		}
		l++
	}
	log.Println("start")
	for i := range b.pairs {
		p := &b.pairs[i]
		log.Printf("%b", p.a)
		t.insert(p.a, p.b)
	}
	log.Println("end")

	return t
}

type Builder[K Key, V any] struct {
	b builder[V]
	keyBuilder[K]
}

type Key interface {
	constraints.Integer | constraints.Float | ~[]byte | ~string
}

func NewBuilder[K Key, V any]() *Builder[K, V] {
	return &Builder[K, V]{
		keyBuilder: newKeyBuilder[K](),
	}
}

func (b *Builder[K, V]) Add(k K, v V) {
	tmp := b.keyBuilder.key(k)
	b.b.add(slices.Clone(tmp), v)
}

func (b *Builder[K, V]) Build() *Trie[K, V] {
	return &Trie[K, V]{t: *b.b.build(), keyBuilder: newKeyBuilder[K]()}
}
