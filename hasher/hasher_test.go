package hasher

import (
	"slices"
	"testing"

	"github.com/valsov/hashmap/types"
)

type testCase[T any] struct {
	value  T
	copy   T // Used to create another hash and comparing it with the one computed from `value`
	reader types.ReaderFunc[T]
}

// Store both a value and its hash
type hashedValue struct {
	value any
	hash  uint64
}

// Sample test struct to check struct and pointer hashing
type testStruct struct {
	A int
	B bool
	C uint64
}

// Function to read `testStruct` raw bytes
func testStructBytesReader(value testStruct) []byte {
	return slices.Concat(
		types.PrimitiveTypeBytesReader(value.A),
		types.PrimitiveTypeBytesReader(value.B),
		types.PrimitiveTypeBytesReader(value.C),
	)
}

func TestHash(t *testing.T) {
	hashHistory := []hashedValue{}

	testHashInternal(testCase[string]{"str1", "str1", types.StringBytesReader}, &hashHistory, t)
	testHashInternal(testCase[string]{"str2", "str2", types.StringBytesReader}, &hashHistory, t)
	testHashInternal(testCase[rune]{'a', 'a', types.PrimitiveTypeBytesReader[rune]}, &hashHistory, t)
	testHashInternal(testCase[rune]{'b', 'b', types.PrimitiveTypeBytesReader[rune]}, &hashHistory, t)
	testHashInternal(testCase[byte]{byte(128), byte(128), types.PrimitiveTypeBytesReader[byte]}, &hashHistory, t)
	testHashInternal(testCase[byte]{byte(120), byte(120), types.PrimitiveTypeBytesReader[byte]}, &hashHistory, t)
	testHashInternal(testCase[int8]{int8(10), int8(10), types.PrimitiveTypeBytesReader[int8]}, &hashHistory, t)
	testHashInternal(testCase[int8]{int8(11), int8(11), types.PrimitiveTypeBytesReader[int8]}, &hashHistory, t)
	testHashInternal(testCase[uint8]{uint8(12), uint8(12), types.PrimitiveTypeBytesReader[uint8]}, &hashHistory, t)
	testHashInternal(testCase[uint8]{uint8(13), uint8(13), types.PrimitiveTypeBytesReader[uint8]}, &hashHistory, t)
	testHashInternal(testCase[int16]{int16(10), int16(10), types.PrimitiveTypeBytesReader[int16]}, &hashHistory, t)
	testHashInternal(testCase[int16]{int16(11), int16(11), types.PrimitiveTypeBytesReader[int16]}, &hashHistory, t)
	testHashInternal(testCase[uint16]{uint16(12), uint16(12), types.PrimitiveTypeBytesReader[uint16]}, &hashHistory, t)
	testHashInternal(testCase[uint16]{uint16(13), uint16(13), types.PrimitiveTypeBytesReader[uint16]}, &hashHistory, t)
	testHashInternal(testCase[int32]{int32(10), int32(10), types.PrimitiveTypeBytesReader[int32]}, &hashHistory, t)
	testHashInternal(testCase[int32]{int32(11), int32(11), types.PrimitiveTypeBytesReader[int32]}, &hashHistory, t)
	testHashInternal(testCase[uint32]{uint32(12), uint32(12), types.PrimitiveTypeBytesReader[uint32]}, &hashHistory, t)
	testHashInternal(testCase[uint32]{uint32(13), uint32(13), types.PrimitiveTypeBytesReader[uint32]}, &hashHistory, t)
	testHashInternal(testCase[int64]{int64(10), int64(10), types.PrimitiveTypeBytesReader[int64]}, &hashHistory, t)
	testHashInternal(testCase[int64]{int64(11), int64(11), types.PrimitiveTypeBytesReader[int64]}, &hashHistory, t)
	testHashInternal(testCase[uint64]{uint64(12), uint64(12), types.PrimitiveTypeBytesReader[uint64]}, &hashHistory, t)
	testHashInternal(testCase[uint64]{uint64(13), uint64(13), types.PrimitiveTypeBytesReader[uint64]}, &hashHistory, t)
	testHashInternal(testCase[int]{int(14), int(14), types.PrimitiveTypeBytesReader[int]}, &hashHistory, t)
	testHashInternal(testCase[int]{int(15), int(15), types.PrimitiveTypeBytesReader[int]}, &hashHistory, t)
	testHashInternal(testCase[uint]{uint(16), uint(16), types.PrimitiveTypeBytesReader[uint]}, &hashHistory, t)
	testHashInternal(testCase[uint]{uint(17), uint(17), types.PrimitiveTypeBytesReader[uint]}, &hashHistory, t)
	testHashInternal(testCase[uintptr]{uintptr(18), uintptr(18), types.PrimitiveTypeBytesReader[uintptr]}, &hashHistory, t)
	testHashInternal(testCase[uintptr]{uintptr(19), uintptr(19), types.PrimitiveTypeBytesReader[uintptr]}, &hashHistory, t)
	testHashInternal(testCase[float32]{float32(10), float32(10), types.PrimitiveTypeBytesReader[float32]}, &hashHistory, t)
	testHashInternal(testCase[float32]{float32(12), float32(12), types.PrimitiveTypeBytesReader[float32]}, &hashHistory, t)
	testHashInternal(testCase[float64]{float64(10), float64(10), types.PrimitiveTypeBytesReader[float64]}, &hashHistory, t)
	testHashInternal(testCase[float64]{float64(12), float64(12), types.PrimitiveTypeBytesReader[float64]}, &hashHistory, t)
	testHashInternal(testCase[complex64]{complex64(10), complex64(10), types.PrimitiveTypeBytesReader[complex64]}, &hashHistory, t)
	testHashInternal(testCase[complex64]{complex64(12), complex64(12), types.PrimitiveTypeBytesReader[complex64]}, &hashHistory, t)
	testHashInternal(testCase[complex128]{complex128(10), complex128(10), types.PrimitiveTypeBytesReader[complex128]}, &hashHistory, t)
	testHashInternal(testCase[complex128]{complex128(12), complex128(12), types.PrimitiveTypeBytesReader[complex128]}, &hashHistory, t)

	testStruct1 := testStruct{123, true, 0}
	testStruct1Copy := testStruct{123, true, 0}
	testStruct2 := testStruct{200, true, 0}
	testStruct2Copy := testStruct{200, true, 0}

	testHashInternal(testCase[testStruct]{testStruct1, testStruct1Copy, testStructBytesReader}, &hashHistory, t)
	testHashInternal(testCase[testStruct]{testStruct2, testStruct2Copy, testStructBytesReader}, &hashHistory, t)

	// Shouldn't collide since pointers differ
	testHashInternal(testCase[*testStruct]{&testStruct1, &testStruct1, types.PointerBytesReader[*testStruct]}, &hashHistory, t)
	testHashInternal(testCase[*testStruct]{&testStruct1Copy, &testStruct1Copy, types.PointerBytesReader[*testStruct]}, &hashHistory, t)
}

func testHashInternal[T any](tc testCase[T], hashHistory *[]hashedValue, t *testing.T) {
	// Compute hash
	bytes := tc.reader(tc.value)
	hash := Hash(bytes)

	// Check for collision
	for _, existingHash := range *hashHistory {
		if hash == existingHash.hash {
			t.Errorf("hash collision detected. value1=%#v (type=%T), value2=%#v (type=%T)", existingHash.value, existingHash.value, tc.value, tc.value)
			break
		}
	}

	// Computer again with value copy and compare
	bytes = tc.reader(tc.copy)
	hash2 := Hash(bytes)
	if hash != hash2 {
		t.Errorf("hash computation yielded different results for the same input. type=%T", tc.value)
	}

	// Store hash for collision detection
	*hashHistory = append(*hashHistory, hashedValue{tc.value, hash})
}
