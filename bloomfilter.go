package bf

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Count() int

	Storage() Storage

	Hasher() Hasher

	Intersect(other BloomFilter) error

	Union(other BloomFilter) error

	Clone() (BloomFilter, error)
}

type bloomFilter struct {
	option  Option
	hasher  Hasher
	storage Storage
	count   int
}

func (b *bloomFilter) Add(item []byte) {
	keys := b.hasher.Hash(item, 1)
	for _, key := range keys[0] {
		index := uint32(key) % b.storage.Capacity()
		b.storage.Set(index)
	}
	b.count++
}

func (b *bloomFilter) Exists(item []byte) bool {
	keys := b.hasher.Hash(item, 1)
	for _, key := range keys[0] {
		index := uint32(key) % b.storage.Capacity()
		if !b.storage.Get(index) {
			return false
		}
	}
	return true
}

func (b *bloomFilter) Count() int {
	return b.count
}

func (b *bloomFilter) Storage() Storage {
	return b.storage
}

func (b *bloomFilter) Hasher() Hasher {
	return b.hasher
}

func (b *bloomFilter) assertOtherBloomFilterIsTheSame(other BloomFilter) error {
	if other == nil {
		return ErrNilBloomFilter
	}

	if !b.storage.Equals(other.Storage()) {
		return ErrStorageDifference
	}

	if !b.hasher.Equals(other.Hasher()) {
		return ErrHasherDifference
	}
	return nil
}

func (b *bloomFilter) Intersect(other BloomFilter) error {
	if err := b.assertOtherBloomFilterIsTheSame(other); err != nil {
		return err
	}

	if bi, ok := b.storage.(BatchIntersect); ok {
		bi.Intersect(other.Storage())
		b.count = -1
		return nil
	}

	oStorage := other.Storage()
	for i := uint32(0); i < oStorage.Capacity(); i++ {
		if !b.storage.Get(i) || !oStorage.Get(i) {
			b.storage.Clear(i)
		}
	}
	b.count = -1
	return nil
}

func (b *bloomFilter) Union(other BloomFilter) error {
	if err := b.assertOtherBloomFilterIsTheSame(other); err != nil {
		return err
	}

	if bi, ok := b.storage.(BatchUnion); ok {
		bi.Union(other.Storage())
		b.count = -1
		return nil
	}

	oStorage := other.Storage()
	for i := uint32(0); i < oStorage.Capacity(); i++ {
		if b.storage.Get(i) || oStorage.Get(i) {
			b.storage.Set(i)
		}
	}
	b.count = -1
	return nil
}

func (b *bloomFilter) Clone() (BloomFilter, error) {
	r, err := newBloomFilter(b.option)
	if err != nil {
		return nil, err
	}

	err = r.Union(b)
	if err != nil {
		return nil, err
	}
	r.count = b.count
	return r, nil
}
