package bf

type Key uint32

type KeySplitter struct {
	Source   []byte
	Count    int
	KeyCount int
	KeySize  int
}

func (ks *KeySplitter) Split() [][]Key {
	result := make([][]Key, ks.Count)
	l := uint32(len(ks.Source) * 8)
	for i := 0; i < ks.Count; i++ {
		result[i] = make([]Key, ks.KeyCount)
		offset := uint32(i * ks.KeySize * ks.KeyCount)

		for j := 0; j < ks.KeyCount; j++ {
			var key Key = 0

			for k := 0; k < ks.KeySize; k++ {
				index := offset + uint32(j*ks.KeySize+k)
				if index >= l {
					continue
				}

				n := index / 8
				m := index % 8
				if ks.Source[n]&(1<<m) > 0 {
					key |= 1 << k
				}
			}

			result[i][j] = key
		}
	}
	return result
}
