package dense

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

type NumericKey interface {
	constraints.Unsigned | constraints.Signed
}

type NumericTrie[K NumericKey, V any] struct {
	keyBuff [16]byte
	raw     Trie[V]
}

func NewNumericTrie[K NumericKey, V any]() *NumericTrie[K, V] {
	return &NumericTrie[K, V]{
		raw: *NewTrie[V](),
	}
}

func (n *NumericTrie[K, V]) bytesFromKey(key K) []byte {
	size := int(unsafe.Sizeof(key))
	buff := n.keyBuff[:size]

	raw := unsafe.Slice((*byte)(unsafe.Pointer(&key)), size)
	size--
	for i := range buff {
		buff[i] = raw[size-i]
	}
	return buff
}

func (n *NumericTrie[K, V]) Insert(key K, val V) {
	rawKey := n.bytesFromKey(key)
	n.raw.Insert(rawKey, val)
}

func (n *NumericTrie[K, V]) MustGet(key K) V {
	rawKey := n.bytesFromKey(key)
	return n.raw.MustGet(rawKey)
}

func (n *NumericTrie[K, V]) Get(key K) (V, bool) {
	rawKey := n.bytesFromKey(key)
	return n.raw.Get(rawKey)
}
