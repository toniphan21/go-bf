package bf

import "crypto/sha256"

const shaSize = 32

type shaHasher struct {
	hasher
}

func (s *shaHasher) Hash(input []byte, configs []ConfigBlock) [][]Key {
	kp := s.hasher.makeKeyPicker(input, configs, shaSize, s.doHash)
	return kp.Pick(configs)
}

func (s *shaHasher) IsCompatible(other Hasher) bool {
	_, ok := other.(*shaHasher)
	return ok
}

func (s *shaHasher) doHash(input *[]byte) []byte {
	var result = make([]byte, shaSize)
	src := sha256.Sum256(*input)
	for i := 0; i < shaSize; i++ {
		result[i] = src[i]
	}
	return result
}
