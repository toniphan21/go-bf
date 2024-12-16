package bf

import (
	"bytes"
	"encoding/hex"
	"testing"
)

type mockHashFn struct {
	hashCalledCount int
	hashCalledWith  map[int][]byte
	hashReturn      map[int][]byte
}

func (m *mockHashFn) Hash(input *[]byte) []byte {
	m.hashCalledCount++
	if m.hashCalledWith == nil {
		m.hashCalledWith = make(map[int][]byte)
	}
	clone := make([]byte, len(*input))
	for i := 0; i < len(clone); i++ {
		clone[i] = (*input)[i]
	}
	m.hashCalledWith[m.hashCalledCount-1] = clone

	if m.hashReturn != nil {
		r, ok := m.hashReturn[m.hashCalledCount-1]
		if !ok {
			return []byte{}
		}
		return r
	}
	return []byte{}
}

func (m *mockHashFn) assertCalledCount(t *testing.T, expected int) {
	if m.hashCalledCount != expected {
		t.Errorf("expected: %v, got: %v", expected, m.hashCalledCount)
	}
}

func (m *mockHashFn) assertCalledWith(t *testing.T, expected map[int][]byte) {
	for i := range m.hashCalledWith {
		if !bytes.Equal(m.hashCalledWith[i], expected[i]) {
			t.Errorf("expected: %v, got: %v", expected[i], m.hashCalledWith[i])
		}
	}
}

func TestHasher_makeKeyPicker(t *testing.T) {
	cases := []struct {
		name             string
		configs          []ConfigBlock
		hashSizeIntBytes int
		mockedReturn     map[int][]byte
		expectedTimes    int
		expectedSource32 []uint32
		expectedSource64 []uint64
	}{
		{
			name: "one block, call hash 1 time",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 2, KeySizeInBits: 8},
			},
			hashSizeIntBytes: 4,
			mockedReturn: map[int][]byte{
				0: {0x11, 0x22, 0x33, 0x44},
			},
			expectedTimes:    1,
			expectedSource32: []uint32{0x44332211},
			expectedSource64: []uint64{0x44332211},
		},
		{
			name: "one block, call hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 4, KeySizeInBits: 16},
			},
			hashSizeIntBytes: 4,
			mockedReturn: map[int][]byte{
				0: {0x11, 0x22, 0x33, 0x44},
				1: {0xaa, 0xbb, 0xcc, 0xdd},
			},
			expectedTimes:    2,
			expectedSource32: []uint32{0x44332211, 0xddccbbaa},
			expectedSource64: []uint64{0xddccbbaa44332211},
		},
		{
			name: "two blocks, same size, call hash 1 time",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 2, KeySizeInBits: 8},
				{NumberOfHashFunctions: 2, KeySizeInBits: 8},
			},
			hashSizeIntBytes: 4,
			mockedReturn: map[int][]byte{
				0: {0x12, 0x34, 0x56, 0x78},
			},
			expectedTimes:    1,
			expectedSource32: []uint32{0x78563412},
			expectedSource64: []uint64{0x78563412},
		},
		{
			name: "two blocks, same size, call hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 2, KeySizeInBits: 16},
				{NumberOfHashFunctions: 2, KeySizeInBits: 16},
			},
			hashSizeIntBytes: 4,
			mockedReturn: map[int][]byte{
				0: {0x12, 0x34, 0x56, 0x78},
				1: {0xa0, 0xb1, 0xc2, 0xd3},
			},
			expectedTimes:    2,
			expectedSource32: []uint32{0x78563412, 0xd3c2b1a0},
			expectedSource64: []uint64{0xd3c2b1a078563412},
		},
		{
			name: "two blocks, different size, call hash 1 time",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 3, KeySizeInBits: 8},
				{NumberOfHashFunctions: 2, KeySizeInBits: 8},
			},
			hashSizeIntBytes: 8,
			mockedReturn: map[int][]byte{
				0: {0x12, 0x34, 0x56, 0x78, 0xa0, 0xb1, 0xc2, 0xd3},
			},
			expectedTimes:    1,
			expectedSource32: []uint32{0x78563412, 0xd3c2b1a0},
			expectedSource64: []uint64{0xd3c2b1a078563412},
		},
		{
			name: "two blocks, different size, call hash 2 times",
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 1, KeySizeInBits: 16},
				{NumberOfHashFunctions: 2, KeySizeInBits: 16},
			},
			hashSizeIntBytes: 4,
			mockedReturn: map[int][]byte{
				0: {0x12, 0x34, 0x56, 0x78},
				1: {0xa0, 0xb1, 0xc2, 0xd3},
			},
			expectedTimes:    2,
			expectedSource32: []uint32{0x78563412, 0xd3c2b1a0},
			expectedSource64: []uint64{0xd3c2b1a078563412},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &hasher{}
			m := &mockHashFn{
				hashReturn: tc.mockedReturn,
			}

			kp := h.makeKeyPicker([]byte("hello"), tc.configs, tc.hashSizeIntBytes, m.Hash)
			m.assertCalledCount(t, tc.expectedTimes)
			assertKeyPickerSourceEquals(t, kp.Source, tc.expectedSource32, tc.expectedSource64)
		})
	}
}

