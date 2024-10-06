package storage

import "errors"

var InvalidCapacity = errors.New("invalid capacity")

type Storage interface {
	Set(index uint)

	Get(index uint) bool

	Capacity() uint

	Bytes() *[]byte
}

func New(capacity int) (Storage, error) {
	if capacity <= 0 {
		return nil, InvalidCapacity
	}

	c := uint(capacity)
	n, m := c/8, c%8
	if m > 0 {
		n += 1
	}
	return newBitset(n, c), nil
}
