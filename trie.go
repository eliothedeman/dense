package dense

import "log"

const nodeMaskEven = 0b0011
const nodeMaskOdd = 0b1100
const oddBit = 0b1
const nodeBitWidth = 2
const nodeChildWidth = 1 << nodeBitWidth
const rootNode = 0

type id = uint16

type tnode[T any] struct {
	value    T
	children [nodeChildWidth]id
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

func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{
		nodes: []tnode[T]{
			{},
		},
	}
}

func (t *Trie[T]) Insert(key []byte, val T) {
	parts := id(len(key) * 2)
	currentNodeId := id(rootNode)
	for i := id(0); i < parts; i++ {
		currentNode := &t.nodes[currentNodeId]
		childIndex := bitsAtDepth(key, i)
		nextChild := currentNode.children[childIndex]
		if nextChild == 0 {
			currentNodeId = t.addNode(currentNodeId, uint16(childIndex))
			continue
		}
		currentNodeId = nextChild
	}
	t.nodes[currentNodeId].value = val
}

func (t *Trie[T]) MustGet(key []byte) (val T) {
	parts := id(len(key) * 2)
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
	parts := id(len(key) * 2)
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

func maskAndShiftAtDepth(depth id) (byte, byte) {
	x := depth & oddBit
	y := nodeMaskEven << x
	y = y << x
	return byte(y), byte(x << x)
}

func bitsAtDepth(data []byte, depth id) byte {
	mask, shift := maskAndShiftAtDepth(depth)
	depth = depth >> 1
	return (data[depth] & mask) >> shift
}
