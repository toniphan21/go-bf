package bf

import (
	"fmt"
	"testing"
)

func assertKeyPickerSourceEquals(t *testing.T, source []uint, expected32 []uint32, expected64 []uint64) {
	if wordSize == 32 {
		if len(source) != len(expected32) {
			t.Errorf("len(source) != len(expected32)")
		}

		for i := 0; i < len(source); i++ {
			r := fmt.Sprintf("%x", source[i])
			e := fmt.Sprintf("%x", expected32[i])
			if r != e {
				t.Errorf("index %v: expected %v, got %v", i, e, r)
			}
		}
	} else {
		if len(source) != len(expected64) {
			t.Errorf("len(source) != len(expected64)")
		}

		for i := 0; i < len(source); i++ {
			r := fmt.Sprintf("%x", source[i])
			e := fmt.Sprintf("%x", expected64[i])
			if r != e {
				t.Errorf("index %v: expected %v, got %v", i, e, r)
			}
		}
	}
}

func TestNewKeyPicker(t *testing.T) {
	cases := []struct {
		name       string
		source     []byte
		expected32 []uint32
		expected64 []uint64
	}{
		{
			name:       "just 1 byte",
			source:     []byte{0xab},
			expected32: []uint32{0xab},
			expected64: []uint64{0xab},
		},
		{
			name:       "2 bytes",
			source:     []byte{0xab, 0x12},
			expected32: []uint32{0x12ab},
			expected64: []uint64{0x12ab},
		},
		{
			name:       "4 bytes",
			source:     []byte{0xab, 0x12, 0x7d, 0x8a},
			expected32: []uint32{0x8a7d12ab},
			expected64: []uint64{0x8a7d12ab},
		},
		{
			name:       "5 bytes",
			source:     []byte{0xab, 0x12, 0x7d, 0x8a, 0x90},
			expected32: []uint32{0x8a7d12ab, 0x90},
			expected64: []uint64{0x908a7d12ab},
		},
		{
			name:       "8 bytes",
			source:     []byte{0xab, 0x12, 0x7d, 0x8a, 0x90, 0x00, 0xff, 0x34},
			expected32: []uint32{0x8a7d12ab, 0x34ff0090},
			expected64: []uint64{0x34ff00908a7d12ab},
		},
		{
			name:       "10 bytes",
			source:     []byte{0xab, 0x12, 0x7d, 0x8a, 0x90, 0x00, 0xff, 0x34, 0x22, 0x10},
			expected32: []uint32{0x8a7d12ab, 0x34ff0090, 0x1022},
			expected64: []uint64{0x34ff00908a7d12ab, 0x1022},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kp := NewKeyPicker(tc.source)
			assertKeyPickerSourceEquals(t, kp.Source, tc.expected32, tc.expected64)
		})
	}
}

func TestKeyPicker_Pick(t *testing.T) {
	cases := []struct {
		name     string
		source   []byte
		configs  []ConfigBlock
		expected [][]Key
	}{
		{
			name:     "single config block",
			source:   []byte{0xab, 0x12, 0x7d, 0x8a},
			configs:  []ConfigBlock{{NumberOfHashFunctions: 3, KeySizeInBits: 8}},
			expected: [][]Key{{0xab, 0x12, 0x7d}},
		},
		{
			name:   "two same config blocks",
			source: []byte{0xab, 0x12, 0x7d, 0x8a, 0x2d, 0x0a, 0x6e, 0x99, 0x78},
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 3, KeySizeInBits: 12},
				{NumberOfHashFunctions: 3, KeySizeInBits: 12},
			},
			expected: [][]Key{
				{0x2ab, 0x7d1, 0xd8a},
				{0x0a2, 0x96e, 0x789},
			},
		},
		{
			name:   "two different config blocks",
			source: []byte{0xab, 0x12, 0x7d, 0x8a, 0x2d, 0x0a, 0x6e, 0x99, 0x78},
			configs: []ConfigBlock{
				{NumberOfHashFunctions: 3, KeySizeInBits: 12},
				{NumberOfHashFunctions: 2, KeySizeInBits: 16},
			},
			expected: [][]Key{
				{0x2ab, 0x7d1, 0xd8a},
				{0xe0a2, 0x8996},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kp := NewKeyPicker(tc.source)
			result := kp.Pick(tc.configs)
			if len(result) != len(tc.expected) {
				t.Errorf("len(result) != len(tc.expected)")
			}
			for i := 0; i < len(result); i++ {
				if len(result[i]) != len(tc.expected[i]) {
					t.Errorf("len(result[%v]) != len(tc.expected[%v])", i, i)
				}
				for j := 0; j < len(result[i]); j++ {
					if result[i][j] != tc.expected[i][j] {
						t.Errorf("result[%v][%v] != tc.expected[%v][%v]", i, j, i, j)
					}
				}
			}
		})
	}
}

