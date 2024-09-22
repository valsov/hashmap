package hasher

import (
	"testing"
	"unsafe"
)

type testCase[T any] struct {
	value T
	copy  T // Used to create another hash and comparing it with the one computed from `value`
}

// Store both a value and its hash
type hashedValue struct {
	value any
	hash  uint64
}

// Struct to check struct and pointer hashing
type testStruct struct {
	A int
	B bool
	C uint64
}

func TestGetHashFunc(t *testing.T) {
	hashHistory := []hashedValue{}

	testGetHashFuncInternal(testCase[string]{"str1", "str1"}, &hashHistory, t)
	testGetHashFuncInternal(testCase[string]{"str2", "str2"}, &hashHistory, t)
	testGetHashFuncInternal(testCase[rune]{'a', 'a'}, &hashHistory, t)
	testGetHashFuncInternal(testCase[rune]{'b', 'b'}, &hashHistory, t)
	testGetHashFuncInternal(testCase[byte]{byte(128), byte(128)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[byte]{byte(120), byte(120)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int8]{int8(10), int8(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int8]{int8(11), int8(11)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint8]{uint8(12), uint8(12)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint8]{uint8(13), uint8(13)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int16]{int16(10), int16(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int16]{int16(11), int16(11)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint16]{uint16(12), uint16(12)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint16]{uint16(13), uint16(13)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int32]{int32(10), int32(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int32]{int32(11), int32(11)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint32]{uint32(12), uint32(12)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint32]{uint32(13), uint32(13)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int64]{int64(10), int64(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int64]{int64(11), int64(11)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint64]{uint64(12), uint64(12)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint64]{uint64(13), uint64(13)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int]{int(14), int(14)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[int]{int(15), int(15)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint]{uint(16), uint(16)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uint]{uint(17), uint(17)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uintptr]{uintptr(18), uintptr(18)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[uintptr]{uintptr(19), uintptr(19)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[float32]{float32(10), float32(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[float32]{float32(12), float32(12)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[float64]{float64(10), float64(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[float64]{float64(12), float64(12)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[complex64]{complex64(10), complex64(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[complex64]{complex64(12), complex64(12)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[complex128]{complex128(10), complex128(10)}, &hashHistory, t)
	testGetHashFuncInternal(testCase[complex128]{complex128(12), complex128(12)}, &hashHistory, t)

	testStruct1 := testStruct{123, true, 0}
	testStruct1Copy := testStruct{123, true, 0}
	testStruct2 := testStruct{200, true, 0}
	testStruct2Copy := testStruct{200, true, 0}

	testGetHashFuncInternal(testCase[testStruct]{testStruct1, testStruct1Copy}, &hashHistory, t)
	testGetHashFuncInternal(testCase[testStruct]{testStruct2, testStruct2Copy}, &hashHistory, t)

	// Shouldn't collide since pointers differ
	testGetHashFuncInternal(testCase[*testStruct]{&testStruct1, &testStruct1}, &hashHistory, t)
	testGetHashFuncInternal(testCase[*testStruct]{&testStruct1Copy, &testStruct1Copy}, &hashHistory, t)
}

func testGetHashFuncInternal[T comparable](tc testCase[T], hashHistory *[]hashedValue, t *testing.T) {
	// Compute hash
	hasher := GetHashFunc[T]()
	seed := GenerateSeed()
	if hasher == nil {
		t.Errorf("hasher retrieval failed. type=%T", tc.value)
		return
	}

	hash := uint64(hasher(uintptr(unsafe.Pointer(&tc.value)), seed))

	// Check for collision
	for _, existingHash := range *hashHistory {
		if hash == existingHash.hash {
			t.Errorf("hash collision detected. value1=%#v (type=%T), value2=%#v (type=%T)", existingHash.value, existingHash.value, tc.value, tc.value)
			break
		}
	}

	// Computer hash again with value copy and compare
	hash2 := uint64(hasher(uintptr(unsafe.Pointer(&tc.copy)), seed))
	if hash != hash2 {
		t.Errorf("hash computation yielded different results for the same input. type=%T", tc.value)
	}

	// Store hash for collision detection
	*hashHistory = append(*hashHistory, hashedValue{tc.value, hash})
}
