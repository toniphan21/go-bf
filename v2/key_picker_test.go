package bf

import (
	"testing"
)

func (p *KeyPicker) pickKeyNaive(size, index int) Key {
	var key Key = 0
	end := index + size
	for i := index; i < end; i++ {
		n := i / uintSize
		m := i % uintSize
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
			input:  []uint{0b1100_0101_1100_1010 << (uintSize - 16)},
			size:   12,
			index:  uintSize - 12,
			output: 0b1100_0101_1100,
		},
		"two uints, end at middle": {
			input:  []uint{0b1100_0101_1100_1010 << (uintSize - 16), 0b1010_0010_1111_1011},
			size:   12,
			index:  uintSize - 6,
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
