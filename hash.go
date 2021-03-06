package dense

import (
	"hash/maphash"
	"unsafe"
)

func hashT[T any](h *maphash.Hash, t T) uint64 {
	h.Reset()
	switch t := any(t).(type) {
	case string:
		h.WriteString(t)
	case []byte:
		h.Write(t)
	default:
		size := unsafe.Sizeof(t)
		ptr := uintptr(unsafe.Pointer(&t))
		for i := uintptr(0); i < size; i++ {
			b := (*byte)(unsafe.Pointer(ptr))
			h.WriteByte(*b)
			ptr++
		}

	}
	return h.Sum64()
}
