package dense

import (
	"math/bits"
	"unsafe"
)

// vector for two bit values
type bibitvec uint8

func (b bibitvec) get(index byte) byte {
	return byte(b&(0b11000000>>(index*2))) >> (6 - index*2)
}

func (b *bibitvec) set(index byte, val byte) {
	index *= 2
	tmp := *b
	writeMask := byte(0b11000000 >> index)
	clearMask := ^writeMask
	tmp &= bibitvec(clearMask)
	t1 := (val & 0b11) << (6 - index)
	tmp |= bibitvec(writeMask & t1)
	*b = tmp
}

func (b bibitvec) byte() byte {
	return byte(b)
}

type bitvec []bibitvec

func (b bitvec) at(index uint32) byte {
	i, off := bits.Div32(0, index, 4)
	return b[i].get(byte(off))
}

func (b bitvec) set(index uint32, val byte) {
	i, off := bits.Div32(0, index, 4)
	b[i].set(byte(off), val)
}

func (b bitvec) bytes() []byte {
	return unsafe.Slice((*byte)(&b[0]), len(b))
}

func bvFromBytes(other []byte) bitvec {
	return unsafe.Slice((*bibitvec)(&other[0]), len(other))
}
