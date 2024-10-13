package main

import (
	"fmt"
	"github.com/toniphan21/go-bf"
)

func infoOfWithCapacity() {
	var capacityInBits uint32 = 65_536
	var numberOfHashFunctions byte = 5
	config := bf.WithCapacity(capacityInBits, numberOfHashFunctions)
	fmt.Println(config.Info())
}

func infoOfWithAccuracy() {
	var errorRate = 0.001
	var numberOfItems uint32 = 10_000_000
	config := bf.WithAccuracy(errorRate, numberOfItems)
	fmt.Println(config.Info())
}

type YourHash struct {
	count byte
	size  byte
}

func (y *YourHash) Hash(bytes []byte) []uint32 {
	// return an array of hash for given bytes input.
	//   - length of the array is count - number of hash functions
	//   - each hash need to >= size - minimum size of a hash in bits
	return []uint32{}
}

type YourHashFactory struct{}

func (y *YourHashFactory) Make(numberOfHashFunctions, hashSizeInBits byte) bf.Hash {
	return &YourHash{
		count: numberOfHashFunctions,
		size:  hashSizeInBits,
	}
}

func customHash() {
	config := bf.WithAccuracy(0.01, 1_000_000)
	filter, err := bf.New(config, bf.WithHash(&YourHashFactory{}))
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))
	// ...
}

type FileStorage struct {
	capacity uint32
}

func (f *FileStorage) Set(index uint32) {
	// set a bit in the given index to true
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

func customStorage() {
	config := bf.WithAccuracy(0.01, 1_000_000)
	filter, err := bf.New(config, bf.WithStorage(&FileStorageFactory{}))
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))
	// ...
}

func newWithCapacity() {
	var capacityInBits uint32 = 65_536
	var numberOfHashFunctions byte = 5
	filter, err := bf.New(bf.WithCapacity(capacityInBits, numberOfHashFunctions))
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))

	if !filter.Exists([]byte("not found")) {
		println("If a bloom filter returns false, it is 100% correct.")
	}

	if filter.Exists([]byte("anything")) {
		println("If a bloom filter returns true it MAYBE correct. Remember to check false positive cases.")
	}
}

func newWithAccuracy() {
	var errorRate = 0.001
	var numberOfItems uint32 = 10_000_000
	filter, err := bf.New(bf.WithAccuracy(errorRate, numberOfItems))
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))

	if !filter.Exists([]byte("not found")) {
		println("If a bloom filter returns false, it is 100% correct.")
	}

	if filter.Exists([]byte("anything")) {
		println("If a bloom filter returns true it MAYBE correct. Remember to check false positive cases.")
	}
}
