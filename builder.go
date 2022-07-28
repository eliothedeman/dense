package dense

import (
	"bytes"
	"sort"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type builder[T any] struct {
	pairs []pair[[]byte, T]
}

var _ sort.Interface = &builder[any]{}

func (b *builder[T]) Len() int {
	return len(b.pairs)
}

func (b *builder[T]) Less(i, j int) bool {
	return bytes.Compare(b.pairs[i].a, b.pairs[j].a) < 1
}

func (b *builder[T]) Swap(i, j int) {
	b.pairs[i], b.pairs[j] = b.pairs[j], b.pairs[i]
}

func (b *builder[T]) Add(key []byte, val T) *builder[T] {
	b.pairs = append(b.pairs, pair[[]byte, T]{key, val})
	return b
}

func (b *builder[T]) Build() *Trie[T] {
	sort.Sort(b)
	l := 0
	t := newTrie[T]()
	keyFound := true
	for keyFound {
		keyFound = false
		for i := range b.pairs {
			p := &b.pairs[i]
			if len(p.a) >= l {
				t.createNodesTo(p.a[:l])
				keyFound = true
			}
		}
		l++
	}
	for i := range b.pairs {
		p := &b.pairs[i]
		t.insert(p.a, p.b)
	}

	return t
}

type Builder[K Key, V any] struct {
	b       builder[V]
	keyBuff [1024]byte
	fn      func(key K, buff []byte) []byte
}

func bytesFromPOD[K Key](key K, keyBuff []byte) []byte {
	size := int(unsafe.Sizeof(key))
	buff := keyBuff[:size]
	// write most significant bit first, so need to reverse
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&key)), size)
	size--
	for i := range buff {
		buff[i] = raw[size-i]
	}
	return buff
}

type Key interface {
	constraints.Integer | constraints.Float | ~[]byte | ~string
}

func NewBuilder[K Key, V any]() *Builder[K, V] {
	var k K
	b := new(Builder[K, V])
	switch any(k).(type) {
	case string:
		b.fn = func(key K, buff []byte) []byte {
			return buff[:copy(buff, any(key).(string))]
		}
	case []byte:
		b.fn = func(key K, buff []byte) []byte {
			return any(key).([]byte)
		}
	default:
		b.fn = bytesFromPOD[K]
	}
	return b

}
