# Hashmap

This project implements a hashmap with the Robin Hood hashing algorithm (see [paper](https://cs.uwaterloo.ca/research/tr/1986/CS-86-14.pdf)).

The default key hash function is Go's standard library implementation. It is extracted at run time using an unsafe pointer to the native `map`. Refer to the source source of `hasher.go` for more details.

## Usage

```go
// Init
m := hashmap.New[string, int]() // Equivalent to make(map[string]int)

// Write
m.Set("key", 123)

// Read
value := m.Get("key") // value == 123

// Try read
value, found := m.TryGet("key") // value == 123, found == true
value, found := m.TryGet("key2") // value == 0, found == false

// Get stored entries count
length := m.Len() // length == 1

// Get all stored key-value pairs
entries := m.GetEntries() // entries == []KeyValue[string, int]{{Key: "key", Value: 123}}

// Delete a specific entry
m.Delete("key")

// Clear the entire map
m.Clear()
```

## Configuration

The hashmap is pre-configured with the following values, but they can be configured:
- Key hasher: **Go's native implementation**.
- Initial capacity: **128 entries**.
- Load factor: **50%**. The load factor is the maximum hashmap load at which point it will resize itself at double its size.

```go
// Custom key hasher
hasher := func(keyPtr, seed uintptr) uintptr {
    // [...]
    return generatedHash
}

// Specify initial capacity (must be a power of 2)
initialCap := 256

// With an increased load factor
loadFactor := 70

// Create configured map
m := hashmap.New[string, int](
    hashmap.WithHashFunc[string, int](hasher),
    hashmap.WithInitialCapacity[string, int](initialCap),
    hashmap.WithMaxLoadPercentage[string, int](loadFactor),
)
```
