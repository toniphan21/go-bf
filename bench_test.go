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

func initForBench(count, m, n int) []BloomFilter {
	result := make([]BloomFilter, count)
	for i := 0; i < count; i++ {
		result[i] = Must(WithAccuracy(0.01, uint32(m)))
		for j := 0; j < n; j++ {
			result[i].Add([]byte(internal.RandString(10)))
		}
	}
	return result
}

func BenchmarkBloomFilter_Intersect(b *testing.B) {
	bfs := initForBench(2, 1_000_000, 1000)
	b.Run("bench", func(pb *testing.B) {
		for i := 0; i < pb.N; i++ {
			_ = bfs[0].Intersect(bfs[1])
		}
	})
}

func BenchmarkBloomFilter_Union(b *testing.B) {
	bfs := initForBench(2, 1_000_000, 1000)
	b.Run("bench", func(pb *testing.B) {
		for i := 0; i < pb.N; i++ {
			_ = bfs[0].Intersect(bfs[1])
		}
	})
}

func BenchmarkBloomFilter_Clone(b *testing.B) {
	bfs := initForBench(1, 1_000_000, 1000)
	b.Run("bench", func(pb *testing.B) {
		for i := 0; i < pb.N; i++ {
			_, _ = bfs[0].Clone()
		}
	})
}
