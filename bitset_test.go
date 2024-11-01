package bf

import (
	"fmt"
	"strings"
	"testing"
)

func TestBitset(t *testing.T) {
	cases := []struct {
		name     string
		size     uint32
		capacity uint32
		start    uint32
		end      uint32
	}{
		{
			name: "one word", size: 1, capacity: 5, start: 0, end: 8,
		},
		{
			name: "one word full", size: 1, capacity: bitsetDataSize, start: 0, end: bitsetDataSize,
		},
		{
			name: "two words", size: 2, capacity: bitsetDataSize + 2, start: 0, end: bitsetDataSize*2 + 4,
		},
		{
			name: "two words full", size: 2, capacity: bitsetDataSize * 2, start: 0, end: bitsetDataSize * 2,
		},
		{
			name: "100 words", size: 100, capacity: bitsetDataSize * 90, start: 0, end: bitsetDataSize * 101,
		},
		{
			name: "100 words full", size: 100, capacity: bitsetDataSize * 100, start: 0, end: bitsetDataSize * 110,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mut := newBitset(tc.size, tc.capacity)
			for i := tc.start; i <= tc.end; i++ {
				before := mut.Get(i)
				mut.Set(i)
				after := mut.Get(i)

				assertBoolForIndex(t, i, before, false)
				if i < tc.capacity {
					assertBoolForIndex(t, i, after, true)
				} else {
					assertBoolForIndex(t, i, after, false)
				}

				b := newBitset(tc.size, tc.capacity)
				before = b.Get(i)
				b.Set(i)
				after = b.Get(i)

				assertBoolForIndex(t, i, before, false)
				if i < tc.capacity {
					assertBoolForIndex(t, i, after, true)
					chars := make([]string, tc.size*bitsetDataSize)
					for ci := range chars {
						chars[ci] = "0"
					}
					chars[i] = "1"
					assertStringEqual(t, sprintfUintInBinary(&b.data), strings.Join(chars, ""))
				} else {
					assertBoolForIndex(t, i, after, false)
				}

				if b.Capacity() != tc.capacity {
					t.Errorf("Expected Capacity %v, got %v", tc.capacity, b.Capacity())
				}

				c := newBitset(tc.size, tc.capacity)
				c.Set(i)
				before = c.Get(i)
				c.Clear(i)
				after = c.Get(i)

				if i < tc.capacity {
					assertBoolForIndex(t, i, before, true)
				} else {
					assertBoolForIndex(t, i, before, false)
				}
				assertBoolForIndex(t, i, after, false)
			}
		})
	}
}

func TestBitset_Equals_ReturnsFalseIfItIsNotBitset(t *testing.T) {
	b := &bitset{capacity: 1}
	o := &mockStorage{capacity: 1}

	if b.Equals(o) {
		t.Errorf("Expected false, got true")
	}
}

func TestBitset_Equals_ReturnsFalseIfCapacityIsNotEqual(t *testing.T) {
	b := &bitset{capacity: 1}
	o := &bitset{capacity: 2}

	if b.Equals(o) {
		t.Errorf("Expected false, got true")
	}
}

func TestBitset_Equals_ReturnsTrueIfItIsABitsetAndHaveSameCapacity(t *testing.T) {
	b := &bitset{capacity: 1}
	o := &bitset{capacity: 1}

	if !b.Equals(o) {
		t.Errorf("Expected true, got false")
	}
}

func TestBitset_Intersect_DoesNothingIfStorageIsNotABitset(t *testing.T) {
	b := &bitset{data: []uint{1, 2}}
	o := &mockStorage{getData: map[uint32]bool{}}
	b.Intersect(o)
	if b.data[0] != 1 || b.data[1] != 2 {
		t.Errorf("Expected do nothing something changed")
	}
}

func TestBitset_Intersect(t *testing.T) {
	a := &bitset{data: []uint{0, 2, 0b00110011}}
	b := &bitset{data: []uint{1, 0, 0b01010101}}
	a.Intersect(b)
	if a.data[0] != 0 || a.data[1] != 0 || a.data[2] != 0b00010001 {
		t.Errorf("Intersect should apply AND operator to all bytes")
	}
	if b.data[0] != 1 || b.data[1] != 0 || b.data[2] != 0b01010101 {
		t.Errorf("Intersect should not changed the given Storage data")
	}
}

func TestBitset_Union_DoesNothingIfStorageIsNotABitset(t *testing.T) {
	b := &bitset{data: []uint{1, 2}}
	o := &mockStorage{getData: map[uint32]bool{}}
	b.Union(o)
	if b.data[0] != 1 || b.data[1] != 2 {
		t.Errorf("Expected do nothing something changed")
	}
}

func TestBitset_Union(t *testing.T) {
	a := &bitset{data: []uint{0, 0, 2, 0b00110011}}
	b := &bitset{data: []uint{0, 1, 0, 0b01010101}}
	a.Union(b)
	if a.data[0] != 0 || a.data[1] != 1 || a.data[2] != 2 && a.data[3] != 0b01110111 {
		t.Errorf("Intersect should apply OR operator to all bytes")
	}
	if b.data[0] != 0 || b.data[1] != 1 || b.data[2] != 0 && b.data[3] != 0b01010101 {
		t.Errorf("Intersect should not changed the given Storage data")
	}
}

func reverseByteBinaryString(b string) string {
	n := len(b)
	sb := strings.Builder{}
	sb.Grow(n)

	for i := n - 1; i >= 0; i-- {
		sb.WriteByte(b[i])
	}
	return sb.String()
}

func sprintfUintInBinary(b *[]uint) string {
	result := make([]string, len(*b))
	format := fmt.Sprintf("%%0%db", bitsetDataSize)
	for i, bt := range *b {
		result[i] = reverseByteBinaryString(fmt.Sprintf(format, bt))
	}
	return strings.Join(result, "")
}

func assertBoolForIndex(t *testing.T, index uint32, result, expected bool) {
	if result != expected {
		t.Errorf("Index %v: Expected %t, got %t", index, expected, result)
	}
}

func assertStringEqual(t *testing.T, result string, expected string) {
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