func (p *KeyPicker) pickKeyNaive(size, index int) Key {
	var key Key = 0
	end := index + size
	for i := index; i < end; i++ {
		n := i / wordSize
		m := i % wordSize
		if p.Source[n]&(1<<m) > 0 {
			key |= 1 << (i - index)
		}
	}
	return key
}

func TestKeyPicker_PickKeyNaive(t *testing.T) {
	cases := map[string]struct {
		input  []uint
		size   int
		index  int
		output uint32
	}{
		"same uint, start at 0": {
			input: []uint{0b1101_0101}, size: 5, index: 0, output: 0b1_0101,
		},
		"same uint, start at 3": {
			input: []uint{0b0101_1100_1010_0010_1111_1011}, size: 8, index: 3, output: 0b101_1111,
		},
		"same uint, end at last": {
			input:  []uint{0b1100_0101_1100_1010 << (wordSize - 16)},
			size:   12,
			index:  wordSize - 12,
			output: 0b1100_0101_1100,
		},
		"two uints, end at middle": {
			input:  []uint{0b1100_0101_1100_1010 << (wordSize - 16), 0b1010_0010_1111_1011},
			size:   12,
			index:  wordSize - 6,
			output: 0b11_1011_1100_01,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			kp := &KeyPicker{
				Source: tc.input,
			}
			expected := kp.pickKeyNaive(tc.size, tc.index)
			if expected != Key(tc.output) {
				t.Errorf("expected %d, got %d", tc.output, expected)
			}

			result := kp.pickKey(tc.size, tc.index)
			if result != Key(tc.output) {
				t.Errorf("expected %d, got %d", tc.output, result)
			}
		})
	}
}

func TestKeyPicker_PickKey(t *testing.T) {
	kp := &KeyPicker{
		Source: []uint{0x2cf24dba5fb0a30e, 0x26e83b2ac5b9e29e, 0x1b161e5c1fa7425e, 0x73043362938b9824},
	}
	for i := 0; i < 48; i++ {
		for size := 1; size <= 32; size++ {
			expected := kp.pickKeyNaive(size, i)
			result := kp.pickKey(size, i)
			if result != expected {
				t.Fatalf("size = %v, index = %v: expected %d, got %d", size, i, expected, result)
			}
		}
	}
}

func BenchmarkKeyPicker_PickKey(b *testing.B) {
	kp := &KeyPicker{
		Source: []uint{0x2cf24dba5fb0a30e, 0x26e83b2ac5b9e29e, 0x1b161e5c1fa7425e, 0x73043362938b9824},
	}
	maxStart := 48
	maxSize := 32

	b.Run("naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for start := 0; start < maxStart; start++ {
				for size := 1; size <= maxSize; size++ {
					kp.pickKeyNaive(size, start)
				}
			}
		}
	})

	b.Run("optimal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for start := 0; start < maxStart; start++ {
				for size := 1; size <= maxSize; size++ {
					kp.pickKey(size, start)
				}
			}
		}
	})
}
