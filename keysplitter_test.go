package bf

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestKeySplitter_Split(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		count    int
		keyCount int
		keySize  int
		expected [][]string
	}{
		{
			name: "8 bits, overflow, count 1", source: "1a2b", count: 1, keySize: 8, keyCount: 3,
			expected: [][]string{
				{"1a000000", "2b000000", "00000000"},
			},
		},

		{
			name: "4 bits, 5 keys, count 1", source: shaHello, count: 1, keySize: 4, keyCount: 5,
			expected: [][]string{
				{"0c000000", "02000000", "02000000", "0f000000", "0d000000"},
			},
		},

		{
			name: "3 bits, 2 keys, count 1", source: shaHello, count: 1, keySize: 3, keyCount: 2,
			expected: [][]string{
				/*
					source in hex:           2c.f2.4d.ba...
					source in bin:           0010-1100.1111-0010.0100-1101.1011-1010
					because bitset stored byte and access via byte index so each byte need to be revered (indexed)
					source in bin - indexed: 0011-0100.0100-1111.1011-0010.0101-1101
					size: 3 count: 2
					key-1: 001 -reverse-> 100 -normalize-> 0100 -> 4
					key-2: 101 -reverse-> 101 -normalize-> 0101 -> 5
				*/
				{"04000000", "05000000"},
			},
		},

		{
			name: "11 bits, 2 keys, count 1", source: shaHello, count: 1, keySize: 11, keyCount: 2,
			expected: [][]string{
				/*
					source in hex:           2c.f2.4d.ba...
					source in bin:           0010-1100.1111-0010.0100-1101.1011-1010
					because bitset stored byte and access via byte index so each byte need to be revered (indexed)
					source in bin - indexed: 0011-0100.0100-1111.1011-0010.0101-1101
					size: 11 count: 2
					key-1: 0011-0100.010  -reverse-> 01000101100 -normalize-> 0010-0010-1100 -little-endian-> 2c02
					key-2: 0-1111.1011-00 -reverse-> 00110111110 -normalize-> 0001-1011-1110 -little-endian-> be01
				*/
				{"2c020000", "be010000"},
			},
		},

		{
			name: "8 bits, overflow, count 2", source: "1a2b", count: 2, keySize: 8, keyCount: 3,
			expected: [][]string{
				{"1a000000", "2b000000", "00000000"},
				{"00000000", "00000000", "00000000"},
			},
		},

		{
			name: "8 bits, overflow on 2, count 2", source: "1a2b3c4d", count: 2, keySize: 8, keyCount: 3,
			expected: [][]string{
				{"1a000000", "2b000000", "3c000000"},
				{"4d000000", "00000000", "00000000"},
			},
		},

		{
			name: "4 bits, 5 keys, count 2", source: shaHello, count: 2, keySize: 4, keyCount: 5,
			expected: [][]string{
				{"0c000000", "02000000", "02000000", "0f000000", "0d000000"},
				{"04000000", "0a000000", "0b000000", "0f000000", "05000000"},
			},
		},

		{
			name: "3 bits, 2 keys, count 2", source: shaHello, count: 2, keySize: 3, keyCount: 2,
			expected: [][]string{
				/*
					source in hex:           2c.f2.4d.ba...
					source in bin:           0010-1100.1111-0010.0100-1101.1011-1010
					because bitset stored byte and access via byte index so each byte need to be revered (indexed)
					source in bin - indexed: 0011-0100.0100-1111.1011-0010.0101-1101
					set: 1 size: 3 count: 2
					key-1: 001 -reverse-> 100 -normalize-> 0100 -> 4
					key-2: 101 -reverse-> 101 -normalize-> 0101 -> 5
					set: 2 size: 3 count: 2
					key-3: 000 -reverse-> 000 -normalize-> 0000 -> 0
					key-4: 100 -reverse-> 001 -normalize-> 0001 -> 1
				*/
				{"04000000", "05000000"},
				{"00000000", "01000000"},
			},
		},

		{
			name: "11 bits, 2 keys, count 2", source: shaHello, count: 2, keySize: 11, keyCount: 2,
			expected: [][]string{
				/*
					source in hex:           2c.f2.4d.ba.5f.b0...
					source in bin:           0010-1100.1111-0010.0100-1101.1011-1010.0101-1111.1011-0000
					because bitset stored byte and access via byte index so each byte need to be revered (indexed)
					source in bin - indexed: 0011-0100.0100-1111.1011-0010.0101-1101.1111-1010.0000-1101
					source in bin - indexed: 0011-0100.010|0-1111.1011-00|10.0101-1101.1|111-1010.0000|-1101
					set:1 size: 11 count: 2
					key-1: 0011-0100.010  -reverse-> 01000101100 -normalize-> 0010-0010-1100 -little-endian-> 2c02
					key-2: 0-1111.1011-00 -reverse-> 00110111110 -normalize-> 0001-1011-1110 -little-endian-> be01
					set:2 size: 11 count: 2
					key-3: 10.0101-1101.1 -reverse-> 11011101001 -normalize-> 0110-1110-1001 -little-endian-> e906
					key-4: 111-1010.0000  -reverse-> 00000101111 -normalize-> 0000-0010-1111 -little-endian-> 2f00

				*/
				{"2c020000", "be010000"},
				{"e9060000", "2f000000"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tc.source)
			ks := &KeySplitter{
				Source:   data,
				Count:    tc.count,
				KeySize:  tc.keySize,
				KeyCount: tc.keyCount,
			}
			r := ks.Split()

			for i := 0; i < len(tc.expected); i++ {
				for j := 0; j < len(tc.expected[i]); j++ {
					le := make([]byte, 4)
					binary.LittleEndian.PutUint32(le, uint32(r[i][j]))

					s := fmt.Sprintf("%x", le)
					assertStringEqual(t, s, tc.expected[i][j])
				}
			}
		})
	}
}
