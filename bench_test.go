package bf

import (
	"fmt"
	"testing"
)

func BenchmarkBloomFilter_WithSHA_Add(b *testing.B) {
	bf, _ := New(WithAccuracy(0.01, 1_000_000), WithSHA())
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkBloomFilter_WithFNV_Add(b *testing.B) {
	bf, _ := New(WithAccuracy(0.01, 1_000_000), WithFNV())
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkBloomFilter_WithSHA_Exists(b *testing.B) {
	bf, _ := New(WithAccuracy(0.01, 1_000_000), WithSHA())
	for i := 0; i < b.N; i++ {
		bf.Exists([]byte(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkBloomFilter_WithFNV_Exists(b *testing.B) {
	bf, _ := New(WithAccuracy(0.01, 1_000_000), WithFNV())
	for i := 0; i < b.N; i++ {
		bf.Exists([]byte(fmt.Sprintf("%d", i)))
	}
}
