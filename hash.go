package bf

import (
	"crypto/sha256"
	"encoding/binary"
)

type Hash interface {
	Hash([]byte) []uint32
}

type HashFactory interface {
	Make(numberOfHashFunctions, hashSizeInBits byte) Hash
}

type KeySplitter struct {
	Source   []byte
	Length   int
	KeyCount int
	KeySize  int
}

func (k *KeySplitter) Split() []uint32 {
	bs := bitset{
		data:     k.Source,
		capacity: uint32(k.Length),
	}
	result := make([]uint32, k.KeyCount)
	for i := 0; i < k.KeyCount; i++ {
		rbs := bitset{data: make([]byte, 4), capacity: 32}
		for j := 0; j < k.KeySize; j++ {
			index := uint32(i*k.KeySize + j)
			if bs.Get(index) {
				rbs.Set(uint32(j))
			}
		}
		result[i] = binary.LittleEndian.Uint32(rbs.data)
	}
	return result
}

const shaSize = 32

type shaHash struct {
	count byte
	size  byte
}

func (h *shaHash) Hash(input []byte) []uint32 {
	count := int(h.count)
	size := int(h.size)
	var length = count * size
	var times = length / 256
	var mod = length % 256
	if mod > 0 {
		times++
	}
	kp := &KeySplitter{
		Source:   h.doHash(byte(times), &input),
		Length:   times * shaSize * 8,
		KeyCount: count,
		KeySize:  size,
	}
	return kp.Split()
}

func (h *shaHash) doHash(n byte, input *[]byte) []byte {
	if n == 1 {
		var result = make([]byte, shaSize)
		src := sha256.Sum256(*input)
		for i := 0; i < shaSize; i++ {
			result[i] = src[i]
		}
		return result
	}

	var result = make([]byte, n*shaSize)
	var i byte = 0
	for i = 0; i < n; i++ {
		l := len(*input)
		item := make([]byte, l+1)
		item[0] = i
		for c := 0; c < l; c++ {
			item[c+1] = (*input)[c]
		}

		src := sha256.Sum256(item)
		for j := 0; j < shaSize; j++ {
			result[shaSize*int(i)+j] = src[j]
		}
	}
	return result
}
