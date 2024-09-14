package hashmap

import "testing"

func TestGet(t *testing.T) {
	testCases := []struct {
		values   []KeyValue[string, int]
		key      string
		expected int
	}{
		{},
	}
	for _, tc := range testCases {
		m := NewStringKeyMap[int]()
		for _, kv := range tc.values {
			m.Set(kv.Key, kv.Value)
		}

		value := m.Get(tc.key)
		if value != tc.expected {
			t.Errorf("retrieved invalid value. expected=%d, got=%d", tc.expected, value)
		}
	}
}

func TestTryGet(t *testing.T) {
	testCases := []struct {
		values        []KeyValue[string, int]
		key           string
		expected      int
		expectedFound bool
	}{
		{},
	}
	for _, tc := range testCases {
		m := NewStringKeyMap[int]()
		for _, kv := range tc.values {
			m.Set(kv.Key, kv.Value)
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

func TestLen(t *testing.T) {
	m := NewStringKeyMap[int]()

	// Initial
	len := m.Len()
	if len != 0 {
		t.Errorf("invalid length. expected=0, got=%d", len)
	}

	// Set 1 key
	m.Set("test1", 44)
	len = m.Len()
	if len != 1 {
		t.Errorf("invalid length. expected=1, got=%d", len)
	}

	// Set same key
	m.Set("test1", 55)
	len = m.Len()
	if len != 1 {
		t.Errorf("invalid length. expected=1, got=%d", len)
	}

	// Set another key
	m.Set("test2", 66)
	len = m.Len()
	if len != 1 {
		t.Errorf("invalid length. expected=2, got=%d", len)
	}

	// Remove key
	m.Delete("test1")
	len = m.Len()
	if len != 1 {
		t.Errorf("invalid length. expected=1, got=%d", len)
	}
}

func TestGetEntries(t *testing.T) {
	panic("todo")
}
