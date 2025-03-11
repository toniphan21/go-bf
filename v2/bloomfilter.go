package bf

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Count() int

	Storage() Storage

	Hasher() Hasher

	Intersect(other BloomFilter) error

	Union(other BloomFilter) error

	Clone() BloomFilter
}

type expansion struct {
	ratio float64
	rate  float64
}

type bloomFilter struct {
	config           Config
	currentConfig    Config
	expansion        expansion
	hasher           Hasher
	storage          Storage
	configBlocks     []ConfigBlock
	count            int
	bitSetByBlocks   []int
	maxBitSetByBlock []int
}

func toConfigBlock(config Config) ConfigBlock {
	return ConfigBlock{
		Capacity:              config.StorageCapacity(),
		NumberOfHashFunctions: config.NumberOfHashFunctions(),
		KeySizeInBits:         config.KeySize(),
	}
}

func newBloomFilter(config Config, hasher Hasher, storage Storage, expansion expansion) (*bloomFilter, error) {
	if config == nil {
		return nil, ErrNilConfig
	}
	if hasher == nil {
		return nil, ErrNilHasher
	}
	if storage == nil {
		return nil, ErrNilStorage
	}

	storage.NewBlock(config.StorageCapacity())
	b := &bloomFilter{
		config:           config,
		currentConfig:    config,
		expansion:        expansion,
		storage:          storage,
		hasher:           hasher,
		count:            0,
		bitSetByBlocks:   []int{0},
		maxBitSetByBlock: []int{int(float64(config.StorageCapacity()) * expansion.ratio)},
		configBlocks:     []ConfigBlock{toConfigBlock(config)},
	}
	return b, nil
}

func (b *bloomFilter) Add(item []byte) {
	lastBlock := len(b.configBlocks) - 1
	keySets := b.hasher.Hash(item, b.configBlocks)
	for i := 0; i < lastBlock; i++ {
		if b.existsInBlock(i, keySets[i]) {
			return
		}
	}

	storageBlock := b.storage.Block(lastBlock)
	exists := true
	for _, key := range keySets[lastBlock] {
		index := uint32(key) % storageBlock.Capacity()
		if !storageBlock.Get(index) {
			exists = false
			storageBlock.Set(index)
			b.bitSetByBlocks[lastBlock]++
		}
	}

	if !exists {
		b.count++
	}

	if b.expansion.rate > 0 && b.bitSetByBlocks[lastBlock] > b.maxBitSetByBlock[lastBlock] {
		next := b.currentConfig.Next(b.expansion.rate)
		b.currentConfig = next

		b.bitSetByBlocks = append(b.bitSetByBlocks, 0)
		b.maxBitSetByBlock = append(b.maxBitSetByBlock, int(float64(next.StorageCapacity())*b.expansion.ratio))
		b.storage.NewBlock(next.StorageCapacity())
		b.configBlocks = append(b.configBlocks, toConfigBlock(next))
	}
}

func (b *bloomFilter) Exists(item []byte) bool {
	keySets := b.hasher.Hash(item, b.configBlocks)
	for i, ks := range keySets {
		if b.existsInBlock(i, ks) {
			return true
		}
	}
	return false
}

func (b *bloomFilter) existsInBlock(block int, keys []Key) bool {
	storageBlock := b.storage.Block(block)
	for _, key := range keys {
		index := uint32(key) % storageBlock.Capacity()
		if !storageBlock.Get(index) {
			return false
		}
	}
	return true
}

func (b *bloomFilter) Count() int {
	return b.count
}

func (b *bloomFilter) Storage() Storage {
	return b.storage
}

func (b *bloomFilter) Hasher() Hasher {
	return b.hasher
}

func (b *bloomFilter) assertOtherBloomFilterIsTheSame(other BloomFilter) error {
	if other == nil {
		return ErrNilBloomFilter
	}

	if !b.storage.IsCompatible(other.Storage()) {
		return ErrStorageDifference
	}

	if !b.hasher.IsCompatible(other.Hasher()) {
		return ErrHasherDifference
	}
	return nil
}

func (b *bloomFilter) Intersect(other BloomFilter) error {
	if err := b.assertOtherBloomFilterIsTheSame(other); err != nil {
		return err
	}

	for i := 0; i < b.storage.BlockCount(); i++ {
		cBlock := b.storage.Block(i)
		oBlock := other.Storage().Block(i)
		if cbi, ok := cBlock.(BatchIntersect); ok {
			cbi.Intersect(oBlock)
			continue
		}

		for j := uint32(0); j < cBlock.Capacity(); j++ {
			if !cBlock.Get(j) || !oBlock.Get(j) {
				cBlock.Clear(j)
			}
		}
	}
	b.count = -1
	return nil
}

func (b *bloomFilter) Union(other BloomFilter) error {
	if err := b.assertOtherBloomFilterIsTheSame(other); err != nil {
		return err
	}

	for i := 0; i < b.storage.BlockCount(); i++ {
		cBlock := b.storage.Block(i)
		oBlock := other.Storage().Block(i)
		if cbi, ok := cBlock.(BatchUnion); ok {
			cbi.Union(oBlock)
			continue
		}

		for j := uint32(0); j < cBlock.Capacity(); j++ {
			if cBlock.Get(j) || oBlock.Get(j) {
				cBlock.Set(j)
			}
		}
	}
	b.count = -1
	return nil
}

func (b *bloomFilter) Clone() BloomFilter {
	cloned := &bloomFilter{
		config:           b.config,
		currentConfig:    b.currentConfig,
		expansion:        b.expansion,
		storage:          b.storage.Clone(),
		hasher:           b.hasher,
		count:            b.count,
		bitSetByBlocks:   make([]int, len(b.bitSetByBlocks)),
		maxBitSetByBlock: make([]int, len(b.maxBitSetByBlock)),
		configBlocks:     make([]ConfigBlock, len(b.configBlocks)),
	}

	copy(cloned.bitSetByBlocks, b.bitSetByBlocks)
	copy(cloned.maxBitSetByBlock, b.maxBitSetByBlock)
	copy(cloned.configBlocks, b.configBlocks)

	return cloned
}

var _ BloomFilter = (*bloomFilter)(nil)
