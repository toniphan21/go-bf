package bf

type ConfigBlock struct {
	Capacity              uint32
	NumberOfHashFunctions byte
	KeySize               byte
}

type Config interface {
	Info() string

	MaxFillRatio() float64

	ExpansionRate() float64

	Get(numberOfBlocks int) []ConfigBlock
}
