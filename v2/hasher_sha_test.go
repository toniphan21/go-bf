package bf

import (
	"testing"
)

const shaHello = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
const shaHelloPad0 = "8a2a5c9b768827de5a9552c38a044c66959c68f6d2f21b5260af54d2f87db827"
const shaHelloPad1 = "cceeb7a985ecc3dabcb4c8f666cd637f16f008e3c963db6aa6f83a7b288c54ef"
const shaHelloPad2 = "29f3ced0b171e52626c66bedaf76469f1efda5c110b47ea24228ef25e61859cc"
const shaHelloPad3 = "0b4d354d56ea9a985571a56b1829f33d072e7902c1afaf981381089b9eb00ffe"

func TestShaHasher_Hash(t *testing.T) {
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
				{0xf22c, 0xba4d, 0xb05f, 0x0ea3, 0xe826},
			},
		},

		{
			name: "2 blocks - same size - hash 1 time",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
			},
			expected: [][]Key{
				{0xf22c, 0xba4d, 0xb05f, 0x0ea3, 0xe826},
				{0x2a3b, 0xb9c5, 0x9ee2, 0x161b, 0x5c1e},
			},
		},

		{
			name: "2 blocks - different size - hash 1 time",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: 16},
				{NumberOfHashFunctions: 4, KeySizeInBits: 20},
			},
			expected: [][]Key{
				{0xf22c, 0xba4d, 0xb05f, 0x0ea3, 0xe826},
				{0x52a3b, 0xe2b9c, 0x61b9e, 0x5c1e1},
			},
		},

		{
			name: "3 blocks - same size - hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 7, KeySizeInBits: 16},
				{NumberOfHashFunctions: 7, KeySizeInBits: 16},
				{NumberOfHashFunctions: 7, KeySizeInBits: 16},
			},
			expected: [][]Key{
				{0xf22c, 0xba4d, 0xb05f, 0x0ea3, 0xe826, 0x2a3b, 0xb9c5},
				{0x9ee2, 0x161b, 0x5c1e, 0xa71f, 0x5e42, 0x0473, 0x6233},
				{0x8b93, 0x2498, 0x2a8a, 0x9b5c, 0x8876, 0xde27, 0x955a},
			},
		},

		{
			name: "3 blocks - different size - hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 7, KeySizeInBits: 16},
				{NumberOfHashFunctions: 4, KeySizeInBits: 20},
				{NumberOfHashFunctions: 6, KeySizeInBits: 24},
			},
			expected: [][]Key{
				{0xf22c, 0xba4d, 0xb05f, 0x0ea3, 0xe826, 0x2a3b, 0xb9c5},
				{0xb9ee2, 0x1e161, 0x71f5c, 0x5e42a},
				{0x330473, 0x8b9362, 0x8a2498, 0x9b5c2a, 0x278876, 0x955ade},
			},
		},

		{
			name: "1 block - hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 9, KeySizeInBits: 32},
			},
			expected: [][]Key{
				{
					0xba4df22c, 0x0ea3b05f, 0x2a3be826, 0x9ee2b9c5,
					0x5c1e161b, 0x5e42a71f, 0x62330473, 0x24988b93,
					0x9b5c2a8a,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := shaHasher{hasher{}}
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

func TestShaHasher_doHash(t *testing.T) {
	h := &shaHasher{hasher{}}
	runTestHasherDoHash(t, h.hasher, shaSize, h.doHash, shaHello, shaHelloPad0, shaHelloPad1, shaHelloPad2, shaHelloPad3)
}
