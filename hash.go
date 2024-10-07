package bf

import (
	"crypto/sha256"
)

type Hash interface {
	Hash([]byte) []uint32
}

type HashFactory interface {
	Make(numberOfHashFunctions, hashSizeInBits byte) Hash
}

const shaSize = 32

type shaHash struct {
	count byte
	size  byte
}

func (h *shaHash) Hash(input []byte) []uint32 {
	var length = int(h.count) * int(h.size)
	var times = length / 256
	var mod = length % 256
	if mod > 0 {
		times++
	}
	_ = bitset{
		data:     h.doHash(byte(times), &input),
		capacity: uint32(times * shaSize * 8),
	}

	return make([]uint32, 0)
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
