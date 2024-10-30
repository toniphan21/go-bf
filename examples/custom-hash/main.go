package main

import "github.com/toniphan21/go-bf"

type YourHash struct {
	count byte
	size  byte
}

func (y *YourHash) Equals(other bf.Hash) bool {
	o, ok := other.(*YourHash)
	if !ok {
		return false
	}
	// check other params
	return y.count == o.count && y.size == o.size
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

func main() {
	config := bf.WithAccuracy(0.01, 1_000_000)
	filter := bf.Must(config, bf.WithHash(&YourHashFactory{}))

	filter.Add([]byte("anything"))
	// ...
}