func TestHasher_HashNTimes(t *testing.T) {
	cases := []struct {
		name               string
		n                  byte
		hashSizeInBytes    int
		input              []byte
		mockedReturn       map[int][]byte
		expectedCount      int
		expectedCalledWith map[int][]byte
		expected           []byte
	}{
		{
			name:            "it should call hashFn() one time with original input if n = 1",
			n:               1,
			hashSizeInBytes: 2,
			input:           []byte{1, 2, 3},
			mockedReturn: map[int][]byte{
				0: {100, 101},
			},
			expectedCount: 1,
			expectedCalledWith: map[int][]byte{
				0: {1, 2, 3},
			},
			expected: []byte{100, 101},
		},
		{
			name:            "it should call hashFn() twice with original input and padded 0 if n = 2",
			n:               2,
			hashSizeInBytes: 2,
			input:           []byte{1, 2, 3},
			mockedReturn: map[int][]byte{
				0: {100, 101},
				1: {102, 103},
			},
			expectedCount: 2,
			expectedCalledWith: map[int][]byte{
				0: {1, 2, 3},
				1: {0, 1, 2, 3},
			},
			expected: []byte{100, 101, 102, 103},
		},
		{
			name:            "it should call hashFn() 2 times with original input and padded 0, 1 if n = 3",
			n:               3,
			hashSizeInBytes: 2,
			input:           []byte{1, 2, 3},
			mockedReturn: map[int][]byte{
				0: {100, 101},
				1: {102, 103},
				2: {104, 105},
			},
			expectedCount: 3,
			expectedCalledWith: map[int][]byte{
				0: {1, 2, 3},
				1: {0, 1, 2, 3},
				2: {1, 1, 2, 3},
			},
			expected: []byte{100, 101, 102, 103, 104, 105},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &hasher{}
			m := &mockHashFn{
				hashReturn: tc.mockedReturn,
			}
			result := h.hashNTimes(tc.n, tc.hashSizeInBytes, &tc.input, m.Hash)

			m.assertCalledCount(t, tc.expectedCount)
			m.assertCalledWith(t, tc.expectedCalledWith)
			if !bytes.Equal(tc.expected, result) {
				t.Errorf("expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}

const hasherInput = "hello"

func runTestHasherDoHash(t *testing.T, hasher hasher, hashSizeInBytes int, hashFn hashFn, h, h0, h1, h2, h3 string) {
	cases := []struct {
		name     string
		input    string
		configs  []ConfigBlock
		n        byte
		count    int
		expected string
	}{
		{
			name:  "1 block - 1 time",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 8, KeySizeInBits: byte(hashSizeInBytes)},
			},
			expected: h,
		},
		{
			name:  "1 block - 2 times",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 2, KeySizeInBits: byte(hashSizeInBytes * 7)},
			},
			expected: h + h0,
		},
		{
			name:  "1 block - 3 times",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 4, KeySizeInBits: byte(hashSizeInBytes * 5)},
			},
			expected: h + h0 + h1,
		},
		{
			name:  "1 block - 4 times",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 4, KeySizeInBits: byte(hashSizeInBytes * 7)},
			},
			expected: h + h0 + h1 + h2,
		},
		{
			name:  "1 block - 5 times",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 5, KeySizeInBits: byte(hashSizeInBytes * 7)},
			},
			expected: h + h0 + h1 + h2 + h3,
		},
		{
			name:  "2 blocks - 1 time",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 2, KeySizeInBits: byte(hashSizeInBytes)},
				{NumberOfHashFunctions: 2, KeySizeInBits: byte(hashSizeInBytes)},
			},
			expected: h,
		},
		{
			name:  "2 blocks - 2 times",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 2, KeySizeInBits: byte(hashSizeInBytes * 4)},
				{NumberOfHashFunctions: 2, KeySizeInBits: byte(hashSizeInBytes * 3)},
			},
			expected: h + h0,
		},
		{
			name:  "2 blocks - 4 times",
			input: hasherInput,
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 4, KeySizeInBits: byte(hashSizeInBytes * 4)},
				{NumberOfHashFunctions: 4, KeySizeInBits: byte(hashSizeInBytes * 3)},
			},
			expected: h + h0 + h1 + h2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := hasher.makeKeyPicker([]byte(tc.input), tc.configs, hashSizeInBytes, hashFn).Source
			b, err := hex.DecodeString(tc.expected)
			if err != nil {
				t.Errorf("expected nil err, got %v", err)
			}

			expected := NewKeyPicker(b).Source
			if len(r) != len(expected) {
				t.Errorf("expected length: %v, got %v", len(expected), len(r))
			}
			for i := 0; i < len(r); i++ {
				if r[i] != expected[i] {
					t.Errorf("index %v: expected %v, got %v", i, expected[i], r[i])
				}
			}
		})
	}
}

func TestBuiltinHashers_IsCompatible(t *testing.T) {
	cases := []struct {
		name     string
		left     string
		right    string
		expected bool
	}{
		{
			name:     "not compatible if not the same type - 1",
			left:     "sha",
			right:    "fnv",
			expected: false,
		},
		{
			name:     "not compatible if not the same type - 2",
			left:     "fnv",
			right:    "sha",
			expected: false,
		},
		{
			name:     "not compatible if HashSizeInBytes are not the same - sha",
			left:     "fnv",
			right:    "fnv",
			expected: true,
		},
		{
			name:     "not compatible if HashSizeInBytes are not the same - sha",
			left:     "sha",
			right:    "sha",
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mk := func(t string) Hasher {
				if t == "sha" {
					return &shaHasher{hasher{}}
				}
				return &fnvHasher{hasher{}}
			}
			left := mk(tc.left)
			right := mk(tc.right)
			result := left.IsCompatible(right)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
