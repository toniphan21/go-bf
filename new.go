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
}

type OptionFunc func(option *Option)

func New(config Config, opts ...OptionFunc) (BloomFilter, error) {
	o := Option{
		config:         config,
		storageFactory: &memoryStorageFactory{},
	}
	for _, opt := range opts {
		opt(&o)
	}

	return nil, nil
}

func WithStorage(sf StorageFactory) OptionFunc {
	return func(o *Option) {
		o.storageFactory = sf
	}
}
