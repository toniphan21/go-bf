package internal

import (
	v1 "github.com/toniphan21/go-bf"
	"testing"
)

func TestBloomFilter_Clone_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000
	cf := v1.WithCapacity(uint32(m), 10)
	filter, _ := v1.New(cf)
	runCloneNoFalseNegativeTest(t, filter, n)
}

func TestBloomFilter_Clone_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	var n = 1_000_000
	cf := v1.WithAccuracy(0.001, uint32(n))
	filter, _ := v1.New(cf)
	runCloneNoFalseNegativeTest(t, filter, n)
}

func runCloneNoFalseNegativeTest(t *testing.T, filter v1.BloomFilter, n int) {
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

	cloned, _ := filter.Clone()
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
