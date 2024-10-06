package bf

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Capacity() uint

	Count() uint

	Data() []byte

	Info() string
}
