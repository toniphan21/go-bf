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
		length   int
		count    int
		size     int
		expected []string
	}{
		{
			name: "8 bits, overflow", source: "0a0b", length: 16, size: 8, count: 3,
			expected: []string{"0a000000", "0b000000", "00000000"},
		},
		{
			name: "8 bits, 5 keys", source: shaRawHashHelloOneTime, length: 256, size: 4, count: 5,
			expected: []string{"0c000000", "02000000", "02000000", "0f000000", "0d000000"},
		},
		{
			name: "3 bits, 2 keys", source: shaRawHashHelloOneTime, length: 256, size: 3, count: 2,
			expected: []string{
				/*
					source in hex:           2c.f2.4d.ba...
					source in bin:           0010-1100.1111-0010.0100-1101.1011-1010
					because bitset stored byte and access via byte index so each byte need to be revered (indexed)
					source in bin - indexed: 0011-0100.0100-1111.1011-0010.0101-1101
					size: 3 count: 2
					key-1: 001 -reverse-> 100 -normalize-> 0100 -> 4
					key-2: 101 -reverse-> 101 -normalize-> 0101 -> 5
				*/
				"04000000", "05000000",
			},
		},
		{
			name: "11 bits, 2 keys", source: shaRawHashHelloOneTime, length: 256, size: 11, count: 2,
			expected: []string{
				/*
					source in hex:           2c.f2.4d.ba...
					source in bin:           0010-1100.1111-0010.0100-1101.1011-1010
					because bitset stored byte and access via byte index so each byte need to be revered (indexed)
					source in bin - indexed: 0011-0100.0100-1111.1011-0010.0101-1101
					size: 11 count: 2
					key-1: 0011-0100.010  -reverse-> 01000101100 -normalize-> 0010-0010-1100 -little-endian-> 2c02
					key-2: 0-1111.1011-00 -reverse-> 00110111110 -normalize-> 0001-1011-1110 -little-endian-> be01
				*/
				"2c020000", "be010000",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tc.source)
			ks := &KeySplitter{
				Source:   data,
				Length:   tc.length,
				KeySize:  tc.size,
				KeyCount: tc.count,
			}
			r := ks.Split()
			for i, ki := range r {
				le := make([]byte, 4)
				binary.LittleEndian.PutUint32(le, uint32(ki))

				s := fmt.Sprintf("%x", le)
				assertStringEqual(t, s, tc.expected[i])
			}
		})
	}
}
