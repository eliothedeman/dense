package dense

import "log"

const all1Byte = 0b11111111
const nodeBitWidth = 2
const partsPerByte = 8 / nodeBitWidth
const nodeChildWidth = 1 << nodeBitWidth
const readMask = all1Byte >> (8 - nodeBitWidth)
const rootNode = 0

// cost = (sizeof(id) * nodeChildWdith) * unique_key_bytes * parts_per_key

type id = uint16

const (
	hasValue = 1 << iota
	deleted
)

type tnode[T any] struct {
	children [nodeChildWidth]id
	flags    uint8
	value    T
}

type Trie[T any] struct {
	nodes []tnode[T]
}

// private

func (t *Trie[T]) addNode(from id, index id) id {
	next := id(len(t.nodes))
	t.nodes = append(t.nodes, tnode[T]{})
	t.nodes[from].children[index] = next
	return next
}

// public

type Option[T any] func(t *Trie[T]) *Trie[T]

func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{
		nodes: []tnode[T]{
			{},
		},
	}
}

func (t *Trie[T]) ForEach(f func([]byte, T)) {

}

func (t *Trie[T]) Insert(key []byte, val T) {
	parts := id(len(key) * partsPerByte)
	currentNodeId := id(rootNode)
	for i := id(0); i < parts; i++ {
		currentNode := &t.nodes[currentNodeId]
		childIndex := bitsAtDepth(key, i)
		nextChild := currentNode.children[childIndex]
		// log.Printf("%b %b", key, childIndex)
		if nextChild == 0 {
			currentNodeId = t.addNode(currentNodeId, id(childIndex))
			continue
		}
		currentNodeId = nextChild
	}
	n := &t.nodes[currentNodeId]

	if n.flags&hasValue > 0 {
		log.Fatalf("Overwriting %v with %v at key %s", n.value, val, key)
	}

	n.value = val
	n.flags |= hasValue
}

func (t *Trie[T]) MustGet(key []byte) (val T) {
	parts := id(len(key) * partsPerByte)
	currentNodeId := id(rootNode)
	for i := id(0); i < parts; i++ {
		currentNode := &t.nodes[currentNodeId]
		childIndex := bitsAtDepth(key, i)
		nextChild := currentNode.children[childIndex]
		if nextChild == 0 {
			log.Panicf("key not found %s", key)
		}
		currentNodeId = nextChild
	}
	val = t.nodes[currentNodeId].value
	return
}

func (t *Trie[T]) Get(key []byte) (val T, found bool) {
	parts := id(len(key) * partsPerByte)
	currentNodeId := id(rootNode)
	for i := id(0); i < parts; i++ {
		currentNode := &t.nodes[currentNodeId]
		childIndex := bitsAtDepth(key, i)
		nextChild := currentNode.children[childIndex]
		if nextChild == 0 {
			return
		}
		currentNodeId = nextChild
	}
	val = t.nodes[currentNodeId].value
	return
}

// helpers

func bitsAtDepth(data []byte, depth id) byte {
	index := depth / partsPerByte
	shift := nodeBitWidth * (depth % partsPerByte)
	mask := byte(readMask << shift)
	return (data[index] & mask) >> shift
}
