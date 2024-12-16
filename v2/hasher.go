package bf

type Hasher interface {
	Hash(input []byte, configs []ConfigBlock) [][]Key

	IsCompatible(other Hasher) bool
}

type hashFn func(*[]byte) []byte

type hasher struct{}

func (h *hasher) makeKeyPicker(input []byte, configs []ConfigBlock, hashSizeInBytes int, hashFn hashFn) *KeyPicker {
	var length = 0
	for _, config := range configs {
		length += int(config.KeySizeInBits) * int(config.NumberOfHashFunctions)
	}

	hashSizeInBits := hashSizeInBytes * 8
	var times = length / hashSizeInBits
	var mod = length % hashSizeInBits
	if mod > 0 {
		times++
	}
	return NewKeyPicker(h.hashNTimes(byte(times), hashSizeInBytes, &input, hashFn))
}

func (h *hasher) hashNTimes(n byte, hashSizeInBytes int, input *[]byte, fn hashFn) []byte {
	if n == 1 {
		return fn(input)
	}

	var result = make([]byte, int(n)*hashSizeInBytes)
	l := len(*input)
	item := make([]byte, l+1)
	for c := 0; c < l; c++ {
		item[c+1] = (*input)[c]
	}

	for i := byte(0); i < n; i++ {
		if i == 0 {
			src := fn(input)
			for j := 0; j < hashSizeInBytes; j++ {
				result[hashSizeInBytes*int(i)+j] = src[j]
			}
			continue
		}

		item[0] = i - 1
		src := fn(&item)
		for j := 0; j < hashSizeInBytes; j++ {
			result[hashSizeInBytes*int(i)+j] = src[j]
		}
	}
	return result
}
