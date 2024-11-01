package bf

type HasherFactory interface {
	Make(numberOfHashFunctions, hashSizeInBits byte) Hasher
}

type Hasher interface {
	Hash(input []byte, count int) [][]Key

	Equals(other Hasher) bool
}

type hasher struct {
	hashSizeInBytes int
	keyCount        byte
	keySize         int
}

func (h *hasher) makeKeySplitter(count int, input []byte, hashFn func(*[]byte) []byte) *KeySplitter {
	keyCount := int(h.keyCount)
	var length = count * keyCount * h.keySize
	hashSizeInBits := h.hashSizeInBytes * 8
	var times = length / hashSizeInBits
	var mod = length % hashSizeInBits
	if mod > 0 {
		times++
	}
	return &KeySplitter{
		Source:   h.hashNTimes(byte(times), &input, hashFn),
		Count:    count,
		KeyCount: keyCount,
		KeySize:  h.keySize,
	}
}

func (h *hasher) hashNTimes(n byte, input *[]byte, fn func(*[]byte) []byte) []byte {
	if n == 1 {
		return fn(input)
	}

	var result = make([]byte, int(n)*h.hashSizeInBytes)
	l := len(*input)
	item := make([]byte, l+1)
	for c := 0; c < l; c++ {
		item[c+1] = (*input)[c]
	}

	for i := byte(0); i < n; i++ {
		if i == 0 {
			src := fn(input)
			for j := 0; j < h.hashSizeInBytes; j++ {
				result[h.hashSizeInBytes*int(i)+j] = src[j]
			}
			continue
		}

		item[0] = i - 1
		src := fn(&item)
		for j := 0; j < h.hashSizeInBytes; j++ {
			result[h.hashSizeInBytes*int(i)+j] = src[j]
		}
	}
	return result
}
