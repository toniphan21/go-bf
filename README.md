## Go Bloom Filter

A classic [Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter) written in GoLang with flexible options and no dependencies.

### Usage

Initialize a Bloom Filter WithAccuracy configuration
```golang
package main

import "github.com/toniphan21/go-bf"

func main() {
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

```

Initialize a Bloom Filter WithCapacity configuration

```golang
package main

import "github.com/toniphan21/go-bf"

func main() {
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
```

### Options

There are 4 option functions could be used from the second param of `bf.New(config Config, opts ...OptionFunc)`:

| Signature                             |         | Description                                            |
|---------------------------------------|---------|--------------------------------------------------------|
| `WithSHA()`                           | default | Use splitted SHA hashing strategy (more uniform hash)  |
| `WithFNV()`                           |         | Use splitted FNV hashing strategy (better performance) |
| `WithHash(factory HashFactory)`       |         | Customize Hashing strategy with a HashFactory          |
| `WithStorage(factory StorageFactory)` |         | Customize Storage strategy with a StorageFactory       |


### Customization

#### Write your own custom storage

By default, all data are stored in memory, you can customize a storage by implement `Storage` and `StorageFactory` 
interface:

```golang
package main

import "github.com/toniphan21/go-bf"

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

func main() {
	config := bf.WithAccuracy(0.01, 1_000_000)
	filter, err := bf.New(config, bf.WithStorage(&FileStorageFactory{}))
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))
	// ...
}
```

#### Write your own hashing strategy

By default, an SHA hash splitted by number of key and key size are used. You can customize the Hash functions by
implement `Hash` and `HashFactory` interface:

```golang
package main

import "github.com/toniphan21/go-bf"

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

func main() {
	config := bf.WithAccuracy(0.01, 1_000_000)
	filter, err := bf.New(config, bf.WithHash(&YourHashFactory{}))
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))
	// ...
}
```

### Implementation Details

#### Error rate, number of hash functions calculation

#### Hashing strategy

### Licence

MIT.
