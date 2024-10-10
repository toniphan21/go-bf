package bf

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

const shaInput = "hello"
const shaRawHashHelloOneTime = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
const shaRawHashHelloTwoTime = "8a2a5c9b768827de5a9552c38a044c66959c68f6d2f21b5260af54d2f87db827cceeb7a985ecc3dabcb4c8f666cd637f16f008e3c963db6aa6f83a7b288c54ef"

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
				binary.LittleEndian.PutUint32(le, ki)

				s := fmt.Sprintf("%x", le)
				assertStringEqual(t, s, tc.expected[i])
			}
		})
	}
}

func TestShaHash_Hash(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		count    byte
		size     byte
		expected []string
	}{
		{
			name: "8 bits, 5 keys", input: shaInput, size: 8, count: 5,
			expected: []string{"2c000000", "f2000000", "4d000000", "ba000000", "5f000000"},
		},
		{
			name: "32 bits, 5 keys", input: shaInput, size: 32, count: 5,
			expected: []string{"2cf24dba", "5fb0a30e", "26e83b2a", "c5b9e29e", "1b161e5c"},
		},
		{
			name: "32 bits, 10 keys", input: shaInput, size: 32, count: 10,
			expected: []string{
				"8a2a5c9b", "768827de", "5a9552c3", "8a044c66", "959c68f6",
				"d2f21b52", "60af54d2", "f87db827", "cceeb7a9", "85ecc3da",
			},
		},
		{
			name: "24 bits, 20 keys", input: shaInput, size: 24, count: 20,
			expected: []string{
				"8a2a5c00", "9b768800", "27de5a00", "9552c300", "8a044c00",
				"66959c00", "68f6d200", "f21b5200", "60af5400", "d2f87d00",
				"b827cc00", "eeb7a900", "85ecc300", "dabcb400", "c8f66600",
				"cd637f00", "16f00800", "e3c96300", "db6aa600", "f83a7b00",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := shaHash{size: tc.size, count: tc.count}
			r := h.Hash([]byte(tc.input))
			for i, ki := range r {
				le := make([]byte, 4)
				binary.LittleEndian.PutUint32(le, ki)

				s := fmt.Sprintf("%x", le)
				assertStringEqual(t, s, tc.expected[i])
			}
		})
	}
}

func TestShaHash_doHash(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		n        byte
		expected string
	}{
		{
			name: "1 time", input: shaInput, n: 1, expected: shaRawHashHelloOneTime,
		},
		{
			name: "2 times", input: shaInput, n: 2, expected: shaRawHashHelloTwoTime,
		},
		{
			name:  "3 times",
			input: shaInput,
			n:     3,
			expected: "8a2a5c9b768827de5a9552c38a044c66959c68f6d2f21b5260af54d2f87db827" +
				"cceeb7a985ecc3dabcb4c8f666cd637f16f008e3c963db6aa6f83a7b288c54ef" +
				"29f3ced0b171e52626c66bedaf76469f1efda5c110b47ea24228ef25e61859cc",
		},
		{
			name:  "4 times",
			input: shaInput,
			n:     4,
			expected: "8a2a5c9b768827de5a9552c38a044c66959c68f6d2f21b5260af54d2f87db827" +
				"cceeb7a985ecc3dabcb4c8f666cd637f16f008e3c963db6aa6f83a7b288c54ef" +
				"29f3ced0b171e52626c66bedaf76469f1efda5c110b47ea24228ef25e61859cc" +
				"0b4d354d56ea9a985571a56b1829f33d072e7902c1afaf981381089b9eb00ffe",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := shaHash{}
			input := []byte(tc.input)
			r := h.doHash(tc.n, &input)

			result := fmt.Sprintf("%x", r)
			assertStringEqual(t, result, tc.expected)
		})
	}
}
