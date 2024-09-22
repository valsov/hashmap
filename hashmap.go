package hashmap

import (
	"unsafe"

	"github.com/valsov/hashmap/hasher"
)

const defaultInitialCapacity uint = 128 // Power of 2
const defaultLoadFactor float32 = 0.5

// Key value pair
type KeyValue[TKey, TValue any] struct {
	Key   TKey
	Value TValue
}

// Internal storage unit of a key value pair
type mapEntry[TKey, TValue any] struct {
	key   TKey
	value TValue
	alive bool
}

// Hashmap struct for fast data lookup
//
// The capacity of the hashmap must be a power of 2. This allows to do: hash & (cap - 1) to compute indexes.
// This way, the use of modulo operator is avoided (which is a much slower operation compared to bitwise AND).
type Hashmap[TKey comparable, TValue any] struct {
	storage    []mapEntry[TKey, TValue]
	length     int     // Number of entries in the hashmap
	loadFactor float32 // Load at which a storage growth will take place
	maxProbe   int     // The maximum number of slots a key search should check, this is the max distance an entry was placed from its ideal index
	hashFunc   func(uintptr, uintptr) uintptr
	hashSeed   uintptr
}

// Instanciate a new hashmap with a custom key bytes reader function.
func New[TKey comparable, TValue any](config ...HashMapConfig[TKey, TValue]) *Hashmap[TKey, TValue] {
	m := &Hashmap[TKey, TValue]{
		loadFactor: defaultLoadFactor,
		hashSeed:   hasher.GenerateSeed(),
	}
	for _, configFunc := range config {
		configFunc(m)
	}

	if m.storage == nil {
		m.storage = make([]mapEntry[TKey, TValue], defaultInitialCapacity)
	}
	if m.hashFunc == nil {
		m.hashFunc = hasher.GetHashFunc[TKey]()
	}

	return m
}

// Get the value associated with the given key. A default value is returned if the key doesn't exist.
func (m *Hashmap[TKey, TValue]) Get(key TKey) TValue {
	index, found := m.tryGetKeyIndex(key)
	if found {
		return m.storage[index].value
	}
	var zeroEntry TValue
	return zeroEntry
}

// Try to get the value associated with the given key.
func (m *Hashmap[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	index, found := m.tryGetKeyIndex(key)
	if found {
		return m.storage[index].value, true
	}
	var zeroEntry TValue
	return zeroEntry, false
}

// Insert or update the given value at the given key.
func (m *Hashmap[TKey, TValue]) Set(key TKey, value TValue) {
	if float64(m.length) >= float64(len(m.storage))*float64(m.loadFactor) {
		m.grow()
	}
	m.length++

	// Find suitable slot
	index := m.getIdealKeyIndex(key)
	var distance int
	for {
		if m.storage[index].key == key {
			m.storage[index].value = value
			m.length-- // Replacement
			return
		}

		if !m.storage[index].alive {
			m.storage[index].key = key
			m.storage[index].value = value
			m.storage[index].alive = true

			m.maxProbe = max(m.maxProbe, distance)
			return
		}

		curSlotIdealIndex := m.getIdealKeyIndex(m.storage[index].key)
		curSlotDistance := (index + len(m.storage) - curSlotIdealIndex) & (len(m.storage) - 1)
		if distance > curSlotDistance {
			// Insert data in this slot and continue to find a new spot for the previous data
			m.storage[index].key, key = key, m.storage[index].key
			m.storage[index].value, value = value, m.storage[index].value

			m.maxProbe = max(m.maxProbe, distance)
			distance = curSlotDistance
		}
		distance++
		index = (index + 1) & (len(m.storage) - 1)
	}
}

// Remove the entry with the given key from the hashmap.
func (m *Hashmap[TKey, TValue]) Delete(key TKey) {
	// Find entry
	index, found := m.tryGetKeyIndex(key)
	if !found {
		return
	}
	m.length--

	// Move next entries if necessary
	previousIndex := index
	index = (index + 1) & (len(m.storage) - 1)
	for {
		if !m.storage[index].alive {
			m.emptySlot(previousIndex)
			return
		}

		curSlotIdealIndex := m.getIdealKeyIndex(m.storage[index].key)
		curSlotDistance := (index + len(m.storage) - curSlotIdealIndex) & (len(m.storage) - 1)
		if curSlotDistance == 0 {
			// Ideal placement
			m.emptySlot(previousIndex)
			return
		}

		// Shift entry one slot back
		m.storage[previousIndex].key = m.storage[index].key
		m.storage[previousIndex].value = m.storage[index].value
		m.storage[previousIndex].alive = true

		previousIndex = index
		index = (index + 1) & (len(m.storage) - 1)
	}
}

// Remove all entries from the hashmap.
func (m *Hashmap[TKey, TValue]) Clear() {
	m.storage = make([]mapEntry[TKey, TValue], len(m.storage))
	m.length = 0
	m.maxProbe = 0
}

// Get the number of entries stored in the hashmap.
func (m *Hashmap[TKey, TValue]) Len() int {
	return int(m.length)
}

// Get all entries stored in the hashmap.
//
// The slice ordering is not guaranteed to be the insertion order.
func (m *Hashmap[TKey, TValue]) GetEntries() []KeyValue[TKey, TValue] {
	entries := make([]KeyValue[TKey, TValue], m.length)
	index := 0
	for _, entry := range m.storage {
		if entry.alive {
			entries[index] = KeyValue[TKey, TValue]{
				Key:   entry.key,
				Value: entry.value,
			}
			index++
		}
	}
	return entries
}

// Main lookup function, try to find the index of the given key.
func (m *Hashmap[TKey, TValue]) tryGetKeyIndex(key TKey) (int, bool) {
	index := m.getIdealKeyIndex(key)
	// The value can only be located within a range of m.maxProbe from its ideal index
	for range m.maxProbe + 1 {
		if m.storage[index].key == key {
			return index, true
		}

		if !m.storage[index].alive {
			return 0, false
		}

		index = (index + 1) & (len(m.storage) - 1)
	}
	return 0, false
}

// Compute the index at which the given key should be located.
func (m *Hashmap[TKey, TValue]) getIdealKeyIndex(key TKey) int {
	hash := uint64(m.hashFunc(uintptr(unsafe.Pointer(&key)), m.hashSeed))
	return int(hash & uint64(len(m.storage)-1))
}

// Set the slot's value to the default, dead slot.
func (m *Hashmap[TKey, TValue]) emptySlot(index int) {
	var zeroVal mapEntry[TKey, TValue]
	m.storage[index] = zeroVal
}

// Allocate a new storage slice, twice as big as previous storage.
// Entries from the previous storage are put into the new storage.
func (m *Hashmap[TKey, TValue]) grow() {
	oldStorage := m.storage
	m.storage = make([]mapEntry[TKey, TValue], len(m.storage)*2)
	m.length = 0 // Reset length, it will be set by m.Set()

	for _, entry := range oldStorage {
		if entry.alive {
			m.Set(entry.key, entry.value)
		}
	}
}
