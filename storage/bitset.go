package storage

type bitset struct {
	data     []byte
	capacity uint
}

func newBitset(n, capacity uint) *bitset {
	return &bitset{data: make([]byte, n), capacity: capacity}
}

func (b *bitset) Capacity() uint {
	return b.capacity
}

func (b *bitset) Set(index uint) {
	if index >= b.capacity {
		return
	}
	n, m := b.Indexing(index)
	d := b.data[n] | m
	b.data[n] = d
}

func (b *bitset) Get(index uint) bool {
	if index >= b.capacity {
		return false
	}
	n, m := b.Indexing(index)
	d := b.data[n] & m
	return d > 0
}

func (b *bitset) Bytes() *[]byte {
	return &b.data
}

func (b *bitset) Indexing(index uint) (uint, byte) {
	n := index / 8
	m := index % 8

	return n, 1 << m
}
