package internal

import (
	v1 "github.com/toniphan21/go-bf"
	v2 "github.com/toniphan21/go-bf/v2"
	"testing"
)

func BenchmarkBloomFilter_WithSHA_Add(b *testing.B) {
	const n = 10000
	inputs := make([][]byte, n)
	for i := 0; i < n; i++ {
		inputs[i] = []byte(RandString(10))
	}
	filterV1 := v1.Must(v1.WithAccuracy(0.01, 1_000_000), v1.WithSHA())
	filterV2 := v2.Must(v2.WithAccuracy(0.01, 1_000_000), v2.WithSHA())

	b.Run("v1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			populateBloomFilter(filterV1, inputs)
		}
	})

	b.Run("v2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			populateBloomFilter(filterV2, inputs)
		}
	})
}

func BenchmarkBloomFilter_WithFNV_Add(b *testing.B) {
	const n = 10000
	inputs := make([][]byte, n)
	for i := 0; i < n; i++ {
		inputs[i] = []byte(RandString(10))
	}
	filterV1 := v1.Must(v1.WithAccuracy(0.01, 1_000_000), v1.WithFNV())
	filterV2 := v2.Must(v2.WithAccuracy(0.01, 1_000_000), v2.WithFNV())

	b.Run("v1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			populateBloomFilter(filterV1, inputs)
		}
	})

	b.Run("v2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			populateBloomFilter(filterV2, inputs)
		}
	})
}

func BenchmarkBloomFilter_WithSHA_Exists(b *testing.B) {
	const n = 10000
	exists := make([][]byte, n)
	for i := 0; i < n; i++ {
		exists[i] = []byte(RandString(10))
	}
	notFounds := make([][]byte, n)
	for i := 0; i < n; i++ {
		notFounds[i] = []byte(RandString(11))
	}
	filterV1 := v1.Must(v1.WithAccuracy(0.01, 1_000_000), v1.WithSHA())
	filterV2 := v2.Must(v2.WithAccuracy(0.01, 1_000_000), v2.WithSHA())
	populateBloomFilter(filterV1, exists)
	populateBloomFilter(filterV2, exists)

	b.Run("v1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, input := range exists {
				filterV1.Exists(input)
			}
			for _, input := range notFounds {
				filterV1.Exists(input)
			}
		}
	})

	b.Run("v2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, input := range exists {
				filterV2.Exists(input)
			}
			for _, input := range notFounds {
				filterV2.Exists(input)
			}
		}
	})
}

func populateBloomFilter(filter BloomFilterSpec, inputs [][]byte) {
	for _, input := range inputs {
		filter.Add(input)
	}
}
