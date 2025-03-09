package internal

import (
	v1 "github.com/toniphan21/go-bf"
	v2 "github.com/toniphan21/go-bf/v2"
	"testing"
)

func TestBloomFilter_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000
	runNoFalseNegativeTest(t, v1.Must(v1.WithCapacity(uint32(m), 10)), n)
	runNoFalseNegativeTest(t, v2.Must(v2.WithCapacity(uint32(m), 10)), n)
}

func TestBloomFilter_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	var n = 1_000_000
	runNoFalseNegativeTest(t, v1.Must(v1.WithAccuracy(0.001, uint32(n))), n)
	runNoFalseNegativeTest(t, v2.Must(v2.WithAccuracy(0.001, uint32(n))), n)
}

func TestBloomFilter_NoFalseNegative_WithCapacity_Expandable(t *testing.T) {
	t.Parallel()
	n := 1_000_000
	runNoFalseNegativeTest(t, v2.Must(v2.WithCapacity(uint32(n/10), 5)), n)
}

func TestBloomFilter_NoFalseNegative_WithAccuracy_Expandable(t *testing.T) {
	t.Parallel()
	var n = 1_000_000
	runNoFalseNegativeTest(t, v2.Must(v2.WithAccuracy(0.001, uint32(n/10))), n)
}

func runNoFalseNegativeTest(t *testing.T, filter BloomFilterSpec, n int) {
	for i := 0; i < n; i++ {
		item := []byte(RandString(10))
		filter.Add(item)
		after := filter.Exists(item)
		if !after {
			t.Fatalf("Bloom Filter has false negative")
		}
	}
}
