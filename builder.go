package hashmap

import (
	"github.com/valsov/hashmap/hasher"
	"github.com/valsov/hashmap/types"
)

// Helper struct to configure and build a hashmap.
type Builder[TKey comparable, TValue any] struct {
	hashFunc        func(TKey) uint64
	initialCapacity uint
	loadFactor      float32
}

// Specify a custom load percentage, beyond which the hashmap will grow in size.
//
// This must be less than 100, or the default percentage will be applied.
func (b Builder[TKey, TValue]) WithMaxLoadPercentage(loadPercentage uint) Builder[TKey, TValue] {
	if loadPercentage >= 100 {
		b.loadFactor = defaultLoadFactor
	} else {
		b.loadFactor = float32(loadPercentage / 100)
	}
	return b
}

// Specify an initial entries storage capacity.
//
// This must be a power of 2, or the default initial capacity will be applied.
func (b Builder[TKey, TValue]) WithInitialCapacity(initialCapacity uint) Builder[TKey, TValue] {
	if initialCapacity == 0 || initialCapacity%2 != 0 {
		b.initialCapacity = defaultInitialCapacity
	} else {
		b.initialCapacity = initialCapacity
	}
	return b
}

// Produce a hashmap configured with the builder's properties.
func (b Builder[TKey, TValue]) Build() *Hashmap[TKey, TValue] {
	return &Hashmap[TKey, TValue]{
		storage:    make([]mapEntry[TKey, TValue], b.initialCapacity),
		hashFunc:   b.hashFunc,
		loadFactor: b.loadFactor,
	}
}

// Get a hashmap builder with a primitive key type.
func NewPrimitiveKeyMapBuilder[TKey types.Primitive, TValue any]() Builder[TKey, TValue] {
	return NewMapBuilder[TKey, TValue](types.PrimitiveTypeBytesReader[TKey])
}

// Get a hashmap builder with a string key type.
func NewStringKeyMapBuilder[TValue any]() Builder[string, TValue] {
	return NewMapBuilder[string, TValue](types.StringBytesReader)
}

// Get a hashmap builder with a pointer key type.
func NewPointerKeyMapBuilder[TKey *TKeyVal, TValue any, TKeyVal any]() Builder[TKey, TValue] {
	return NewMapBuilder[TKey, TValue](types.PointerBytesReader[TKey])
}

// Get a hashmap builder with a custom key bytes reader function.
func NewMapBuilder[TKey comparable, TValue any](keyBytesReader types.ReaderFunc[TKey]) Builder[TKey, TValue] {
	return Builder[TKey, TValue]{
		hashFunc: createHashFunc(keyBytesReader),
	}
}

// Helper to generate a hashing function from the given bytes reader function.
func createHashFunc[TKey comparable](readerFunc types.ReaderFunc[TKey]) func(TKey) uint64 {
	return func(key TKey) uint64 {
		bytes := readerFunc(key)
		return hasher.Hash(bytes)
	}
}
