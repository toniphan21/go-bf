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
	if config == nil {
		return nil, ErrNilConfig
	}

	o := Option{
		config:         config,
		storageFactory: memoryStorageFactory{},
		hasherFactory:  shaHasherFactory{},
	}
	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilOptionFunc
		}
		opt(&o)
	}

	if o.storageFactory == nil {
		return nil, ErrNilStorageFactory
	}

	if o.hasherFactory == nil {
		return nil, ErrNilHasherFactory
	}

	r, err := newBloomFilter(o)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func newBloomFilter(o Option) (*bloomFilter, error) {
	storage, err := o.storageFactory.Make(o.config.StorageCapacity())
	if err != nil {
		return nil, err
	}
	if storage == nil {
		return nil, ErrNilStorage
	}

	h := o.hasherFactory.Make(o.config.NumberOfHashFunctions(), o.config.KeySize())
	if h == nil {
		return nil, ErrNilHasher
	}
	return &bloomFilter{option: o, storage: storage, hasher: h, count: 0}, nil
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
		o.hasherFactory = shaHasherFactory{}
	}
}

func WithFNV() OptionFunc {
	return func(o *Option) {
		o.hasherFactory = fnvHasherFactory{}
	}
}
