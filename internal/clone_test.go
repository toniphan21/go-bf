package internal

import (
	v1 "github.com/toniphan21/go-bf"
	v2 "github.com/toniphan21/go-bf/v2"
	"testing"
)

func TestBloomFilter_Clone_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000

	v1Cf := v1.WithCapacity(uint32(m), 10)
	v1Filter, _ := v1.New(v1Cf)
	runCloneNoFalseNegativeTest(t, v1Filter, n, func() (BloomFilterSpec, error) {
		return v1Filter.Clone()
	})

	v2Cf := v2.WithCapacity(uint32(m), 10)
	v2Filter, _ := v2.New(v2Cf)
	runCloneNoFalseNegativeTest(t, v2Filter, n, func() (BloomFilterSpec, error) {
		return v2Filter.Clone(), nil
	})
}

func TestBloomFilter_Clone_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	var n = 1_000_000

	v1Cf := v1.WithAccuracy(0.001, uint32(n))
	v1Filter, _ := v1.New(v1Cf)
	runCloneNoFalseNegativeTest(t, v1Filter, n, func() (BloomFilterSpec, error) {
		return v1Filter.Clone()
	})

	v2Cf := v2.WithAccuracy(0.001, uint32(n))
	v2Filter, _ := v2.New(v2Cf)
	runCloneNoFalseNegativeTest(t, v2Filter, n, func() (BloomFilterSpec, error) {
		return v2Filter.Clone(), nil
	})
}

func runCloneNoFalseNegativeTest(t *testing.T, filter BloomFilterSpec, n int, fn func() (BloomFilterSpec, error)) {
	keys := make([][]byte, n)
	for i := 0; i < n; i++ {
		item := []byte(RandString(10))
		keys[i] = item
		filter.Add(item)
		after := filter.Exists(item)
		if !after {
			t.Fatalf("Bloom Filter has false negative")
		}
	}

	cloned, _ := fn()
	for i := 0; i < n; i++ {
		check := cloned.Exists(keys[i])
		if !check {
			t.Fatalf("Bloom Filter has false negative")
		}
	}
	if cloned.Count() != filter.Count() {
		t.Fatalf("cloned Bloom Filter has different count")
	}
}
