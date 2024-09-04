package types

import "unsafe"

type Primitive interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~complex64 | ~complex128 | ~bool
}

// Function to produce a byte slice from an input value.
//
// Bytes in the padding section(s) of structs may be set to values other than 0 and can be modified during runtime.
// In consequence, identical structs could have different padding bytes values.
//
// This should be implemented by users for their structs to provide a reliable access to the struct's raw bytes.
// `reflect` is not used since it would have a negative performance overhead.
//
// Example struct reader implementation:
//
//	type testStruct struct {
//	    A int
//	    B bool
//	    C uint64
//	}
//
//	 func testStructBytesReader(value testStruct) []byte {
//	     return slices.Concat(
//	        types.PrimitiveTypeBytesReader(value.A),
//	        types.PrimitiveTypeBytesReader(value.B),
//	        types.PrimitiveTypeBytesReader(value.C),
//	     )
//	 }
type ReaderFunc[T any] func(input T) []byte

// Read raw bytes from a primitive type value
func PrimitiveTypeBytesReader[T Primitive](value T) []byte {
	ptr := unsafe.Pointer(&value)
	length := unsafe.Sizeof(value)
	return unsafe.Slice((*byte)(ptr), length)
}

// Read bytes from a string
func StringBytesReader(value string) []byte {
	return []byte(value)
}

// Read a pointer's `uintptr` bytes
func PointerBytesReader[TPointer *T, T any](pointer TPointer) []byte {
	uintPtr := uintptr(unsafe.Pointer(pointer))
	return PrimitiveTypeBytesReader(uintPtr)
}
