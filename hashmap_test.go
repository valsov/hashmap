package hashmap

import "testing"

type keyValue struct {
	key   string
	value int
}

func TestGet(t *testing.T) {
	testCases := []struct {
		values   []keyValue
		key      string
		expected int
	}{
		{},
	}
	for _, tc := range testCases {
		m := NewHmap[string, int]()
		for _, kv := range tc.values {
			m.Set(kv.key, kv.value)
		}

		value := m.Get(tc.key)
		if value != tc.expected {
			t.Errorf("retrieved invalid value. expected=%d, got=%d", tc.expected, value)
		}
	}
}

func TestTryGet(t *testing.T) {
	testCases := []struct {
		values        []keyValue
		key           string
		expected      int
		expectedFound bool
	}{
		{},
	}
	for _, tc := range testCases {
		m := NewHmap[string, int]()
		for _, kv := range tc.values {
			m.Set(kv.key, kv.value)
		}

		value, found := m.TryGet(tc.key)
		if value != tc.expected {
			t.Errorf("retrieved invalid value. expected=%d, got=%d", tc.expected, value)
		}
		if found != tc.expectedFound {
			t.Errorf("unexpected found state. expected=%t, got=%t", tc.expectedFound, found)
		}
	}
}

func TestSet(t *testing.T) {
	panic("todo")
}

func TestDelete(t *testing.T) {
	panic("todo")
}
