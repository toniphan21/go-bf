package bf

import "math"

type memoryStorageBlock struct {
	data     []uint
	capacity uint32
}

func newMemoryStorageBlock(n, capacity uint32) *memoryStorageBlock {
	return &memoryStorageBlock{data: make([]uint, n), capacity: capacity}
}

func (b *memoryStorageBlock) Capacity() uint32 {
	return b.capacity
}

func (b *memoryStorageBlock) Set(index uint32) {
	if index >= b.capacity {
		return
	}

	n, m := b.indexing(index)
	d := b.data[n] | m
	b.data[n] = d
}

func (b *memoryStorageBlock) Clear(index uint32) {
	if index >= b.capacity {
		return
	}

	n, m := b.indexing(index)
	d := b.data[n] & (m ^ math.MaxUint)
	b.data[n] = d
}

func (b *memoryStorageBlock) Get(index uint32) bool {
	if index >= b.capacity {
		return false
	}

	n, m := b.indexing(index)
	d := b.data[n] & m
	return d > 0
}

func (b *memoryStorageBlock) Equals(other StorageBlock) bool {
	o, ok := other.(*memoryStorageBlock)
	if !ok {
		return false
	}
	return o.capacity == b.capacity
}

func (b *memoryStorageBlock) indexing(i uint32) (uint32, uint) {
	n := i / uintSize
	m := i % uintSize

	return n, 1 << m
}

func (b *memoryStorageBlock) Intersect(other StorageBlock) {
	o, ok := other.(*memoryStorageBlock)
	if !ok {
		return
	}

	l := len(b.data)
	for i := 0; i < l; i++ {
		b.data[i] &= o.data[i]
	}
}

func (b *memoryStorageBlock) Union(other StorageBlock) {
	o, ok := other.(*memoryStorageBlock)
	if !ok {
		return
	}

	l := len(b.data)
	for i := 0; i < l; i++ {
		b.data[i] |= o.data[i]
	}
}
