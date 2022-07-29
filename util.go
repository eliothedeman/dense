package dense

import (
	"reflect"
	"unsafe"
)

func unsafeAsString(b []byte) string {
	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data)

	var s string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(p)
	hdr.Len = len(b)
	return s
}

func cast[T any](v any) T {
	return v.(T)
}

type stack[T any] struct {
	data []T
}

type pair[A, B any] struct {
	a A
	b B
}

func (s *stack[T]) push(t T) {
	s.data = append(s.data, t)
}
func (s *stack[T]) pop() T {
	top := len(s.data) - 1
	val := s.data[top]
	s.data = s.data[:top]
	return val
}

func (s *stack[T]) peek() T {
	return s.data[len(s.data)-1]
}

func (s *stack[T]) len() int {
	return len(s.data)
}

type keyBuilder[K Key] struct {
	keyBuff [1024]byte
	fn      func(key K, buff []byte) []byte
}

func newKeyBuilder[K Key]() keyBuilder[K] {
	return keyBuilder[K]{
		fn: getKeyFunc[K](),
	}
}

func (k *keyBuilder[K]) key(key K) []byte {
	return k.fn(key, k.keyBuff[:])
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
func getKeyFunc[K Key]() func(K, []byte) []byte {
	var k K
	switch any(k).(type) {
	case string:
		return func(key K, buff []byte) []byte {
			return buff[:copy(buff, any(key).(string))]
		}
	case []byte:
		return func(key K, buff []byte) []byte {
			return any(key).([]byte)
		}
	default:
		return bytesFromPOD[K]
	}

}
