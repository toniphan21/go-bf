package bf

import "hash/fnv"

const fnvSize = 16

type fnvHasher struct {
	hasher
}

func (s *fnvHasher) Hash(input []byte, configs []ConfigBlock) [][]Key {
	kp := s.hasher.makeKeyPicker(input, configs, fnvSize, s.doHash)
	return kp.Pick(configs)
}

func (s *fnvHasher) IsCompatible(other Hasher) bool {
	_, ok := other.(*fnvHasher)
	return ok
}

func (s *fnvHasher) doHash(input *[]byte) []byte {
	hash := fnv.New128()
	hash.Write(*input)
	return hash.Sum(nil)
}
