package hashmap

import "github.com/valsov/hashmap/types"

const defaultInitialCapacity uint = 128
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
type Hashmap[TKey comparable, TValue any] struct {
	storage    []mapEntry[TKey, TValue]
	length     int     // Number of entries in the hashmap
	loadFactor float32 // Load at which a storage growth will take place
	hashFunc   func(TKey) uint64
}

func NewPrimitiveKeyMap[TKey types.Primitive, TValue any]() *Hashmap[TKey, TValue] {
	return NewMap[TKey, TValue](types.PrimitiveTypeBytesReader[TKey])
}

func NewStringKeyMap[TValue any]() *Hashmap[string, TValue] {
	return NewMap[string, TValue](types.StringBytesReader)
}

func NewPointerKeyMap[TKey *TKeyVal, TValue any, TKeyVal any]() *Hashmap[TKey, TValue] {
	return NewMap[TKey, TValue](types.PointerBytesReader[TKey])
}

func NewMap[TKey comparable, TValue any](keyBytesReader types.ReaderFunc[TKey]) *Hashmap[TKey, TValue] {
	return &Hashmap[TKey, TValue]{
		storage:    make([]mapEntry[TKey, TValue], defaultInitialCapacity),
		hashFunc:   createHashFunc(keyBytesReader),
		loadFactor: defaultLoadFactor,
	}
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
			return
		}

		curSlotIdealIndex := m.getIdealKeyIndex(m.storage[index].key)
		curSlotDistance := (index + len(m.storage) - curSlotIdealIndex) % len(m.storage)
		if distance > curSlotDistance {
			// Insert data in this slot and continue to find a new spot for the previous data
			m.storage[index].key, key = key, m.storage[index].key
			m.storage[index].value, value = value, m.storage[index].value
			distance = curSlotDistance
		}
		distance++
		index = (index + 1) % len(m.storage)
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
	index = (index + 1) % len(m.storage)
	for {
		if !m.storage[index].alive {
			m.emptySlot(previousIndex)
			return
		}

		curSlotIdealIndex := m.getIdealKeyIndex(m.storage[index].key)
		curSlotDistance := (index + len(m.storage) - curSlotIdealIndex) % len(m.storage)
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
		index = (index + 1) % len(m.storage)
	}
}

// Remove all entries from the hashmap.
func (m *Hashmap[TKey, TValue]) Clear() {
	m.length = 0
	m.storage = make([]mapEntry[TKey, TValue], len(m.storage))
}

// Get the number of entries stored in the hashmap.
func (m *Hashmap[TKey, TValue]) Len() int {
	return int(m.length)
}

// Get all entries stored in the hashmap.
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
	for {
		if m.storage[index].key == key {
			return index, true
		}

		if !m.storage[index].alive {
			return 0, false
		}

		index = (index + 1) % len(m.storage)
	}
}

// Compute the index at which the given key should be located.
func (m *Hashmap[TKey, TValue]) getIdealKeyIndex(key TKey) int {
	hash := m.hashFunc(key)
	return int(hash % uint64(len(m.storage)))
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
