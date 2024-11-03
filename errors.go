package bf

import "errors"

var ErrInvalidStorageCapacity = errors.New("invalid storage capacity")

var ErrStorageDifference = errors.New("storage is not the same")
var ErrHasherDifference = errors.New("hasher is not the same")

var ErrNilConfig = errors.New("implementation of Config is nil")
var ErrNilOptionFunc = errors.New("implementation of OptionFunc is nil")
var ErrNilStorageFactory = errors.New("implementation of StorageFactory is nil")
var ErrNilHasherFactory = errors.New("implementation of HasherFactory is nil")
var ErrNilStorage = errors.New("implementation of Storage is nil")
var ErrNilHasher = errors.New("implementation of Hasher is nil")
var ErrNilBloomFilter = errors.New("implementation of BloomFilter is nil")
