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
