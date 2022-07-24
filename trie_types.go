package dense

import (
	"reflect"
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
	ptr := unsafe.Pointer(&key)
	header := reflect.SliceHeader{Data: uintptr(ptr), Len: size, Cap: size}
	copy(buff, *(*[]byte)(unsafe.Pointer(&header)))
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
