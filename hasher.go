package bf

type HasherFactory interface {
	Make(numberOfHashFunctions, hashSizeInBits byte) Hasher
}

type Hasher interface {
	Hash(input []byte, size int) [][]Key

	Equals(other Hasher) bool
}
