package bf

type Option struct {
	config         Config
	storageFactory StorageFactory
	hashFactory    HashFactory
}

type OptionFunc func(option *Option)

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

func WithFNV() OptionFunc {
	return func(o *Option) {
		o.hashFactory = &fnvHashFactory{}
	}
}
