package hasher

import (
	"math/rand"
	"unsafe"
)

// Get a hash function for the given type.
//
// The function relies on Go's internal hasher implementation for this type.
func GetHashFunc[T comparable]() func(uintptr, uintptr) uintptr {
	nativeHmap := any((map[T]struct{})(nil))
	return (*emptyInterface)(unsafe.Pointer(&nativeHmap))._type.Hasher
}

// Produce a random seed
func GenerateSeed() uintptr {
	return uintptr(rand.Uint64())
}

// Internal interface type
//
// Source: runtime/runtime2.go
type emptyInterface struct {
	_type *mapType
	data  unsafe.Pointer
}

// Internal map type
//
// Note: in Go's source code, Hasher takes an unsafe.Pointer as its first argument,
// It has been replaced by a uintptr here to avoid the argument escaping to the Heap.
// The call is seemless but it is a hack.
//
// Source: internal/abi/map_swiss.go
type mapType struct {
	internalType
	Key    *internalType
	Elem   *internalType
	Bucket *internalType // internal type representing a hash bucket
	// function for hashing keys (ptr to key, seed) -> hash
	Hasher     func(uintptr, uintptr) uintptr
	KeySize    uint8  // size of key slot
	ValueSize  uint8  // size of elem slot
	BucketSize uint16 // size of bucket
	Flags      uint32
}

// Internal Type representation
//
// Source: internal/abi/type.go
type internalType struct {
	Size_       uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // hash of type; avoids computation in hash tables
	TFlag       uint8   // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, GCData is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	GCData    *byte
	Str       int32 // string form
	PtrToThis int32 // type for pointer to this type, may be zero
}
