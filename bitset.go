package bf

type bitset struct {
	data     []byte
	capacity uint32
}

func newBitset(n, capacity uint32) *bitset {
	return &bitset{data: make([]byte, n), capacity: capacity}
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
	d := b.data[n] & (m ^ 0xFF)
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

func (b *bitset) indexing(i uint32) (uint32, byte) {
	n := i / 8
	m := i % 8

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
