package hashmap

type Hmap[TKey comparable, TValue any] struct{}

func NewHmap[TKey comparable, TValue any]() Hmap[TKey, TValue] {
	return Hmap[TKey, TValue]{}
}

func (m *Hmap[TKey, TValue]) Get(key TKey) TValue {
	panic("todo")
}

func (m *Hmap[TKey, TValue]) TryGet(key TKey) (TValue, bool) {
	panic("todo")
}

func (m *Hmap[TKey, TValue]) Set(key TKey, value TValue) {
	panic("todo")
}

func (m *Hmap[TKey, TValue]) Delete(key TKey) {
	panic("todo")
}
