package hashmap

// Configuration function to customize internal properties of a Hashmap.
type HashMapConfig[TKey comparable, TValue any] func(*Hashmap[TKey, TValue])

// Specify a custom load percentage, beyond which the hashmap will grow in size.
//
// This must be less than 100, or the default percentage will be applied.
func WithMaxLoadPercentage[TKey comparable, TValue any](loadPercentage uint) HashMapConfig[TKey, TValue] {
	var loadFactor float32
	if loadPercentage >= 100 {
		loadFactor = defaultLoadFactor
	} else {
		loadFactor = float32(loadPercentage) / 100
	}

	return func(hmap *Hashmap[TKey, TValue]) {
		hmap.loadFactor = loadFactor
	}
}

// Specify an initial entries storage capacity.
//
// This must be a power of 2, or the default initial capacity will be applied.
func WithInitialCapacity[TKey comparable, TValue any](initialCapacity uint) HashMapConfig[TKey, TValue] {
	if initialCapacity == 0 || initialCapacity%2 != 0 {
		initialCapacity = defaultInitialCapacity
	}

	return func(hmap *Hashmap[TKey, TValue]) {
		hmap.storage = make([]mapEntry[TKey, TValue], initialCapacity)
	}
}

// Specify a custom hash function that will be used for key hashing operations.
func WithHashFunc[TKey comparable, TValue any](hashFunc func(uintptr, uintptr) uintptr) HashMapConfig[TKey, TValue] {
	return func(hmap *Hashmap[TKey, TValue]) {
		hmap.hashFunc = hashFunc
	}
}
