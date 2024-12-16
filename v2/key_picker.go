package bf

import (
	"math"
)

type Key uint32

const wordSize = 32 << (^uint(0) >> 63) // 32 or 64

func NewKeyPicker(source []byte) *KeyPicker {
	l := len(source)
	m := wordSize >> 3
	size := l / m
	if l%m != 0 {
		size++
	}

	kp := &KeyPicker{
		Source: make([]uint, size),
	}
	for i := 0; i < l; i++ {
		index := i / m
		mod := i % m

		kp.Source[index] |= uint(source[i]) << (mod * 8)
	}
	return kp
}

type KeyPicker struct {
	Source []uint
}

func (p *KeyPicker) Pick(configs []ConfigBlock) [][]Key {
	var result = make([][]Key, len(configs))
	index := 0
	for i, cf := range configs {
		k := int(cf.NumberOfHashFunctions)
		size := int(cf.KeySizeInBits)

		result[i] = make([]Key, k)
		for j := 0; j < k; j++ {
			result[i][j] = p.pickKey(size, index)
			index += size
		}
	}
	return result
}

func (p *KeyPicker) pickKey(size, index int) Key {
	end := index + size
	indexStart := index / wordSize
	indexEnd := end / wordSize

	remainderStart := index % wordSize
	remainderEnd := end % wordSize

	var maskStart uint = math.MaxUint << remainderStart
	var maskEnd uint = math.MaxUint >> (wordSize - remainderEnd)
	if indexStart == indexEnd {
		return Key(p.Source[indexStart] & maskStart & maskEnd >> remainderStart)
	}

	result := Key(p.Source[indexStart] & maskStart >> remainderStart)
	if remainderEnd == 0 {
		return result
	}
	return result | Key(p.Source[indexEnd]&maskEnd)<<(wordSize-remainderStart)
}
