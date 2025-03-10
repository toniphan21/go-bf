package internal

import (
	v1 "github.com/toniphan21/go-bf"
	v2 "github.com/toniphan21/go-bf/v2"
	"testing"
)

func TestBloomFilter_Union_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000

	v1Cf := v1.WithCapacity(uint32(m), 10)
	v1Target := v1.Must(v1Cf)
	v1Other := v1.Must(v1Cf)
	runNoFalseNegativeAfterUnionTest(t, n, v1Target, v1Other, func() error {
		return v1Target.Union(v1Other)
	})

	v2Cf := v2.WithCapacity(uint32(m), 10)
	v2Target := v2.Must(v2Cf)
	v2Other := v2.Must(v2Cf)
	runNoFalseNegativeAfterUnionTest(t, n, v2Target, v2Other, func() error {
		return v2Target.Union(v2Other)
	})
}

func TestBloomFilter_Union_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	n := 1_000_000

	v1Cf := v1.WithAccuracy(0.001, uint32(n))
	v1Target := v1.Must(v1Cf)
	v1Other := v1.Must(v1Cf)
	runNoFalseNegativeAfterUnionTest(t, n, v1Target, v1Other, func() error {
		return v1Target.Union(v1Other)
	})

	v2Cf := v2.WithAccuracy(0.001, uint32(n))
	v2Target := v2.Must(v2Cf)
	v2Other := v2.Must(v2Cf)
	runNoFalseNegativeAfterUnionTest(t, n, v2Target, v2Other, func() error {
		return v2Target.Union(v2Other)
	})
}

func runNoFalseNegativeAfterUnionTest(t *testing.T, n int, target, other BloomFilterSpec, fn func() error) {
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

	err := fn()
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
