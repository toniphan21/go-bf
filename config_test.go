package bf

import (
	"testing"
)

func TestWithCapacity(t *testing.T) {
	cases := []struct {
		name     string
		capacity uint32
		k        byte
		expected config
	}{
		{
			name: "invalid capacity", capacity: 0, k: 5,
			expected: config{
				mode:            "capacity",
				storageCapacity: DefaultSizeInBits,
				k:               5,
			},
		},

		{
			name: "invalid number of hash function", capacity: 1000, k: 0,
			expected: config{
				mode:            "capacity",
				storageCapacity: 1000,
				k:               DefaultNumberOfHasFunction,
			},
		},

		{
			name: "custom values", capacity: 10000, k: 10,
			expected: config{
				mode:            "capacity",
				storageCapacity: 10000,
				k:               10,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := WithCapacity(tc.capacity, tc.k)
			assertConfigEqual(t, c, tc.expected)
		})
	}
}

func TestWithAccuracy(t *testing.T) {
	cases := []struct {
		name             string
		e                float64
		n                uint32
		expectedE        float64
		expectedN        uint32
		expectedK        byte
		expectedCapacity uint32
	}{
		{
			name:             "invalid capacity",
			e:                0,
			n:                5,
			expectedK:        14,
			expectedCapacity: 96,
			expectedE:        DefaultErrorRate,
			expectedN:        5,
		},

		{
			name: "invalid number of item", e: 0.001, n: 0,
			expectedK:        10,
			expectedCapacity: 14350730,
			expectedE:        0.001,
			expectedN:        DefaultNumberOfItem,
		},

		{
			name: "custom values", e: 0.01, n: 10000,
			expectedK:        7,
			expectedCapacity: 95672,
			expectedE:        0.01,
			expectedN:        10000,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := WithAccuracy(tc.e, tc.n)

			if c.StorageCapacity() != tc.expectedCapacity {
				t.Errorf("got %v, want %v", c.StorageCapacity(), tc.expectedCapacity)
			}

			if c.NumberOfHashFunctions() != tc.expectedK {
				t.Errorf("got %v, want %v", c.NumberOfHashFunctions(), tc.expectedK)
			}

			cf, ok := c.(config)
			if !ok {
				t.Errorf("%v is not instance of config", c)
			}
			if cf.requestedE != tc.expectedE {
				t.Errorf("got %v, want %v", cf.e, tc.expectedE)
			}
			if cf.n != tc.expectedN {
				t.Errorf("got %v, want %v", cf.n, tc.expectedN)
			}
		})
	}
}

func assertConfigEqual(t *testing.T, result Config, expected config) {
	c, ok := result.(config)
	if !ok {
		t.Errorf("%v is not instance of config", result)
	}
	if c != expected {
		t.Errorf("expected %v, got %v", expected, c)
	}
}

func TestConfig_Info_WithAccuracy(t *testing.T) {
	var errorRate = 0.001
	var numberOfItems uint32 = 10_000_000
	cf := WithAccuracy(errorRate, numberOfItems)
	expected := `Config WithAccuracy()
  - Requested error rate: 0.10000%
  - Expected number of items: 10000000
  - Bits per item: 14.351
  - Number of hash functions: 10
  - Size in bits of each has function: 28
  - Storage capacity: 143507294 bits = 17938412 bytes = 17517.98KB = 17.11MB
  - Estimated error rate: 0.10130%`

	if cf.Info() != expected {
		t.Errorf("expected %v, got %v", expected, cf.Info())
	}
}

func TestConfig_Info_WithCapacity(t *testing.T) {
	var capacityInBits uint32 = 65_536
	var numberOfHashFunctions byte = 5
	cf := WithCapacity(capacityInBits, numberOfHashFunctions)
	expected := `Config WithCapacity()
  - Storage capacity: 65536 bits = 8192 bytes = 8.00KB = 0.01MB
  - Number of hash functions: 5
  - Size in bits of each has function: 16
  - Estimated error rate by n - number of added items:
      n=   100; estimated error rate: 0.00000%
      n=   200; estimated error rate: 0.00000%
      n=   500; estimated error rate: 0.00001%
      n=  1000; estimated error rate: 0.00021%
      n=  2000; estimated error rate: 0.00568%
      n=  5000; estimated error rate: 0.32083%
      n= 10000; estimated error rate: 4.33023%
      n= 20000; estimated error rate: 29.35056%
      n= 50000; estimated error rate: 89.45317%
      n=100000; estimated error rate: 99.75726%
      n=200000; estimated error rate: 99.99988%
      n=500000; estimated error rate: 100.00000%`

	if cf.Info() != expected {
		t.Errorf("expected %v, got %v", expected, cf.Info())
	}
}

func TestConfig_Info_WithCapacitySmall(t *testing.T) {
	var capacityInBits uint32 = 512
	var numberOfHashFunctions byte = 5
	cf := WithCapacity(capacityInBits, numberOfHashFunctions)
	expected := `Config WithCapacity()
  - Storage capacity: 512 bits = 64 bytes = 0.06KB = 0.00MB
  - Number of hash functions: 5
  - Size in bits of each has function: 9
  - Estimated error rate by n - number of added items:
      n=  10; estimated error rate: 0.00070%
      n=  20; estimated error rate: 0.01758%
      n=  50; estimated error rate: 0.86047%
      n= 100; estimated error rate: 9.41504%
      n= 200; estimated error rate: 46.54427%
      n= 500; estimated error rate: 96.26912%
      n=1000; estimated error rate: 99.97131%
      n=2000; estimated error rate: 100.00000%
      n=5000; estimated error rate: 100.00000%`
	if cf.Info() != expected {
		t.Errorf("expected %v, got %v", expected, cf.Info())
	}
}
