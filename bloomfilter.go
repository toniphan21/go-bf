package bf

import "errors"

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Count() uint

	Storage() Storage

	Hash() Hash

	Intersect(other BloomFilter) error

	Union(other BloomFilter) error
}

var ErrStorageDifference = errors.New("storage is not the same")
var ErrHashDifference = errors.New("hash is not the same")

type bloomFilter struct {
	hash    Hash
	storage Storage
	count   uint
}

func (b *bloomFilter) Add(item []byte) {
	keys := b.hash.Hash(item)
	for _, key := range keys {
		index := uint32(key) % b.storage.Capacity()
		b.storage.Set(index)
	}
	b.count++
}

func (b *bloomFilter) Exists(item []byte) bool {
	keys := b.hash.Hash(item)
	for _, key := range keys {
		index := uint32(key) % b.storage.Capacity()
		if !b.storage.Get(index) {
			return false
		}
	}
	return true
}

func (b *bloomFilter) Count() uint {
	return b.count
}

func (b *bloomFilter) Storage() Storage {
	return b.storage
}

func (b *bloomFilter) Hash() Hash {
	return b.hash
}

func (b *bloomFilter) Intersect(other BloomFilter) error {
	if !b.storage.Equals(other.Storage()) {
		return ErrStorageDifference
	}
	if !b.hash.Equals(other.Hash()) {
		return ErrHashDifference
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
	return nil
}

func (b *bloomFilter) Union(other BloomFilter) error {
	if !b.storage.Equals(other.Storage()) {
		return ErrStorageDifference
	}
	if !b.hash.Equals(other.Hash()) {
		return ErrHashDifference
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
	return nil
}
