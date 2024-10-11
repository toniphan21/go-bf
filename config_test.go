package bf

import "testing"

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

func assertConfigEqual(t *testing.T, result Config, expected config) {
	c, ok := result.(*config)
	if !ok {
		t.Errorf("%v is not instance of config", result)
	}
	if *c != expected {
		t.Errorf("expected %v, got %v", expected, c)
	}
}
