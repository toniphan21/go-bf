package bf

import (
	"fmt"
	"strings"
	"testing"
)

type mockStorageBlock struct {
	setIndex   []uint32
	clearIndex []uint32
	getData    map[uint32]bool
	capacity   uint32
}

func (m *mockStorageBlock) Set(index uint32) {
	m.setIndex = append(m.setIndex, index)
}

func (m *mockStorageBlock) Clear(index uint32) {
	m.clearIndex = append(m.clearIndex, index)
}

func (m *mockStorageBlock) Get(index uint32) bool {
	l := len(m.getData)
	if index >= uint32(l) {
		return false
	}
	return m.getData[index]
}

func (m *mockStorageBlock) Capacity() uint32 {
	return m.capacity
}

func (m *mockStorageBlock) Equals(other StorageBlock) bool {
	o, ok := other.(*mockStorageBlock)
	if !ok {
		return false
	}
	return o.capacity == m.capacity
}

func (m *mockStorageBlock) assertSetCalledWith(t *testing.T, indices []uint32) {
	if len(m.setIndex) != len(indices) {
		t.Errorf("Set is not called with %v", indices)
	}
	for i, b := range indices {
		if b != m.setIndex[i] {
			t.Errorf("Set is not called with %v", indices)
		}
	}
}

func (m *mockStorageBlock) assertClearCalledWith(t *testing.T, indices []uint32) {
	if len(m.clearIndex) != len(indices) {
		t.Errorf("Clear is not called with %v", indices)
	}
	for i, b := range indices {
		if b != m.clearIndex[i] {
			t.Errorf("Clear is not called with %v", indices)
		}
	}
}

func TestMemoryStorageBlock(t *testing.T) {
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
			name: "one word full", size: 1, capacity: uintSize, start: 0, end: uintSize,
		},
		{
			name: "two words", size: 2, capacity: uintSize + 2, start: 0, end: uintSize*2 + 4,
		},
		{
			name: "two words full", size: 2, capacity: uintSize * 2, start: 0, end: uintSize * 2,
		},
		{
			name: "100 words", size: 100, capacity: uintSize * 90, start: 0, end: uintSize * 101,
		},
		{
			name: "100 words full", size: 100, capacity: uintSize * 100, start: 0, end: uintSize * 110,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mut := newMemoryStorageBlock(tc.size, tc.capacity)
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

				b := newMemoryStorageBlock(tc.size, tc.capacity)
				before = b.Get(i)
				b.Set(i)
				after = b.Get(i)

				assertBoolForIndex(t, i, before, false)
				if i < tc.capacity {
					assertBoolForIndex(t, i, after, true)
					chars := make([]string, tc.size*uintSize)
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

				c := newMemoryStorageBlock(tc.size, tc.capacity)
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

func TestMemoryStorageBlock_Equals_ReturnsFalseIfItIsNotMemoryStorageBlock(t *testing.T) {
	b := &memoryStorageBlock{capacity: 1}
	o := &mockStorageBlock{capacity: 1}

	if b.Equals(o) {
		t.Errorf("Expected false, got true")
	}
}

func TestMemoryStorageBlock_Equals_ReturnsFalseIfCapacityIsNotEqual(t *testing.T) {
	b := &memoryStorageBlock{capacity: 1}
	o := &memoryStorageBlock{capacity: 2}

	if b.Equals(o) {
		t.Errorf("Expected false, got true")
	}
}

func TestMemoryStorageBlock_Equals_ReturnsTrueIfItIsAMemoryStorageBlockAndHaveSameCapacity(t *testing.T) {
	b := &memoryStorageBlock{capacity: 1}
	o := &memoryStorageBlock{capacity: 1}

	if !b.Equals(o) {
		t.Errorf("Expected true, got false")
	}
}

func TestMemoryStorageBlock_Intersect_DoesNothingIfStorageIsNotAMemoryStorageBlock(t *testing.T) {
	b := &memoryStorageBlock{data: []uint{1, 2}}
	o := &mockStorageBlock{getData: map[uint32]bool{}}
	b.Intersect(o)
	if b.data[0] != 1 || b.data[1] != 2 {
		t.Errorf("Expected do nothing something changed")
	}
}

func TestMemoryStorageBlock_Intersect(t *testing.T) {
	a := &memoryStorageBlock{data: []uint{0, 2, 0b00110011}}
	b := &memoryStorageBlock{data: []uint{1, 0, 0b01010101}}
	a.Intersect(b)
	if a.data[0] != 0 || a.data[1] != 0 || a.data[2] != 0b00010001 {
		t.Errorf("Intersect should apply AND operator to all bytes")
	}
	if b.data[0] != 1 || b.data[1] != 0 || b.data[2] != 0b01010101 {
		t.Errorf("Intersect should not changed the given Storage data")
	}
}

func TestMemoryStorageBlock_Union_DoesNothingIfStorageIsNotAMemoryStorageBlock(t *testing.T) {
	b := &memoryStorageBlock{data: []uint{1, 2}}
	o := &mockStorageBlock{getData: map[uint32]bool{}}
	b.Union(o)
	if b.data[0] != 1 || b.data[1] != 2 {
		t.Errorf("Expected do nothing something changed")
	}
}

func TestMemoryStorageBlock_Union(t *testing.T) {
	a := &memoryStorageBlock{data: []uint{0, 0, 2, 0b00110011}}
	b := &memoryStorageBlock{data: []uint{0, 1, 0, 0b01010101}}
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
	format := fmt.Sprintf("%%0%db", uintSize)
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
