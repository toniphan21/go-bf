package internal

import (
	v1 "github.com/toniphan21/go-bf"
	v2 "github.com/toniphan21/go-bf/v2"
	"testing"
)

func TestBloomFilter_Intersect_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000

	v1Cf := v1.WithCapacity(uint32(m), 10)
	v1Target := v1.Must(v1Cf)
	v1Other := v1.Must(v1Cf)
	runNoFalseNegativeAfterIntersectTest(t, n/2, v1Target, v1Other, func() error {
		return v1Target.Intersect(v1Other)
	})

	v2Cf := v2.WithCapacity(uint32(m), 10)
	v2Target := v2.Must(v2Cf)
	v2Other := v2.Must(v2Cf)
	runNoFalseNegativeAfterIntersectTest(t, n/2, v2Target, v2Other, func() error {
		return v2Target.Intersect(v2Other)
	})
}

func TestBloomFilter_Intersect_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	n := 1_000_000

	v1Cf := v1.WithAccuracy(0.001, uint32(n))
	v1Target := v1.Must(v1Cf)
	v1Other := v1.Must(v1Cf)
	runNoFalseNegativeAfterIntersectTest(t, n/2, v1Target, v1Other, func() error {
		return v1Target.Intersect(v1Other)
	})

	v2Cf := v2.WithAccuracy(0.001, uint32(n))
	v2Target := v2.Must(v2Cf)
	v2Other := v2.Must(v2Cf)
	runNoFalseNegativeAfterIntersectTest(t, n/2, v2Target, v2Other, func() error {
		return v2Target.Intersect(v2Other)
	})
}

func runNoFalseNegativeAfterIntersectTest(t *testing.T, n int, target, other BloomFilterSpec, fn func() error) {
	targetKeys := make([]string, n)
	sharedKeys := make([]string, n)
	otherKeys := make([]string, n)

	for i := 0; i < n; i++ {
		targetKeys[i] = RandString(10)
		sharedKeys[i] = RandString(11)
		otherKeys[i] = RandString(12)

		target.Add([]byte(targetKeys[i]))
		target.Add([]byte(sharedKeys[i]))

		other.Add([]byte(otherKeys[i]))
		other.Add([]byte(sharedKeys[i]))

		if !target.Exists([]byte(targetKeys[i])) || !target.Exists([]byte(sharedKeys[i])) {
			t.Fatalf("Bloom Filter has false negative")
		}
		if !other.Exists([]byte(otherKeys[i])) || !other.Exists([]byte(sharedKeys[i])) {
			t.Fatalf("Bloom Filter has false negative")
		}
	}

	err := fn()
	if err != nil {
		t.Fatalf("Bloom Filter intersect failed")
	}
	for _, key := range sharedKeys {
		if !target.Exists([]byte(key)) {
			t.Fatalf("Bloom Filter has false negative in target after Intersect() with key=%v", key)
		}
	}
}
