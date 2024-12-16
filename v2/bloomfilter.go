package bf

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Count() int

	Storage() Storage

	Hasher() Hasher

	Intersect(other BloomFilter) error

	Union(other BloomFilter) error

	Clone() (BloomFilter, error)
}

type bloomFilter struct {
	config           Config
	currentConfig    Config
	expansion        Expansion
	hasher           Hasher
	storage          Storage
	configBlocks     []ConfigBlock
	count            int
	countByBlocks    []int
	maxCountByBlocks []int
}

func toConfigBlock(config Config) ConfigBlock {
	return ConfigBlock{
		Capacity:              config.StorageCapacity(),
		NumberOfHashFunctions: config.NumberOfHashFunctions(),
		KeySizeInBits:         config.KeySize(),
	}
}

func newBloomFilter(config Config, hasher Hasher, storage Storage, expansion Expansion) (*bloomFilter, error) {
	if config == nil {
		return nil, ErrNilConfig
	}
	if hasher == nil {
		return nil, ErrNilHasher
	}
	if storage == nil {
		return nil, ErrNilStorage
	}
	if expansion == nil {
		return nil, ErrNilExpansion
	}

	b := &bloomFilter{
		config:           config,
		currentConfig:    config,
		expansion:        expansion,
		storage:          storage,
		hasher:           hasher,
		count:            0,
		countByBlocks:    []int{0},
		maxCountByBlocks: []int{int(float64(config.StorageCapacity()) * expansion.ExpansionRatio())},
		configBlocks:     []ConfigBlock{toConfigBlock(config)},
	}
	return b, nil
}

func (b *bloomFilter) Add(item []byte) {
	lastBlock := len(b.configBlocks) - 1
	keySets := b.hasher.Hash(item, b.configBlocks)
	storageBlock := b.storage.Block(lastBlock)

	exists := true
	for _, key := range keySets[lastBlock] {
		index := uint32(key) % storageBlock.Capacity()
		if !storageBlock.Get(index) {
			exists = false
		}
		storageBlock.Set(index)
	}

	if exists {
		return
	}

	b.count++
	b.countByBlocks[lastBlock]++

	if b.expansion.ExpansionRate() != 0 && b.countByBlocks[lastBlock] > b.maxCountByBlocks[lastBlock] {
		next := b.currentConfig.Next(b.expansion.ExpansionRate())
		b.currentConfig = next

		b.countByBlocks = append(b.countByBlocks, 0)
		b.maxCountByBlocks = append(b.maxCountByBlocks, int(float64(next.StorageCapacity())*b.expansion.ExpansionRatio()))
		b.storage.NewBlock(next.StorageCapacity())
		b.configBlocks = append(b.configBlocks, toConfigBlock(next))
	}
}

func (b *bloomFilter) Exists(item []byte) bool {
	keySets := b.hasher.Hash(item, b.configBlocks)
	for i, ks := range keySets {
		storageBlock := b.storage.Block(i)
		exists := true

		for _, key := range ks {
			index := uint32(key) % storageBlock.Capacity()
			if !storageBlock.Get(index) {
				exists = false
				break
			}
		}

		if exists {
			return true
		}
	}
	return false
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

func (b *bloomFilter) Intersect(other BloomFilter) error {
	//TODO implement me
	panic("implement me")
}

func (b *bloomFilter) Union(other BloomFilter) error {
	//TODO implement me
	panic("implement me")
}

func (b *bloomFilter) Clone() (BloomFilter, error) {
	//TODO implement me
	panic("implement me")
}

var _ BloomFilter = (*bloomFilter)(nil)

//func getConfigBlock(config Config, expansionCf ExpansionConfig, numberOfBlocks int) []ConfigBlock {
//	expansionRate := expansionCf.ExpansionRate()
//	if expansionCf == nil || expansionRate == 0 || numberOfBlocks < 1 {
//		return []ConfigBlock{
//			{
//				Capacity:              config.StorageCapacity(),
//				NumberOfHashFunctions: config.NumberOfHashFunctions(),
//				KeySizeInBits:         config.KeySize(),
//			},
//		}
//	}
//
//	result := make([]ConfigBlock, numberOfBlocks)
//	current := config
//	for i := 0; i < numberOfBlocks; i++ {
//		result[i] = ConfigBlock{
//			Capacity:              current.StorageCapacity(),
//			NumberOfHashFunctions: current.NumberOfHashFunctions(),
//			KeySizeInBits:         current.KeySize(),
//		}
//		current = current.Next(expansionRate)
//	}
//	return result
//}
