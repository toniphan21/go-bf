package bf

import (
	"fmt"
	"math"
	"strings"
)

const DefaultErrorRate = 0.0001
const DefaultNumberOfItem = 1000000
const DefaultSizeInBits = 8192
const DefaultNumberOfHasFunction = 5

type Config interface {
	Info() string
	NumberOfHashFunctions() byte
	StorageCapacity() uint32
	KeySize() byte
}

func calcKeyMinSizeFromCapacity(capacity uint32) byte {
	return byte(math.Ceil(math.Log2(float64(capacity))))
}

func calcEstimatedErrorRate(k byte, n int, m uint32) float64 {
	return math.Pow(1-math.Pow(math.E, (0-float64(k)*float64(n))/float64(m)), float64(k))
}

type config struct {
	mode            string
	k               byte
	n               uint32
	m               float64
	e               float64
	requestedE      float64
	storageCapacity uint32
}

func (c config) NumberOfHashFunctions() byte {
	return c.k
}

func (c config) StorageCapacity() uint32 {
	return c.storageCapacity
}

func (c config) KeySize() byte {
	return calcKeyMinSizeFromCapacity(c.storageCapacity)
}

func (c config) Info() string {
	cB := c.storageCapacity / 8
	d := c.storageCapacity % 8
	if d > 0 {
		cB += 1
	}
	cKB := float64(cB) / 1024
	cMB := float64(cKB) / 1024

	var info []string
	switch c.mode {
	case "accuracy":
		info = append(info, "Config WithAccuracy()")
		info = append(info, fmt.Sprintf("  - Requested error rate: %#.5f%%", c.requestedE*100))
		info = append(info, fmt.Sprintf("  - Expected number of items: %d", c.n))
		info = append(info, fmt.Sprintf("  - Bits per item: %#.3f", c.m))
		info = append(info, fmt.Sprintf("  - Number of hash functions: %v", c.k))
		info = append(info, fmt.Sprintf("  - Size in bits of each has function: %v", c.KeySize()))

		info = append(info, fmt.Sprintf("  - Storage capacity: %v bits = %v bytes = %#.2fKB = %#.2fMB", c.storageCapacity, cB, cKB, cMB))
		info = append(info, fmt.Sprintf("  - Estimated error rate: %#.5f%%", c.e*100))
	default:
		info = append(info, "Config WithCapacity()")
		info = append(info, fmt.Sprintf("  - Storage capacity: %v bits = %v bytes = %#.2fKB = %#.2fMB", c.storageCapacity, cB, cKB, cMB))
		info = append(info, fmt.Sprintf("  - Number of hash functions: %v", c.k))
		info = append(info, fmt.Sprintf("  - Size in bits of each has function: %v", c.KeySize()))

		log := int(math.Ceil(math.Log10(float64(c.storageCapacity))))
		info = append(info, fmt.Sprintf("  - Estimated error rate by n - number of added items:"))
		fmtString := fmt.Sprintf("      n=%%%vd; estimated error rate: %%#.5f%%%%", log+1)
		var i = log - 3
		if i < 1 {
			i = 1
		}
		base := int(math.Pow(10, float64(i)))
		for ; i <= log; i++ {
			n := base
			info = append(info, fmt.Sprintf(fmtString, n, 100*calcEstimatedErrorRate(c.k, n, c.storageCapacity)))

			n = 2 * base
			info = append(info, fmt.Sprintf(fmtString, n, 100*calcEstimatedErrorRate(c.k, n, c.storageCapacity)))

			n = 5 * base
			info = append(info, fmt.Sprintf(fmtString, n, 100*calcEstimatedErrorRate(c.k, n, c.storageCapacity)))

			base *= 10
		}
	}

	return strings.Join(info, "\n")
}

func WithAccuracy(errorRate float64, numberOfItems uint32) Config {
	if numberOfItems == 0 {
		numberOfItems = DefaultNumberOfItem
	}
	if errorRate <= 0 {
		errorRate = DefaultErrorRate
	}

	noi := float64(numberOfItems)
	log2 := math.Abs(math.Log2(errorRate))
	k := log2
	bitPerItem := 1.44 * log2
	capacity := uint32(math.Ceil(noi * bitPerItem))

	nK := byte(math.Ceil(k))
	return config{
		mode:            "accuracy",
		k:               nK,
		m:               bitPerItem,
		n:               numberOfItems,
		e:               calcEstimatedErrorRate(nK, int(numberOfItems), capacity),
		requestedE:      errorRate,
		storageCapacity: capacity,
	}
}

func WithCapacity(capacityInBits uint32, numberOfHashFunctions byte) Config {
	if capacityInBits == 0 {
		capacityInBits = DefaultSizeInBits
	}
	if numberOfHashFunctions == 0 {
		numberOfHashFunctions = DefaultNumberOfHasFunction
	}

	return config{
		mode:            "capacity",
		k:               numberOfHashFunctions,
		storageCapacity: capacityInBits,
	}
}
