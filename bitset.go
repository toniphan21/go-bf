package bf

import "math"

const bitsetDataSize = 32 << (^uint(0) >> 63) // 32 or 64

type bitset struct {
	data     []uint
	capacity uint32
}

func newBitset(n, capacity uint32) *bitset {
	return &bitset{data: make([]uint, n), capacity: capacity}
}

func (b *bitset) Capacity() uint32 {
	return b.capacity
}

func (b *bitset) Set(index uint32) {
	if index >= b.capacity {
		return
	}

	n, m := b.indexing(index)
	d := b.data[n] | m
	b.data[n] = d
}

func (b *bitset) Clear(index uint32) {
	if index >= b.capacity {
		return
	}

	n, m := b.indexing(index)
	d := b.data[n] & (m ^ math.MaxUint)
	b.data[n] = d
}

func (b *bitset) Get(index uint32) bool {
	if index >= b.capacity {
		return false
	}

	n, m := b.indexing(index)
	d := b.data[n] & m
	return d > 0
}

func (b *bitset) Equals(other Storage) bool {
	o, ok := other.(*bitset)
	if !ok {
		return false
	}
	return o.capacity == b.capacity
}

func (b *bitset) indexing(i uint32) (uint32, uint) {
	n := i / bitsetDataSize
	m := i % bitsetDataSize

	return n, 1 << m
}

func (b *bitset) Intersect(other Storage) {
	o, ok := other.(*bitset)
	if !ok {
		return
	}

	l := len(b.data)
	for i := 0; i < l; i++ {
		b.data[i] &= o.data[i]
	}
}

func (b *bitset) Union(other Storage) {
	o, ok := other.(*bitset)
	if !ok {
		return
	}

	l := len(b.data)
	for i := 0; i < l; i++ {
		b.data[i] |= o.data[i]
	}
}
