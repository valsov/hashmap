package hashmap

import (
	"slices"
	"testing"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		values   []KeyValue[string, int]
		key      string
		expected int
	}{
		{
			values:   []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:      "notfound",
			expected: 0,
		},
		{
			values:   []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:      "key1",
			expected: 123,
		},
		{
			values:   []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:      "key2",
			expected: 456,
		},
	}
	for _, tc := range testCases {
		m := New[string, int]()
		for _, kv := range tc.values {
			m.Set(kv.Key, kv.Value)
		}

		value := m.Get(tc.key)
		if value != tc.expected {
			t.Errorf("retrieved invalid value for key=%s. expected=%d, got=%d", tc.key, tc.expected, value)
		}
	}
}

func TestTryGet(t *testing.T) {
	testCases := []struct {
		values        []KeyValue[string, int]
		key           string
		expectedValue int
		expectedFound bool
	}{
		{
			values:        []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:           "notfound",
			expectedFound: false,
		},
		{
			values:        []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:           "key1",
			expectedValue: 123,
			expectedFound: true,
		},
		{
			values:        []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:           "key2",
			expectedValue: 456,
			expectedFound: true,
		},
	}
	for _, tc := range testCases {
		m := New[string, int]()
		for _, kv := range tc.values {
			m.Set(kv.Key, kv.Value)
		}

		value, found := m.TryGet(tc.key)
		if found != tc.expectedFound {
			t.Errorf("unexpected found state for key=%s. expected=%t, got=%t", tc.key, tc.expectedFound, found)
		}
		if value != tc.expectedValue {
			t.Errorf("retrieved invalid value for key=%s. expected=%d, got=%d", tc.key, tc.expectedValue, value)
		}
	}
}

func TestSet(t *testing.T) {
	testCases := []struct {
		values   []KeyValue[string, int]
		keyVal   KeyValue[string, int]
		expected []KeyValue[string, int]
	}{
		{
			values:   []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			keyVal:   KeyValue[string, int]{"key1", 99999},
			expected: []KeyValue[string, int]{{"key1", 99999}, {"key2", 456}},
		},
		{
			values:   []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			keyVal:   KeyValue[string, int]{"key3", 99999},
			expected: []KeyValue[string, int]{{"key1", 123}, {"key2", 456}, {"key3", 99999}},
		},
		{
			values:   []KeyValue[string, int]{},
			keyVal:   KeyValue[string, int]{"key1", 123},
			expected: []KeyValue[string, int]{{"key1", 123}},
		},
	}
	for _, tc := range testCases {
		m := New[string, int]()
		for _, kv := range tc.values {
			m.Set(kv.Key, kv.Value)
		}

		m.Set(tc.keyVal.Key, tc.keyVal.Value)

		for _, kv := range tc.expected {
			value, found := m.TryGet(kv.Key)
			if !found {
				t.Errorf("key=%s not found", kv.Key)
			}
			if value != kv.Value {
				t.Errorf("retrieved invalid value for key=%s. expected=%d, got=%d", kv.Key, kv.Value, value)
			}
		}
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		values   []KeyValue[string, int]
		key      string
		expected []KeyValue[string, int]
	}{
		{
			values:   []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:      "notfound",
			expected: []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
		},
		{
			values:   []KeyValue[string, int]{{"key1", 123}, {"key2", 456}},
			key:      "key1",
			expected: []KeyValue[string, int]{{"key2", 456}},
		},
		{
			values:   []KeyValue[string, int]{{"key1", 123}},
			key:      "key1",
			expected: []KeyValue[string, int]{},
		},
	}
	for _, tc := range testCases {
		m := New[string, int]()
		for _, kv := range tc.values {
			m.Set(kv.Key, kv.Value)
		}

		m.Delete(tc.key)

		if _, found := m.TryGet(tc.key); found {
			t.Errorf("key=%s was found", tc.key)
		}

		for _, kv := range tc.expected {
			value, found := m.TryGet(kv.Key)
			if !found {
				t.Errorf("key=%s not found", kv.Key)
			}
			if value != kv.Value {
				t.Errorf("retrieved invalid value for key=%s. expected=%d, got=%d", kv.Key, kv.Value, value)
			}
		}
	}
}

func TestClear(t *testing.T) {
	testCases := [][]KeyValue[string, int]{
		{},
		{{"key1", 123}},
		{{"key1", 123}, {"key2", 456}, {"key3", 789}},
	}
	for _, tc := range testCases {
		m := New[string, int]()
		for _, kv := range tc {
			m.Set(kv.Key, kv.Value)
		}

		m.Clear()
		len := m.Len()
		if len != 0 {
			t.Errorf("invalid length. expected=0, got=%d", len)
		}
		for _, kv := range tc {
			if _, found := m.TryGet(kv.Key); found {
				t.Errorf("key=%s was found", kv.Key)
			}
		}
	}
}

func TestLen(t *testing.T) {
	m := New[string, int]()

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
	if len != 2 {
		t.Errorf("invalid length. expected=2, got=%d", len)
	}

	// Remove key
	m.Delete("test1")
	len = m.Len()
	if len != 1 {
		t.Errorf("invalid length. expected=1, got=%d", len)
	}

	// Remove last key
	m.Delete("test2")
	len = m.Len()
	if len != 0 {
		t.Errorf("invalid length. expected=0, got=%d", len)
	}

	// Remove non-existant key
	m.Delete("test2")
	len = m.Len()
	if len != 0 {
		t.Errorf("invalid length. expected=0, got=%d", len)
	}
}

func TestGetEntries(t *testing.T) {
	testCases := [][]KeyValue[string, int]{
		{},
		{{"key1", 123}},
		{{"key1", 123}, {"key2", 456}, {"key3", 789}},
	}
	for _, tc := range testCases {
		m := New[string, int]()
		for _, kv := range tc {
			m.Set(kv.Key, kv.Value)
		}

		entries := m.GetEntries()
		if len(entries) != len(tc) {
			t.Errorf("invalid length. expected=%d, got=%d", len(tc), len(entries))
		}

		if len(entries) == 0 {
			continue
		}

		slices.SortFunc(entries, func(a, b KeyValue[string, int]) int {
			return a.Value - b.Value
		})

		for i := 0; i < len(entries); i++ {
			if entries[i].Key != tc[i].Key {
				t.Errorf("invalid key. expected=%s got=%s", tc[i].Key, entries[i].Key)
			}
			if entries[i].Value != tc[i].Value {
				t.Errorf("invalid value for key=%s. expected=%d got=%d", tc[i].Key, tc[i].Value, entries[i].Value)
			}
		}
	}
}
