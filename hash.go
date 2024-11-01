package bf

import (
	"crypto/sha256"
	"hash/fnv"
)

type Hash interface {
	Hash([]byte) []Key

	Equals(other Hash) bool
}

type HashFactory interface {
	Make(numberOfHashFunctions, hashSizeInBits byte) Hash
}

type genericHash struct {
	hashSizeInBytes int
	count           byte
	size            int
}

func (g *genericHash) makeKeySplitter(input []byte, hashFn func(*[]byte) []byte) *KeySplitter {
	count := int(g.count)
	var length = count * g.size
	hashSizeInBits := g.hashSizeInBytes * 8
	var times = length / hashSizeInBits
	var mod = length % hashSizeInBits
	if mod > 0 {
		times++
	}
	return &KeySplitter{
		Source:   g.hashNTimes(byte(times), &input, hashFn),
		KeyCount: count,
		KeySize:  g.size,
	}
}

func (g *genericHash) hashNTimes(n byte, input *[]byte, fn func(*[]byte) []byte) []byte {
	if n == 1 {
		return fn(input)
	}

	var result = make([]byte, int(n)*g.hashSizeInBytes)
	for i := byte(0); i < n; i++ {
		l := len(*input)
		item := make([]byte, l+1)
		item[0] = i
		for c := 0; c < l; c++ {
			item[c+1] = (*input)[c]
		}

		src := fn(&item)
		for j := 0; j < g.hashSizeInBytes; j++ {
			result[g.hashSizeInBytes*int(i)+j] = src[j]
		}
	}
	return result
}

type shaHash struct {
	genericHash
}

func (h *shaHash) Equals(other Hash) bool {
	o, ok := other.(*shaHash)
	if !ok {
		return false
	}
	return o.genericHash == h.genericHash
}

const shaSize = 32

func (h *shaHash) Hash(input []byte) []Key {
	return h.genericHash.makeKeySplitter(input, h.doHash).Split()
}

func (h *shaHash) doHash(item *[]byte) []byte {
	var result = make([]byte, shaSize)
	src := sha256.Sum256(*item)
	for i := 0; i < shaSize; i++ {
		result[i] = src[i]
	}
	return result
}

type shaHashFactory struct{}

func (s *shaHashFactory) Make(numberOfHashFunctions, hashSizeInBits byte) Hash {
	return &shaHash{
		genericHash: genericHash{
			hashSizeInBytes: shaSize,
			count:           numberOfHashFunctions,
			size:            int(hashSizeInBits),
		},
	}
}

type fnvHash struct {
	genericHash
}

func (h *fnvHash) Equals(other Hash) bool {
	o, ok := other.(*fnvHash)
	if !ok {
		return false
	}
	return o.genericHash == h.genericHash
}

const fnvSize = 16

func (h *fnvHash) Hash(input []byte) []Key {
	return h.genericHash.makeKeySplitter(input, h.doHash).Split()
}

func (h *fnvHash) doHash(item *[]byte) []byte {
	hash := fnv.New128()
	hash.Write(*item)
	return hash.Sum(nil)
}

type fnvHashFactory struct{}

func (f *fnvHashFactory) Make(numberOfHashFunctions, hashSizeInBits byte) Hash {
	return &fnvHash{
		genericHash: genericHash{
			hashSizeInBytes: fnvSize,
			count:           numberOfHashFunctions,
			size:            int(hashSizeInBits),
		},
	}
}
