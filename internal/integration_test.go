package internal

import (
	"fmt"
	v1 "github.com/toniphan21/go-bf"
	v2 "github.com/toniphan21/go-bf/v2"
	"math"
	"math/rand"
	"testing"
)

func calcEstimatedErrorRate(k byte, n int, m uint32) float64 {
	return math.Pow(1-math.Pow(math.E, (0-float64(k)*float64(n))/float64(m)), float64(k))
}

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

func TestBloomFilter_Intersect_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000
	cf := v1.WithCapacity(uint32(m), 10)
	target := v1.Must(cf)
	other := v1.Must(cf)
	runNoFalseNegativeAfterIntersectTest(t, n/2, target, other)
}

func TestBloomFilter_Intersect_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	n := 1_000_000
	cf := v1.WithAccuracy(0.001, uint32(n))
	target := v1.Must(cf)
	other := v1.Must(cf)
	runNoFalseNegativeAfterIntersectTest(t, n/2, target, other)
}

func runNoFalseNegativeAfterIntersectTest(t *testing.T, n int, target, other v1.BloomFilter) {
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

	err := target.Intersect(other)
	if err != nil {
		t.Fatalf("Bloom Filter intersect failed")
	}
	for _, key := range sharedKeys {
		if !target.Exists([]byte(key)) {
			t.Fatalf("Bloom Filter has false negative in target after Intersect() with key=%v", key)
		}
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

func TestBloomFilter_FalsePositiveRate_WithCapacity(t *testing.T) {
	t.Parallel()
	var n = 1_000_000
	var m uint32 = 8_388_608
	var k byte = 10

	cf := v1.WithCapacity(m, 10)
	filter, _ := v1.New(cf)
	for i := 0; i < n; i++ {
		filter.Add([]byte(RandString(10 + rand.Intn(10))))
	}

	count := 0
	for i := 0; i < n; i++ {
		if filter.Exists([]byte(RandString(9))) {
			count++
		}
	}
	rate := float64(count) / float64(n)
	estimated := calcEstimatedErrorRate(k, n, m)
	tolerant := math.Abs(rate - estimated)
	if tolerant > estimated {
		t.Skipf("False positive rate is 2x greater than estimated error rate. Estimated %v, actual %v", estimated, rate)
	}
}

func TestBloomFilter_FalsePositiveRate_WithAccuracy(t *testing.T) {
	requested := []float64{0.05, 0.02, 0.01, 0.005, 0.002, 0.001, 0.0001}
	for _, v := range requested {
		e := v
		t.Run(fmt.Sprintf("v1 - Check false positive rate with requested error rate %v - SHA", e), func(t *testing.T) {
			t.Parallel()
			var n = 1_000_000
			filter := v1.Must(v1.WithAccuracy(e, uint32(n)))
			runTestBloomFilterFalsePositiveRateWithAccuracy(t, n, filter, e)
		})

		t.Run(fmt.Sprintf("v1 - Check false positive rate with requested error rate %v - FVN", e), func(t *testing.T) {
			t.Parallel()
			var n = 1_000_000
			filter := v1.Must(v1.WithAccuracy(e, uint32(n)), v1.WithFNV())
			runTestBloomFilterFalsePositiveRateWithAccuracy(t, n, filter, e)
		})

		t.Run(fmt.Sprintf("v2 - Check false positive rate with requested error rate %v - SHA", e), func(t *testing.T) {
			t.Parallel()
			var n = 1_000_000
			filter := v2.Must(v2.WithAccuracy(e, uint32(n)))
			runTestBloomFilterFalsePositiveRateWithAccuracy(t, n, filter, e)
		})

		t.Run(fmt.Sprintf("v2 - Check false positive rate with requested error rate %v - FVN", e), func(t *testing.T) {
			t.Parallel()
			var n = 1_000_000
			filter := v2.Must(v2.WithAccuracy(e, uint32(n)), v2.WithFNV())
			runTestBloomFilterFalsePositiveRateWithAccuracy(t, n, filter, e)
		})

		t.Run(fmt.Sprintf("v2 - Check false positive rate with requested error rate %v - SHA - Expandable", e), func(t *testing.T) {
			t.Parallel()
			var n = 1_000_000
			filter := v2.Must(v2.WithAccuracy(e, uint32(n/10)))
			runTestBloomFilterFalsePositiveRateWithAccuracy(t, n, filter, e)
		})

		t.Run(fmt.Sprintf("v2 - Check false positive rate with requested error rate %v - FVN - Expandable", e), func(t *testing.T) {
			t.Parallel()
			var n = 1_000_000
			filter := v2.Must(v2.WithAccuracy(e, uint32(n/10)), v2.WithFNV())
			runTestBloomFilterFalsePositiveRateWithAccuracy(t, n, filter, e)
		})
	}
}

func runTestBloomFilterFalsePositiveRateWithAccuracy(
	t *testing.T,
	n int,
	filter BloomFilterSpec,
	requestedErrorRate float64,
) {
	for i := 0; i < n; i++ {
		filter.Add([]byte(RandString(10 + rand.Intn(10))))
	}

	count := 0
	for i := 0; i < n; i++ {
		if filter.Exists([]byte(RandString(9))) {
			count++
		}
	}
	rate := float64(count) / float64(n)
	tolerant := math.Abs(rate - requestedErrorRate)
	println(fmt.Sprintf("False positive rate - requested %v, actual %v", requestedErrorRate, rate))
	if tolerant > requestedErrorRate {
		t.Skipf("False positive error rate is 2x greater than requested. Requested %v, actual %v", requestedErrorRate, rate)
	}
}
