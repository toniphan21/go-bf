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
	K() byte
	N() uint32
	M() float64
	E() float64
	Info() string
	StorageCapacity() uint32
	KeySize() byte
}

type config struct {
	mode            string
	k               byte
	n               uint32
	m               float64
	e               float64
	requestedE      float64
	storageCapacity uint32
	keySize         byte
}

func (c *config) K() byte {
	return c.k
}

func (c *config) N() uint32 {
	return c.n
}

func (c *config) M() float64 {
	return c.m
}

func (c *config) E() float64 {
	return c.e
}

func (c *config) StorageCapacity() uint32 {
	return c.storageCapacity
}

func (c *config) KeySize() byte {
	return c.keySize
}

func (c *config) Info() string {
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
		info = append(info, fmt.Sprintf("  - Size in bits of each has function: %v", c.keySize))

		info = append(info, fmt.Sprintf("  - Storage capacity: %v bits = %v bytes = %#.2fKB = %#.2fMB", c.storageCapacity, cB, cKB, cMB))
		info = append(info, fmt.Sprintf("  - Estimated error rate: %#.5f%%", c.e*100))
	default:
		info = append(info, "Config WithCapacity()")
		info = append(info, fmt.Sprintf("  - Storage capacity: %v bits = %v bytes = %#.2fKB = %#.2fMB", c.storageCapacity, cB, cKB, cMB))
		info = append(info, fmt.Sprintf("  - Number of hash functions: %v", c.k))
		info = append(info, fmt.Sprintf("  - Size in bits of each has function: %v", c.keySize))
		//points := []int{100, 200, 500, 1000, 2000, 5000, 10_000, 20_000, 50_000, 100_000, 200_000, 500_000}
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
	capacityInBits := noi * bitPerItem
	capacity := uint32(math.Ceil(noi * bitPerItem))

	nK := byte(math.Ceil(k))
	estimatedErrorRate := math.Pow(1-math.Pow(math.E, (0-float64(nK)*noi)/math.Ceil(capacityInBits)), float64(nK))
	return &config{
		mode:            "accuracy",
		k:               nK,
		m:               bitPerItem,
		n:               numberOfItems,
		e:               estimatedErrorRate,
		requestedE:      errorRate,
		storageCapacity: capacity,
		keySize:         byte(math.Ceil(math.Log2(float64(capacity)))),
	}
}

func WithCapacity(capacityInBits uint32, numberOfHashFunctions byte) Config {
	if capacityInBits == 0 {
		capacityInBits = DefaultSizeInBits
	}
	if numberOfHashFunctions == 0 {
		numberOfHashFunctions = DefaultNumberOfHasFunction
	}

	return &config{
		mode:            "capacity",
		k:               numberOfHashFunctions,
		m:               0,
		n:               0,
		e:               0,
		storageCapacity: capacityInBits,
		keySize:         byte(math.Ceil(math.Log2(float64(capacityInBits)))),
	}
}
