package bf

import (
	"fmt"
	"github.com/toniphan21/go-bf/internal"
	"testing"
)

func BenchmarkBloomFilter_WithSHA_Add(b *testing.B) {
	bf := Must(WithAccuracy(0.01, 1_000_000), WithSHA())
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkBloomFilter_WithFNV_Add(b *testing.B) {
	bf := Must(WithAccuracy(0.01, 1_000_000), WithFNV())
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkBloomFilter_WithSHA_Exists(b *testing.B) {
	bf := Must(WithAccuracy(0.01, 1_000_000), WithSHA())
	for i := 0; i < b.N; i++ {
		bf.Exists([]byte(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkBloomFilter_WithFNV_Exists(b *testing.B) {
	bf := Must(WithAccuracy(0.01, 1_000_000), WithFNV())
	for i := 0; i < b.N; i++ {
		bf.Exists([]byte(fmt.Sprintf("%d", i)))
	}
}

func BenchmarkBloomFilter_Intersect(b *testing.B) {
	bf := Must(WithAccuracy(0.01, 1_000_000))
	right := Must(WithAccuracy(0.01, 1_000_000))
	for i := 0; i < 1000; i++ {
		bf.Add([]byte(internal.RandString(10)))
		right.Add([]byte(internal.RandString(10)))
	}
	b.Run("bench", func(pb *testing.B) {
		for i := 0; i < pb.N; i++ {
			_ = bf.Intersect(right)
		}
	})
}

func BenchmarkBloomFilter_Union(b *testing.B) {
	bf := Must(WithAccuracy(0.01, 1_000_000))
	right := Must(WithAccuracy(0.01, 1_000_000))
	for i := 0; i < 1000; i++ {
		bf.Add([]byte(internal.RandString(10)))
		right.Add([]byte(internal.RandString(10)))
	}
	b.Run("bench", func(pb *testing.B) {
		for i := 0; i < pb.N; i++ {
			_ = bf.Union(right)
		}
	})
}
