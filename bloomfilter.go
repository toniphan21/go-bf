package bf

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Count() uint

	Data() Storage
}

type bloomFilter struct {
	hash    Hash
	storage Storage
	count   uint
}

func (b *bloomFilter) Add(item []byte) {
	keys := b.hash.Hash(item)
	for _, key := range keys {
		index := key % b.storage.Capacity()
		b.storage.Set(index)
	}
	b.count++
}

func (b *bloomFilter) Exists(item []byte) bool {
	keys := b.hash.Hash(item)
	for _, key := range keys {
		index := key % b.storage.Capacity()
		if !b.storage.Get(index) {
			return false
		}
	}
	return true
}

func (b *bloomFilter) Count() uint {
	return b.count
}

func (b *bloomFilter) Data() Storage {
	return b.storage
}
