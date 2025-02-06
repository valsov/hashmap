package hashmap

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	for _, entriesCount := range []int{100, 1000, 10_000, 100_000, 1_000_000} {
		key := strconv.Itoa(rand.Intn(entriesCount))

		b.Run(fmt.Sprintf("size_%d", entriesCount), func(b *testing.B) {
			b.Run("Hmap", func(b *testing.B) {
				m := New[string, int]()
				for i := range entriesCount {
					m.Set(fmt.Sprint(i), 0)
				}
				b.ResetTimer()

				for range b.N {
					for i := 0; i < b.N; i++ {
						_ = m.Get(key)
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
					for i := 0; i < b.N; i++ {
						_ = m[key]
					}
				}
			})
		})
	}
}
