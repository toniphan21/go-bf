package main

import "github.com/toniphan21/go-bf"

type FileStorage struct {
	capacity uint32
}

func (f *FileStorage) Equals(other bf.Storage) bool {
	o, ok := other.(*FileStorage)
	if !ok {
		return false
	}
	// check other params
	return o.capacity == f.capacity
}

func (f *FileStorage) Set(index uint32) {
	// set a bit in the given index to true
}

func (f *FileStorage) Clear(index uint32) {
	// clear a bit in the given index to true
}

func (f *FileStorage) Get(index uint32) bool {
	// return a boolean in the given index
	return false
}

func (f *FileStorage) Capacity() uint32 {
	// return the capacity of the storage in bits
	return f.capacity
}

type FileStorageFactory struct{}

func (f *FileStorageFactory) Make(capacity uint32) (bf.Storage, error) {
	return &FileStorage{capacity}, nil
}

func main() {
	config := bf.WithAccuracy(0.01, 1_000_000)
	filter := bf.Must(config, bf.WithStorage(&FileStorageFactory{}))

	filter.Add([]byte("anything"))
	// ...
}
