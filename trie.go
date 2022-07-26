package dense

import (
	"log"
	"unsafe"
)

const all1Byte = 0b11111111
const nodeBitWidth = 2
const partsPerByte = 8 / nodeBitWidth
const nodeChildWidth = 1 << nodeBitWidth
const readMask = all1Byte >> (8 - nodeBitWidth)
const rootNode = 0

// cost = (sizeof(id) * nodeChildWdith) * unique_key_bytes * parts_per_key

type id = uint32

func init() {
	switch nodeBitWidth {
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

func (t *tnode[T]) hasValue() bool {
	return t.flags&hasValue > 0
}

type Trie[T any] struct {
	nodes []tnode[T]
}

// private

func (t *Trie[T]) addNode(from int, index id) id {
	next := len(t.nodes)
	t.nodes = append(t.nodes, tnode[T]{})
	offset := id(next - int(from))
	t.nodes[from].children[index] = offset
	return offset
}

type visitor[T any] func(parent id, n *tnode[T]) bool

func (t *Trie[T]) dfsNodes(f visitor[T]) bool {
	stack := stack[pair[id, id]]{}
	for _, c := range t.nodes[rootNode].children {
		if c != 0 {
			stack.push(pair[id, id]{rootNode, c})
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
				stack.push(pair[id, id]{args.b, c})
			}
		}
	}
	return true
}

// public

type Option[T any] func(t *Trie[T]) *Trie[T]

func NewTrie[T any]() *Trie[T] {
	t := &Trie[T]{}
	t.nodes = make([]tnode[T], 0, 1024)
	t.nodes = append(t.nodes, tnode[T]{})
	return t
}

func (t *Trie[T]) ForEach(f func([]byte, T)) {
	// buff := make([]byte, 1024)
}

func (t *Trie[T]) sizeBytes() int {
	return int(unsafe.Sizeof(t.nodes[rootNode])) * len(t.nodes)
}

func (t *Trie[T]) findNextPart(from int, key []byte, depth id) (int, byte) {
	currentNode := &t.nodes[from]
	childIndex := bitsAtDepth(key, depth)
	offset := currentNode.children[childIndex]
	if offset == 0 {
		return -1, childIndex
	}
	return from + int(offset), childIndex
}

func (t *Trie[T]) Insert(key []byte, val T) {
	parts := id(len(key) * partsPerByte)
	index := rootNode
	for i := id(0); i < parts; i++ {
		tmp, childIndex := t.findNextPart(index, key, i)
		if tmp < 1 {
			index += int(t.addNode(index, id(childIndex)))
			continue
		}
		index = tmp
	}
	n := &t.nodes[index]

	if n.flags&hasValue > 0 {
		log.Fatalf("Overwriting %v with %v at key '%b'", n.value, val, key)
	}

	n.value = val
	n.flags |= hasValue
}

func (t *Trie[T]) MustGet(key []byte) (val T) {
	parts := id(len(key) * partsPerByte)
	index := rootNode
	for i := id(0); i < parts; i++ {
		index, _ = t.findNextPart(index, key, i)
		if index < 1 {
			log.Panicf("key not found %s", key)
		}

	}
	val = t.nodes[index].value
	return
}

func (t *Trie[T]) Get(key []byte) (val T, found bool) {
	parts := id(len(key) * partsPerByte)
	index := rootNode
	for i := id(0); i < parts; i++ {
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
	shift := nodeBitWidth * (depth % partsPerByte)
	mask := byte(readMask << shift)
	return (data[index] & mask) >> shift
}
