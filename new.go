package bf

type BloomFilter interface {
	Add(item []byte)

	Exists(item []byte) bool

	Count() uint

	Data() Storage
}

type Option struct {
	config         Config
	storageFactory StorageFactory
	hashFactory    HashFactory
}

type OptionFunc func(option *Option)

type bloomFilter struct {
	hash    Hash
	storage Storage
	count   uint
}

func (b *bloomFilter) Add(item []byte) {
	keys := b.hash.Hash(item)
	for _, key := range keys {
		index := key % b.storage.Capacity()
		b.storage.Set(index)
	}
	b.count++
}

func (b *bloomFilter) Exists(item []byte) bool {
	keys := b.hash.Hash(item)
	for _, key := range keys {
		index := key % b.storage.Capacity()
		if !b.storage.Get(index) {
			return false
		}
	}
	return true
}

func (b *bloomFilter) Count() uint {
	return b.count
}

func (b *bloomFilter) Data() Storage {
	return b.storage
}

func New(config Config, opts ...OptionFunc) (BloomFilter, error) {
	o := Option{
		config:         config,
		storageFactory: &memoryStorageFactory{},
		hashFactory:    &shaHashFactory{},
	}
	for _, opt := range opts {
		opt(&o)
	}

	storage, err := o.storageFactory.Make(config.StorageCapacity())
	if err != nil {
		return nil, err
	}

	hash := o.hashFactory.Make(config.NumberOfHashFunctions(), calcKeyMinSizeFromCapacity(config.StorageCapacity()))

	return &bloomFilter{storage: storage, hash: hash, count: 0}, nil
}

func WithStorage(sf StorageFactory) OptionFunc {
	return func(o *Option) {
		o.storageFactory = sf
	}
}

func WithHash(hf HashFactory) OptionFunc {
	return func(o *Option) {
		o.hashFactory = hf
	}
}

func WithSHA() OptionFunc {
	return func(o *Option) {
		o.hashFactory = &shaHashFactory{}
	}
}
