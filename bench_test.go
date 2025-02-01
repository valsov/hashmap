package hashmap

import (
	"fmt"
	"testing"
)

var mapSizes []int = []int{100, 1000, 10_000, 100_000, 1_000_000}

func BenchmarkGet(b *testing.B) {
	for _, entriesCount := range mapSizes {
		b.Run(fmt.Sprintf("size_%d", entriesCount), func(b *testing.B) {
			b.Run("Hmap", func(b *testing.B) {
				m := New[string, int]()
				for i := range entriesCount {
					m.Set(fmt.Sprint(i), 0)
				}
				b.ResetTimer()

				for range b.N {
					for key := 0; key < b.N; key++ {
						_ = m.Get(fmt.Sprint(key))
					}
				}
			})
			b.Run("Native map", func(b *testing.B) {
				m := map[string]int{}
				for i := range entriesCount {
					m[fmt.Sprint(i)] = 0
				}
				b.ResetTimer()

				for range b.N {
					for key := 0; key < b.N; key++ {
						_ = m[fmt.Sprint(key)]
					}
				}
			})
		})
	}
}

func BenchmarkTryGet(b *testing.B) {
	for _, entriesCount := range mapSizes {
		var found bool
		b.Run(fmt.Sprintf("size_%d", entriesCount), func(b *testing.B) {
			b.Run("Hmap", func(b *testing.B) {
				m := New[string, int]()
				for i := range entriesCount {
					m.Set(fmt.Sprint(i), 0)
				}
				b.ResetTimer()

				for range b.N {
					_, found = m.TryGet("50")
				}
			})
			b.Run("Native map", func(b *testing.B) {
				m := map[string]int{}
				for i := range entriesCount {
					m[fmt.Sprint(i)] = 0
				}
				b.ResetTimer()

				for range b.N {
					_, found = m["50"]
				}
			})
		})
		// Dummy code to avoid compilation warnings
		if found {
			found = false
		}
	}
}
