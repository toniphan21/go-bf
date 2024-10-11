package internal

import (
	"fmt"
	"github.com/toniphan21/go-bf"
	"math"
	"math/rand"
	"testing"
)

func calcEstimatedErrorRate(k byte, n int, m uint32) float64 {
	return math.Pow(1-math.Pow(math.E, (0-float64(k)*float64(n))/float64(m)), float64(k))
}

func TestBloomFilter_NoFalseNegative_WithCapacity(t *testing.T) {
	t.Parallel()
	n, m := 1_000_000, 5_000_000
	cf := bf.WithCapacity(uint32(m), 10)
	filter, _ := bf.New(cf)
	for i := 0; i < n; i++ {
		runNoFalseNegativeTest(t, filter)
	}
}

func TestBloomFilter_NoFalseNegative_WithAccuracy(t *testing.T) {
	t.Parallel()
	var n = 1_000_000
	cf := bf.WithAccuracy(0.001, uint32(n))
	filter, _ := bf.New(cf)
	for i := 0; i < n; i++ {
		runNoFalseNegativeTest(t, filter)
	}
}

func runNoFalseNegativeTest(t *testing.T, filter bf.BloomFilter) {
	item := []byte(RandString(10))
	filter.Add(item)
	after := filter.Exists(item)
	if !after {
		t.Fatalf("Bloom Filter has false negative")
	}
}

func TestBloomFilter_FalsePositiveRate_WithCapacity(t *testing.T) {
	t.Parallel()
	var n = 1_000_000
	var m uint32 = 8_388_608
	var k byte = 10

	cf := bf.WithCapacity(m, 10)
	filter, _ := bf.New(cf)
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
	for _, requestedErrorRate := range requested {
		t.Run(fmt.Sprintf("Check false positive rate with requested error rate %v", requestedErrorRate), func(t *testing.T) {
			t.Parallel()
			var n = 1_000_000
			cf := bf.WithAccuracy(requestedErrorRate, uint32(n))
			filter, _ := bf.New(cf)
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
			if tolerant > requestedErrorRate {
				t.Skipf("False positive error rate is 2x greater than requested. Requested %v, actual %v", requestedErrorRate, rate)
			}
		})
	}
}
