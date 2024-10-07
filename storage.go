package bf

import (
	"errors"
)

var InvalidStorageCapacity = errors.New("invalid storage capacity")

type Storage interface {
	Set(index uint32)

	Get(index uint32) bool

	Capacity() uint32

	Bytes() *[]byte
}

type StorageFactory interface {
	Make(capacity uint32) (Storage, error)
}

type memoryStorageFactory struct{}

func (msf *memoryStorageFactory) Make(capacity uint32) (Storage, error) {
	if capacity <= 0 {
		return nil, InvalidStorageCapacity
	}

	n, m := capacity/8, capacity%8
	if m > 0 {
		n += 1
	}
	return newBitset(n, capacity), nil
}
