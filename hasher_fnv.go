package bf

import "hash/fnv"

const fnvSize = 16

type fnvHasher struct {
	hasher
}

func (s *fnvHasher) Hash(input []byte, count int) [][]Key {
	kp := s.hasher.makeKeySplitter(count, input, s.doHash)
	return kp.Split2()
}

func (s *fnvHasher) Equals(other Hasher) bool {
	o, ok := other.(*fnvHasher)
	if !ok {
		return false
	}
	return o.hasher == s.hasher
}

func (s *fnvHasher) doHash(input *[]byte) []byte {
	hash := fnv.New128()
	hash.Write(*input)
	return hash.Sum(nil)
}

type fnvHasherFactory struct{}

func (s *fnvHasherFactory) Make(numberOfHashFunctions, hashSizeInBits byte) Hasher {
	return &fnvHasher{
		hasher: hasher{
			hashSizeInBytes: fnvSize,
			keyCount:        numberOfHashFunctions,
			keySize:         int(hashSizeInBits),
		},
	}
}
