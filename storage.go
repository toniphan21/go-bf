package bf

import (
	"errors"
)

var ErrInvalidStorageCapacity = errors.New("invalid storage capacity")

type Storage interface {
	Set(index uint32)

	Clear(index uint32)

	Get(index uint32) bool

	Capacity() uint32

	Equals(other Storage) bool
}

type BatchIntersect interface {
	Intersect(other Storage)
}

type BatchUnion interface {
	Union(other Storage)
}

type StorageFactory interface {
	Make(capacity uint32) (Storage, error)
}

type memoryStorageFactory struct{}

func (msf *memoryStorageFactory) Make(capacity uint32) (Storage, error) {
	if capacity <= 0 {
		return nil, ErrInvalidStorageCapacity
	}

	n, m := capacity/bitsetDataSize, capacity%bitsetDataSize
	if m > 0 {
		n += 1
	}
	return newBitset(n, capacity), nil
}
