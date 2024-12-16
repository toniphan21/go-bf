package bf

import (
	"testing"
)

const fnvHello = "f14b58486483d94f708038798c29697f"
const fnvHelloPad0 = "09a825debb3c64bf6dc6a3066cccba81"
const fnvHelloPad1 = "0176dd4ddc3c64bf6dc6a13606a6d1ee"
const fnvHelloPad2 = "f24359bb713c64bf6dc69d3d1973fccf"
const fnvHelloPad3 = "eb01a908923c64bf6dc69bb684a13444"

func TestFnvHasher_Hash(t *testing.T) {
	cases := []struct {
		name     string
		configs  []ConfigBlock
		expected [][]Key
	}{
		{
			name: "1 block - hash 1 time",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
			},
			expected: [][]Key{
				{0x4bf1, 0x4858, 0x8364, 0x4fd9, 0x8070},
			},
		},

		{
			name: "2 blocks - same size - hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
			},
			expected: [][]Key{
				{0x4bf1, 0x4858, 0x8364, 0x4fd9, 0x8070},
				{0x7938, 0x298c, 0x7f69, 0xa809, 0xde25},
			},
		},

		{
			name: "2 blocks - different size - hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
				{NumberOfHashFunctions: 4, KeySizeInBits: 20},
			},
			expected: [][]Key{
				{0x4bf1, 0x4858, 0x8364, 0x4fd9, 0x8070},
				{0xc7938, 0x69298, 0x8097f, 0xde25a},
			},
		},

		{
			name: "3 blocks - same size - hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
			},
			expected: [][]Key{
				{0x4bf1, 0x4858, 0x8364, 0x4fd9, 0x8070},
				{0x7938, 0x298c, 0x7f69, 0xa809, 0xde25},
				{0x3cbb, 0xbf64, 0xc66d, 0x06a3, 0xcc6c},
			},
		},

		{
			name: "3 blocks - different size - hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
				{NumberOfHashFunctions: 4, KeySizeInBits: 20},
				{NumberOfHashFunctions: 6, KeySizeInBits: 12},
			},
			expected: [][]Key{
				{0x4bf1, 0x4858, 0x8364, 0x4fd9, 0x8070},
				{0xc7938, 0x69298, 0x8097f, 0xde25a},
				{0xcbb, 0x643, 0xdbf, 0xc66, 0x6a3, 0x6c0},
			},
		},

		{
			name: "1 block - hash 3 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 9, KeySizeInBits: 32},
			},
			expected: [][]Key{
				{
					0x48584bf1, 0x4fd98364, 0x79388070, 0x7f69298c,
					0xde25a809, 0xbf643cbb, 0x06a3c66d, 0x81bacc6c,
					0x4ddd7601,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := fnvHasher{hasher{}}
			r := s.Hash([]byte(hasherInput), tc.configs)
			if len(r) != len(tc.expected) {
				t.Errorf("expected %d keyset, got %d", len(tc.expected), len(r))
			}

			for i := 0; i < len(tc.expected); i++ {
				if len(r) != len(tc.expected) {
					t.Errorf("keyset %v: expected %d key, got %d", i, len(tc.expected), len(r))
				}
				for j := 0; j < len(tc.expected[i]); j++ {
					if r[i][j] != tc.expected[i][j] {
						t.Errorf("index (%v, %v): got %v, want %v", i, j, r[i][j], tc.expected[i][j])
					}
				}
			}
		})
	}
}

func TestFnvHasher_doHash(t *testing.T) {
	h := &fnvHasher{hasher{}}
	runTestHasherDoHash(t, h.hasher, fnvSize, h.doHash, fnvHello, fnvHelloPad0, fnvHelloPad1, fnvHelloPad2, fnvHelloPad3)
}
