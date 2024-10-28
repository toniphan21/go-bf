## Go Bloom Filter
[![codecov](https://codecov.io/github/toniphan21/go-bf/graph/badge.svg?token=20ILOD9CPG)](https://codecov.io/github/toniphan21/go-bf)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b80cab5415324a4e91ad5cd6cdec1fb0)](https://app.codacy.com/gh/toniphan21/go-bf/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/toniphan21/go-bf)](https://goreportcard.com/report/github.com/toniphan21/go-bf)
[![Go Reference](https://pkg.go.dev/badge/github.com/toniphan21/go-bf?status.svg)](https://pkg.go.dev/github.com/toniphan21/go-bf?tab=doc)

A [Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter) written in GoLang with flexible options and no dependencies.

### Usage

Initialize a Bloom Filter `WithAccuracy` configuration
```golang
package main

import "github.com/toniphan21/go-bf"

func main() {
	var errorRate = 0.001
	var numberOfItems uint32 = 10_000_000
	filter := bf.Must(bf.WithAccuracy(errorRate, numberOfItems))

	filter.Add([]byte("anything"))

	if !filter.Exists([]byte("not found")) {
		println("If a bloom filter returns false, it is 100% correct.")
	}

	if filter.Exists([]byte("anything")) {
		println("If a bloom filter returns true it MAYBE correct. Remember to check false positive cases.")
	}
}

```

Initialize a Bloom Filter `WithCapacity` configuration

```golang
package main

import "github.com/toniphan21/go-bf"

func main() {
	var capacityInBits uint32 = 65_536
	var numberOfHashFunctions byte = 5
	filter := bf.Must(bf.WithCapacity(capacityInBits, numberOfHashFunctions))

	filter.Add([]byte("anything"))

	if !filter.Exists([]byte("not found")) {
		println("If a bloom filter returns false, it is 100% correct.")
	}

	if filter.Exists([]byte("anything")) {
		println("If a bloom filter returns true it MAYBE correct. Remember to check false positive cases.")
	}
}
```

### APIs

#### Constructors

- `New(Config, ...Options) (BloomFilter, error)` initialize new instance
- `Must(Config, ...Options) BloomFilter` initialize new instance, panic if encounter any error

#### BloomFilter interface

The `BloomFilter` interface has 4 main methods:

| Method                | Description                               |
|-----------------------|-------------------------------------------|
| `Add([]byte)`         | Add an item into the filter               |
| `Exists([]byte) bool` | Check existence of an item in the filter  |
| `Count() int`         | Get number of items added into the filter |
| `Data() Storage`      | Get filter's Storage                      |


#### Options

There are 4 option functions could be used from the second param of `bf.New(Config, ...OptionFunc)`:

| Signature                       |           | Description                                            |
|---------------------------------|-----------|--------------------------------------------------------|
| `WithSHA()`                     | _default_ | Use splitted SHA hashing strategy (more uniform hash)  |
| `WithFNV()`                     |           | Use splitted FNV hashing strategy (better performance) |
| `WithHash(f HashFactory)`       |           | Customize Hashing strategy with a HashFactory          |
| `WithStorage(f StorageFactory)` |           | Customize Storage strategy with a StorageFactory       |


### Implementation Details

#### Error rate, number of hash functions calculation

_This is just a summary of how to calculate estimated error rate of a Bloom Filter, for proof please check
academic paper or wikipedia page._

Given:

- `n` estimated number of items in a Bloom Filter
- `e` false positive error rate
- `k` number of hash functions needed
- `m` bits of memory

If you have `n` and `e` (`WithAccuracy` config):

```math
\huge k \approx -log_2 e
```


```math
\huge m \approx -1.44 n log_2 e
```


`WithAccuracy` has a builtin `.Info()` returns calculated `m`, `k` and rounded values.

```golang
package main

import (
	"fmt"
	"github.com/toniphan21/go-bf"
)

func main() {
	var errorRate = 0.001
	var numberOfItems uint32 = 10_000_000
	config := bf.WithAccuracy(errorRate, numberOfItems)
	fmt.Println(config.Info())
}

/** OUTPUT:

Config WithAccuracy()
  - Requested error rate: 0.10000%
  - Expected number of items: 10000000
  - Bits per item: 14.351
  - Number of hash functions: 10
  - Size in bits of each has function: 28
  - Storage capacity: 143507294 bits = 17938412 bytes = 17517.98KB = 17.11MB
  - Estimated error rate: 0.10130%

*/
```

If you have `m` and `k` (`WithCapacity` config), estimated error rate when `n` items are added into a Bloom Filter:


```math
\huge \epsilon \approx (1 - e^-\frac{k n}{m} )^k
```


`WithCapacity` configuration has a builtin `.Info()` to display the estimated error rate when config by `m` and `k`. 

```golang
package main

import (
	"fmt"
	"github.com/toniphan21/go-bf"
)

func main() {
	var capacityInBits uint32 = 65_536
	var numberOfHashFunctions byte = 5
	config := bf.WithCapacity(capacityInBits, numberOfHashFunctions)
	fmt.Println(config.Info())
}

/** OUTPUT:

Config WithCapacity()
  - Storage capacity: 65536 bits = 8192 bytes = 8.00KB = 0.01MB
  - Number of hash functions: 5
  - Size in bits of each has function: 16
  - Estimated error rate by n - number of added items:
      n=   100; estimated error rate: 0.00000%
      n=   200; estimated error rate: 0.00000%
      n=   500; estimated error rate: 0.00001%
      n=  1000; estimated error rate: 0.00021%
      n=  2000; estimated error rate: 0.00568%
      n=  5000; estimated error rate: 0.32083%
      n= 10000; estimated error rate: 4.33023%
      n= 20000; estimated error rate: 29.35056%
      n= 50000; estimated error rate: 89.45317%
      n=100000; estimated error rate: 99.75726%
      n=200000; estimated error rate: 99.99988%
      n=500000; estimated error rate: 100.00000%

 */
```

> If you spot something wrong with the calculation, no worries - you can [write your own config](https://github.com/toniphan21/go-bf?tab=readme-ov-file#write-your-own-config).

#### Hashing strategy

This library has builtin 2 hash functions with the same strategy:

- From the config we could know: `size` minimum key size (in bits) and `count` number of hash needed.
- Use `SHA-256` or `FNV-128` to generate hash bytes from the input. If `size * count` > `256` when use SHA (or `128`
  when use FNV), will do the hash multiple times with a byte prefixed
- Pick `size` bits from the hash bytes in previous step, discard all remaining bits.

Example 1: `size = 25`, `count = 10`, use `SHA-256`:

- Because 25*10 = 250 bits, we only need to hash 1 time
- hash = `sha_hash(input)`
- pick key 0 = bit 0-24
- pick key 1 = bit 25-49
- ...
- pick key 9 = bit 225-249
- bit 250-255 is discarded
- return `[10]uint32{key0, key1...key9}`

Example 2: `size = 25`, `count = 10`, use `FNV-128`:

- Because 25*10 > 128, we hash input 2 times
- hash = `fnv_128(byte(0) + input)` + `fnv_128(byte(1) + input)`
- pick key 0 = bit 0-24
- pick key 1 = bit 25-49
- ...
- pick key 9 = bit 225-249
- bit 250-255 is discarded
- return `[10]uint32{key0, key1...key9}`

### Customization

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
	filter := bf.Must(config, bf.WithHash(&YourHashFactory{}))

	filter.Add([]byte("anything"))
	// ...
}
```

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
	filter := bf.Must(config, bf.WithStorage(&FileStorageFactory{}))

	filter.Add([]byte("anything"))
	// ...
}
```

#### Write your own config

If you don't like `WithCapacity()` or `WithAccuracy()` configuration, you can write your own:

```golang
package main

import "github.com/toniphan21/go-bf"

type YourConfig struct {
}

func (y *YourConfig) Info() string {
	return "info about your config"
}

func (y *YourConfig) NumberOfHashFunctions() byte {
	return 5
}

func (y *YourConfig) StorageCapacity() uint32 {
	return 1_000_000
}

func main() {
	config := &YourConfig{}
	filter, err := bf.New(config)
	if err != nil {
		panic("Something went wrong")
	}

	filter.Add([]byte("anything"))
	// ...
}
```

### Benchmark

```
BenchmarkBloomFilter_WithSHA_Add-12       	 1026219	      1120  ns/op
BenchmarkBloomFilter_WithFNV_Add-12       	 2048646	      593.9 ns/op
BenchmarkBloomFilter_WithSHA_Exists-12    	 1000000	      1127  ns/op
BenchmarkBloomFilter_WithFNV_Exists-12    	 2114071	      569.9 ns/op
```
