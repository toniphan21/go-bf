package bf

import "crypto/sha256"

const shaSize = 32

type shaHasher struct {
	hasher
}

func (s *shaHasher) Hash(input []byte, count int) [][]Key {
	kp := s.hasher.makeKeySplitter(count, input, s.doHash)
	return kp.Split2()
}

func (s *shaHasher) Equals(other Hasher) bool {
	o, ok := other.(*shaHasher)
	if !ok {
		return false
	}
	return o.hasher == s.hasher
}

func (s *shaHasher) doHash(input *[]byte) []byte {
	var result = make([]byte, shaSize)
	src := sha256.Sum256(*input)
	for i := 0; i < shaSize; i++ {
		result[i] = src[i]
	}
	return result
}

type shaHasherFactory struct{}

func (s *shaHasherFactory) Make(numberOfHashFunctions, hashSizeInBits byte) Hasher {
	return &shaHasher{
		hasher: hasher{
			hashSizeInBytes: shaSize,
			keyCount:        numberOfHashFunctions,
			keySize:         int(hashSizeInBits),
		},
	}
}
