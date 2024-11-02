package bf

import (
	"encoding/binary"
	"fmt"
	"reflect"
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
				{"f14b0000", "58480000", "64830000", "d94f0000", "70800000"},
			},
		},
		{
			name:     "count 2 - hash 2 times",
			keyCount: 5,
			keySize:  16,
			count:    2,
			expected: [][]string{
				{"f14b0000", "58480000", "64830000", "d94f0000", "70800000"},
				{"38790000", "8c290000", "697f0000", "09a80000", "25de0000"},
			},
		},
		{
			name:     "count 3 - hash 2 times",
			keyCount: 5,
			keySize:  16,
			count:    3,
			expected: [][]string{
				{"f14b0000", "58480000", "64830000", "d94f0000", "70800000"},
				{"38790000", "8c290000", "697f0000", "09a80000", "25de0000"},
				{"bb3c0000", "64bf0000", "6dc60000", "a3060000", "6ccc0000"},
			},
		},
		{
			name:     "count 4 - hash 3 times",
			keyCount: 5,
			keySize:  16,
			count:    4,
			expected: [][]string{
				{"f14b0000", "58480000", "64830000", "d94f0000", "70800000"},
				{"38790000", "8c290000", "697f0000", "09a80000", "25de0000"},
				{"bb3c0000", "64bf0000", "6dc60000", "a3060000", "6ccc0000"},
				{"ba810000", "01760000", "dd4d0000", "dc3c0000", "64bf0000"},
			},
		},
		{
			name:     "count 1 - hash 3 times",
			keyCount: 9,
			keySize:  32,
			count:    1,
			expected: [][]string{
				{"f14b5848", "6483d94f", "70803879", "8c29697f", "09a825de", "bb3c64bf", "6dc6a306", "6cccba81", "0176dd4d"},
			},
		},
		{
			name:     "count 2 - hash 3 times",
			keyCount: 9,
			keySize:  32,
			count:    3,
			expected: [][]string{
				{"f14b5848", "6483d94f", "70803879", "8c29697f", "09a825de", "bb3c64bf", "6dc6a306", "6cccba81", "0176dd4d"},
				{"dc3c64bf", "6dc6a136", "06a6d1ee", "f24359bb", "713c64bf", "6dc69d3d", "1973fccf", "eb01a908", "923c64bf"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := fnvHasher{hasher{hashSizeInBytes: fnvSize, keySize: tc.keySize, keyCount: tc.keyCount}}
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

func TestFnvHasher_doHash(t *testing.T) {
	h := &fnvHasher{
		hasher{
			hashSizeInBytes: fnvSize,
			keySize:         fnvSize * 8,
		},
	}
	runTestHasherDoHash(t, h.hasher, h.doHash, fnvHello, fnvHelloPad0, fnvHelloPad1, fnvHelloPad2, fnvHelloPad3)
}

func TestFnvHasherFactory_Make(t *testing.T) {
	f := fnvHasherFactory{}
	r := f.Make(5, 10)

	result, ok := r.(*fnvHasher)
	if !ok {
		t.Errorf("expected *fnvHasher, got %v", reflect.TypeOf(r))
	}
	if result.hashSizeInBytes != fnvSize {
		t.Errorf("expected hashSizeInBytes to be %d, got %d", fnvSize, result.hashSizeInBytes)
	}
	if result.keySize != 10 {
		t.Errorf("expected size 10, got %v", result.keySize)
	}
	if result.keyCount != 5 {
		t.Errorf("expected count 5, got %v", result.keyCount)
	}
}
