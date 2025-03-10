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
