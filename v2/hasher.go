package bf

type Hasher interface {
	Hash(input []byte, configs []ConfigBlock) [][]Key

	Equals(other Hasher) bool
}
