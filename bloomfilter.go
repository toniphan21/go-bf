package bf

import "errors"

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Count() int

	Storage() Storage

	Hash() Hasher

	Intersect(other BloomFilter) error

	Union(other BloomFilter) error
}

var ErrStorageDifference = errors.New("storage is not the same")
var ErrHasherDifference = errors.New("hasher is not the same")

type bloomFilter struct {
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

func (b *bloomFilter) Hash() Hasher {
	return b.hasher
}

func (b *bloomFilter) Intersect(other BloomFilter) error {
	if !b.storage.Equals(other.Storage()) {
		return ErrStorageDifference
	}
	if !b.hasher.Equals(other.Hash()) {
		return ErrHasherDifference
	}

	if bi, ok := b.storage.(BatchIntersect); ok {
		bi.Intersect(other.Storage())
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
	if !b.storage.Equals(other.Storage()) {
		return ErrStorageDifference
	}
	if !b.hasher.Equals(other.Hash()) {
		return ErrHasherDifference
	}

	if bi, ok := b.storage.(BatchUnion); ok {
		bi.Union(other.Storage())
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
