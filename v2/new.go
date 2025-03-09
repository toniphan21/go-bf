package bf

type Option struct {
	config    Config
	storage   Storage
	hasher    Hasher
	expansion expansion
}

type OptionFunc func(option *Option)

/*
New BloomFilter instance with Config could be the built-in WithAccuracy or
WithCapacity configuration. Options including WithStorage, WithHasher or a
built-in hash strategy WithSHA (default) and WithFNV.

Default expand when filled ratio reaches 0.5 with expansion rate is 2. The values can be config
using WithExpansion, WithFilledRatio, WithExpansionRate or WithoutExpansion, WithNoExpansion.
*/
func New(config Config, opts ...OptionFunc) (BloomFilter, error) {
	o := Option{
		config:    config,
		storage:   &memoryStorage{},
		hasher:    &shaHasher{},
		expansion: expansion{ratio: DefaultFilledRatio, rate: DefaultExpansionRate},
	}
	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilOptionFunc
		}
		opt(&o)
	}

	r, err := newBloomFilter(o.config, o.hasher, o.storage, o.expansion)
	if err != nil {
		return nil, err
	}
	return r, nil
}

/*
Must create new BloomFilter instance with Config could be the built-in
WithAccuracy or WithCapacity configuration. Options including WithStorage,
WithHasher or a built-in hash strategy WithSHA (default) and WithFNV.

Default expand when filled ratio reaches 0.5 with expansion rate is 2. The values can be config
using WithExpansion, WithFilledRatio, WithExpansionRate or WithoutExpansion, WithNoExpansion.
*/
func Must(config Config, opts ...OptionFunc) BloomFilter {
	f, err := New(config, opts...)
	if err != nil {
		panic(err)
	}
	return f
}

func WithStorage(s Storage) OptionFunc {
	return func(o *Option) {
		o.storage = s
	}
}

func WithHasher(h Hasher) OptionFunc {
	return func(o *Option) {
		o.hasher = h
	}
}

func WithSHA() OptionFunc {
	return func(o *Option) {
		o.hasher = &shaHasher{}
	}
}

func WithFNV() OptionFunc {
	return func(o *Option) {
		o.hasher = &fnvHasher{}
	}
}

func WithExpansion(filledRatio, expansionRate float64) OptionFunc {
	return func(o *Option) {
		o.expansion = expansion{ratio: filledRatio, rate: expansionRate}
	}
}

func WithFilledRatio(v float64) OptionFunc {
	return func(o *Option) {
		o.expansion = expansion{ratio: v, rate: DefaultExpansionRate}
	}
}

func WithExpansionRate(v float64) OptionFunc {
	return func(o *Option) {
		o.expansion = expansion{ratio: DefaultFilledRatio, rate: v}
	}
}

func WithoutExpansion() OptionFunc {
	return func(o *Option) {
		o.expansion = expansion{}
	}
}

func WithNoExpansion() OptionFunc {
	return func(o *Option) {
		o.expansion = expansion{}
	}
}
