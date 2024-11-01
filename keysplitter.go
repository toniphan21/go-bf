package bf

type Key uint32

type KeySplitter struct {
	Source   []byte
	Length   int
	KeyCount int
	KeySize  int
}

func (k *KeySplitter) Split() []Key {
	result := make([]Key, k.KeyCount)
	l := uint32(len(k.Source) * 8)
	for i := 0; i < k.KeyCount; i++ {
		var key Key = 0

		for j := 0; j < k.KeySize; j++ {
			index := uint32(i*k.KeySize + j)
			if index >= l {
				continue
			}

			n := index / 8
			m := index % 8
			if k.Source[n]&(1<<m) > 0 {
				key |= 1 << j
			}
		}

		result[i] = key
	}
	return result
}
