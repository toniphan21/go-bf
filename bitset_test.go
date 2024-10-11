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
			name: "one byte", size: 1, capacity: 5, start: 0, end: 8,
		},
		{
			name: "one byte full", size: 1, capacity: 8, start: 0, end: 8,
		},
		{
			name: "two bytes", size: 2, capacity: 10, start: 0, end: 12,
		},
		{
			name: "two bytes full", size: 2, capacity: 16, start: 0, end: 16,
		},
		{
			name: "100 bytes", size: 100, capacity: 700, start: 0, end: 810,
		},
		{
			name: "100 bytes full", size: 100, capacity: 800, start: 0, end: 810,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mut := newBitset(tc.size, tc.capacity)
			var i uint32 = 0
			for i = tc.start; i <= tc.end; i++ {
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
					chars := make([]string, tc.size*8)
					for ci := range chars {
						chars[ci] = "0"
					}
					chars[i] = "1"
					assertStringEqual(t, sprintfBytesInBinary(b.Bytes()), strings.Join(chars, ""))
				} else {
					assertBoolForIndex(t, i, after, false)
				}

				if b.Capacity() != tc.capacity {
					t.Errorf("Expected Capacity %v, got %v", tc.capacity, b.Capacity())
				}
			}
		})
	}
}

func reverseByteBinaryString(b string) string {
	return string([]byte{b[7], b[6], b[5], b[4], b[3], b[2], b[1], b[0]})
}

func sprintfBytesInBinary(b *[]byte) string {
	result := make([]string, len(*b))
	for i, bt := range *b {
		result[i] = reverseByteBinaryString(fmt.Sprintf("%08b", bt))
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
