package bf

type Option struct {
	config         Config
	storageFactory StorageFactory
	hasherFactory  HasherFactory
}

type OptionFunc func(option *Option)

/*
New BloomFilter instance with Config could be the built-in WithAccuracy or
WithCapacity configuration. Options including WithStorage, WithHasher or a
built-in hash strategy WithSHA (default) and WithFNV.
*/
func New(config Config, opts ...OptionFunc) (BloomFilter, error) {
	o := Option{
		config:         config,
		storageFactory: &memoryStorageFactory{},
		hasherFactory:  &shaHasherFactory{},
	}
	for _, opt := range opts {
		opt(&o)
	}

	storage, err := o.storageFactory.Make(config.StorageCapacity())
	if err != nil {
		return nil, err
	}

	hash := o.hasherFactory.Make(config.NumberOfHashFunctions(), calcKeyMinSizeFromCapacity(config.StorageCapacity()))

	return &bloomFilter{storage: storage, hasher: hash, count: 0}, nil
}

/*
Must create new BloomFilter instance with Config could be the built-in
WithAccuracy or WithCapacity configuration. Options including WithStorage,
WithHasher or a built-in hash strategy WithSHA (default) and WithFNV.
*/
func Must(config Config, opts ...OptionFunc) BloomFilter {
	f, err := New(config, opts...)
	if err != nil {
		panic(err)
	}
	return f
}

func WithStorage(sf StorageFactory) OptionFunc {
	return func(o *Option) {
		o.storageFactory = sf
	}
}

func WithHasher(hf HasherFactory) OptionFunc {
	return func(o *Option) {
		o.hasherFactory = hf
	}
}

func WithSHA() OptionFunc {
	return func(o *Option) {
		o.hasherFactory = &shaHasherFactory{}
	}
}

func WithFNV() OptionFunc {
	return func(o *Option) {
		o.hasherFactory = &fnvHasherFactory{}
	}
}
