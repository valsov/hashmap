package hashmap

import (
	"github.com/valsov/hashmap/types"
)

type KeyValue[TKey, TValue any] struct {
	Key   TKey
	Value TValue
}

type Hashmap[TKey, TValue any] struct {
	keyBytesReader types.ReaderFunc[TKey]
	Length         int
}

func NewPrimitiveKeyMap[TKey types.Primitive, TValue any]() *Hashmap[TKey, TValue] {
	return &Hashmap[TKey, TValue]{keyBytesReader: types.PrimitiveTypeBytesReader[TKey]}
}

func NewStringKeyMap[TValue any]() *Hashmap[string, TValue] {
	return &Hashmap[string, TValue]{keyBytesReader: types.StringBytesReader}
}

func NewPointerKeyMap[TKey *TKeyVal, TValue any, TKeyVal any]() *Hashmap[TKey, TValue] {
	return &Hashmap[TKey, TValue]{keyBytesReader: types.PointerBytesReader[TKey]}
}

func NewMap[TKey comparable, TValue any](keyBytesReader types.ReaderFunc[TKey]) *Hashmap[TKey, TValue] {
	return &Hashmap[TKey, TValue]{keyBytesReader: keyBytesReader}
}

func (m Hashmap[TKey, TValue]) Get(key TKey) TValue {
	NewPrimitiveKeyMap[int, string]()
	NewStringKeyMap[byte]()
	NewPointerKeyMap[*int, string]()

	panic("todo")
}

func (m Hashmap[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	panic("todo")
}

func (m Hashmap[TKey, TValue]) Set(key TKey, value TValue) {
	panic("todo")
}

func (m Hashmap[TKey, TValue]) Delete(key TKey) {
	panic("todo")
}

func (m Hashmap[TKey, TValue]) Len() int {
	return m.Length
}

func (m Hashmap[TKey, TValue]) GetEntries() []KeyValue[TKey, TValue] {
	panic("todo")
}
