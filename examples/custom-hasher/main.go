package main

import "github.com/toniphan21/go-bf"

type YourHasher struct {
	keyCount byte
	keySize  byte
}

func (y *YourHasher) Equals(other bf.Hasher) bool {
	o, ok := other.(*YourHasher)
	if !ok {
		return false
	}
	// check other params
	return y.keyCount == o.keyCount && y.keySize == o.keySize
}

func (y *YourHasher) Hash(input []byte, count int) [][]bf.Key {
	// return an array of a hash array for given bytes input.
	//   - length of the array is count
	//   - each of subarray will have length = keyCount - number of hash functions
	//   - each hash need to >= keySize - minimum size of a hash in bits
	// For example: given keyCount = 5, keySize = 16
	//   - count = 1 requires you to returns:
	//     [][]Key{
	//       { "key0: at least 16 bits long", "key1:...", "key2:...", "key3:...", "key4:..."},
	//     }
	//   - count = 2 requires you to returns:
	//     [][]Key{
	//       { "key0: at least 16 bits long", "key1:...", "key2:...", "key3:...", "key4:..."},
	//       { "key5: at least 16 bits long", "key6:...", "key7:...", "key7:...", "key8:..."},
	//     }
	return [][]bf.Key{}
}

type YourHasherFactory struct{}

func (y *YourHasherFactory) Make(numberOfHashFunctions, hashSizeInBits byte) bf.Hasher {
	return &YourHasher{
		keyCount: numberOfHashFunctions,
		keySize:  hashSizeInBits,
	}
}

func main() {
	config := bf.WithAccuracy(0.01, 1_000_000)
	filter := bf.Must(config, bf.WithHasher(&YourHasherFactory{}))

	filter.Add([]byte("anything"))
	// ...
}
