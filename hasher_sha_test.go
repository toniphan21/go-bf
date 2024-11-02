package bf

import (
	"encoding/binary"
	"fmt"
	"reflect"
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
		count    int
		keySize  int
		keyCount byte
		expected [][]string
	}{
		{
			name:     "count 1 - hash 1 time",
			keyCount: 5,
			keySize:  16,
			count:    1,
			expected: [][]string{
				{"2cf20000", "4dba0000", "5fb00000", "a30e0000", "26e80000"},
			},
		},
		{
			name:     "count 2 - hash 1 time",
			keyCount: 5,
			keySize:  16,
			count:    2,
			expected: [][]string{
				{"2cf20000", "4dba0000", "5fb00000", "a30e0000", "26e80000"},
				{"3b2a0000", "c5b90000", "e29e0000", "1b160000", "1e5c0000"},
			},
		},
		{
			name:     "count 3 - hash 1 time",
			keyCount: 5,
			keySize:  16,
			count:    3,
			expected: [][]string{
				{"2cf20000", "4dba0000", "5fb00000", "a30e0000", "26e80000"},
				{"3b2a0000", "c5b90000", "e29e0000", "1b160000", "1e5c0000"},
				{"1fa70000", "425e0000", "73040000", "33620000", "938b0000"},
			},
		},
		{
			name:     "count 4 - hash 2 times",
			keyCount: 5,
			keySize:  16,
			count:    4,
			expected: [][]string{
				{"2cf20000", "4dba0000", "5fb00000", "a30e0000", "26e80000"},
				{"3b2a0000", "c5b90000", "e29e0000", "1b160000", "1e5c0000"},
				{"1fa70000", "425e0000", "73040000", "33620000", "938b0000"},
				{"98240000", "8a2a0000", "5c9b0000", "76880000", "27de0000"},
			},
		},
		{
			name:     "count 1 - hash 2 times",
			keyCount: 9,
			keySize:  32,
			count:    1,
			expected: [][]string{
				{"2cf24dba", "5fb0a30e", "26e83b2a", "c5b9e29e", "1b161e5c", "1fa7425e", "73043362", "938b9824", "8a2a5c9b"},
			},
		},
		{
			name:     "count 2 - hash 3 times",
			keyCount: 9,
			keySize:  32,
			count:    3,
			expected: [][]string{
				{"2cf24dba", "5fb0a30e", "26e83b2a", "c5b9e29e", "1b161e5c", "1fa7425e", "73043362", "938b9824", "8a2a5c9b"},
				{"768827de", "5a9552c3", "8a044c66", "959c68f6", "d2f21b52", "60af54d2", "f87db827", "cceeb7a9", "85ecc3da"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := shaHasher{hasher{hashSizeInBytes: shaSize, keySize: tc.keySize, keyCount: tc.keyCount}}
			r := s.Hash([]byte(hashInput), tc.count)

			for i := 0; i < len(tc.expected); i++ {
				for j := 0; j < len(tc.expected[i]); j++ {
					le := make([]byte, 4)
					binary.LittleEndian.PutUint32(le, uint32(r[i][j]))

					str := fmt.Sprintf("%x", le)
					assertStringEqual(t, str, tc.expected[i][j])
				}
			}
		})
	}
}

func TestShaHasher_doHash(t *testing.T) {
	h := &shaHasher{
		hasher{
			hashSizeInBytes: shaSize,
			keySize:         shaSize * 8,
		},
	}
	runTestHasherDoHash(t, h.hasher, h.doHash, shaHello, shaHelloPad0, shaHelloPad1, shaHelloPad2, shaHelloPad3)
}

func TestShaHasherFactory_Make(t *testing.T) {
	f := shaHasherFactory{}
	r := f.Make(5, 10)

	result, ok := r.(*shaHasher)
	if !ok {
		t.Errorf("expected *shaHasher, got %v", reflect.TypeOf(r))
	}
	if result.hashSizeInBytes != shaSize {
		t.Errorf("expected hashSizeInBytes to be %d, got %d", shaSize, result.hashSizeInBytes)
	}
	if result.keySize != 10 {
		t.Errorf("expected size 10, got %v", result.keySize)
	}
	if result.keyCount != 5 {
		t.Errorf("expected count 5, got %v", result.keyCount)
	}
}
