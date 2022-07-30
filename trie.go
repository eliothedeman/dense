package dense

import (
	"log"
	"unsafe"
)

const all1Byte = 0b11111111
const branchFactor = 2
const partsPerByte = 8 / branchFactor
const nodeChildWidth = 1 << branchFactor
const readMask = all1Byte >> (8 - branchFactor)
const rootNode = 0

// cost = (sizeof(id) * nodeChildWdith) * unique_key_bytes * parts_per_key

type id = uint16

func init() {
	switch branchFactor {
	case 1, 2, 4, 8:

	default:
		log.Panicf("nodeBitWidth must be a power of 2 between 1 and 8")

	}
}

const (
	hasValue = 1 << iota
	deleted
)

type tnode[T any] struct {
	children [nodeChildWidth]id
	flags    uint8
	value    T
}

type Trie[K Key, V any] struct {
	keyBuilder[K]
	t trie[V]
}

func (t *Trie[K, V]) Get(key K) (V, bool) {
	return t.t.Get(t.keyBuilder.key(key))
}

func (t *Trie[K, V]) MustGet(key K) V {
	return t.t.MustGet(t.keyBuilder.key(key))
}

func (t *tnode[T]) hasValue() bool {
	return t.flags&hasValue > 0
}

type trie[T any] struct {
	nodes []tnode[T]
}

// private

func (t *trie[T]) addNode(from int, index id) id {
	next := len(t.nodes)
	t.nodes = append(t.nodes, tnode[T]{})
	offset := id(next - int(from))
	t.nodes[from].children[index] = offset
	return offset
}

type visitor[T any] func(parent int, n *tnode[T]) bool

func (t *trie[T]) dfsNodes(f visitor[T]) bool {
	stack := stack[pair[int, int]]{}
	for _, c := range t.nodes[rootNode].children {
		if c != 0 {
			stack.push(pair[int, int]{rootNode, int(c)})
		}
	}
	for stack.len() > 0 {
		args := stack.pop()
		n := &t.nodes[args.b]
		if !f(args.a, n) {
			return false
		}
		for _, c := range n.children {
			if c != 0 {
				stack.push(pair[int, int]{args.b, args.b + int(c)})
			}
		}
	}
	return true
}

// public

type Option[T any] func(t *trie[T]) *trie[T]

func newTrie[T any]() *trie[T] {
	t := &trie[T]{}
	t.nodes = make([]tnode[T], 0, 1024)
	t.nodes = append(t.nodes, tnode[T]{})
	return t
}

func (t *trie[T]) ForEach(f func([]byte, T)) {
	// buff := make([]byte, 1024)
}

func (t *trie[T]) sizeBytes() int {
	return int(unsafe.Sizeof(t.nodes[rootNode])) * len(t.nodes)
}

func (t *trie[T]) findNextPart(from int, key bitvec, depth uint32) (int, byte) {
	currentNode := &t.nodes[from]
	childIndex := key.at(depth)
	offset := currentNode.children[childIndex]
	if offset == 0 {
		return -1, childIndex
	}
	return from + int(offset), childIndex
}

func (t *trie[T]) insert(key bitvec, val T) {
	parts := uint32(len(key) * partsPerByte)
	index := rootNode
	for i := uint32(0); i < parts; i++ {
		tmp, childIndex := t.findNextPart(index, key, i)
		if tmp < 1 {
			index += int(t.addNode(index, id(childIndex)))
			continue
		}
		index = tmp
	}
	n := &t.nodes[index]

	if n.hasValue() {
		log.Fatalf("Overwriting %v with %v at key '%b'", n.value, val, key)
	}

	n.value = val
	n.flags |= hasValue
}

func (t *trie[T]) createNodesTo(key bitvec) {
	parts := uint32(len(key) * partsPerByte)
	index := rootNode
	for i := uint32(0); i < parts; i++ {
		tmp, childIndex := t.findNextPart(index, key, i)
		if tmp < 1 {
			index += int(t.addNode(index, id(childIndex)))
			continue
		}
		index = tmp
	}
}

func (t *trie[T]) MustGet(key bitvec) (val T) {
	parts := uint32(len(key) * partsPerByte)
	index := rootNode
	for i := uint32(0); i < parts; i++ {
		index, _ = t.findNextPart(index, key, i)
		if index < 1 {
			log.Panicf("key not found %s", key)
		}

	}
	val = t.nodes[index].value
	return
}

func (t *trie[T]) Get(key bitvec) (val T, found bool) {
	parts := uint32(len(key) * partsPerByte)
	index := rootNode
	for i := uint32(0); i < parts; i++ {
		index, _ = t.findNextPart(index, key, i)
		if index < 1 {
			return
		}

	}
	val = t.nodes[index].value
	found = true
	return
}

// helpers

func bitsAtDepth(data []byte, depth id) byte {
	index := depth / partsPerByte
	shift := branchFactor * (depth % partsPerByte)
	mask := byte(readMask << shift)
	return (data[index] & mask) >> shift
}
