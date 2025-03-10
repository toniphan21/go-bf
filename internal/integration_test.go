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

func TestBloomFilter_Union_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000
	cf := v1.WithCapacity(uint32(m), 10)
	target := v1.Must(cf)
	other := v1.Must(cf)
	runNoFalseNegativeAfterUnionTest(t, n, target, other)
}

func TestBloomFilter_Union_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	n := 1_000_000
	cf := v1.WithAccuracy(0.001, uint32(n))
	target := v1.Must(cf)
	other := v1.Must(cf)
	runNoFalseNegativeAfterUnionTest(t, n, target, other)
}

func runNoFalseNegativeAfterUnionTest(t *testing.T, n int, target, other v1.BloomFilter) {
	targetKeys := make([]string, n)
	otherKeys := make([]string, n)

	for i := 0; i < n; i++ {
		targetKeys[i] = RandString(10)
		otherKeys[i] = RandString(12)

		target.Add([]byte(targetKeys[i]))
		other.Add([]byte(otherKeys[i]))

		if !target.Exists([]byte(targetKeys[i])) {
			t.Fatalf("Bloom Filter has false negative")
		}
		if !other.Exists([]byte(otherKeys[i])) {
			t.Fatalf("Bloom Filter has false negative")
		}
	}

	err := target.Union(other)
	if err != nil {
		t.Fatalf("Bloom Filter union failed")
	}
	for _, key := range targetKeys {
		if !target.Exists([]byte(key)) {
			t.Fatalf("Bloom Filter has false negative in target after Union() with key=%v", key)
		}
	}
	for _, key := range otherKeys {
		if !target.Exists([]byte(key)) {
			t.Fatalf("Bloom Filter has false negative when using key from other after Union() with key=%v", key)
		}
	}
}
