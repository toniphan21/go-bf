package internal

type BloomFilterSpec interface {
	Add([]byte)

	Exists([]byte) bool

	Count() int
}
