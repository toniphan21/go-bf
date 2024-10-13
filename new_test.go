package bf

import (
	"errors"
	"testing"
)

type dummyConfig struct {
	info     string
	k        byte
	capacity uint32
}

func (d *dummyConfig) Info() string {
	return d.info
}

func (d *dummyConfig) NumberOfHashFunctions() byte {
	return d.k
}

func (d *dummyConfig) StorageCapacity() uint32 {
	return d.capacity
}

func TestNew_CallsStorageFactoryAndReturnErrorIfExists(t *testing.T) {
	expected := errors.New("whatever")
	cf := &dummyConfig{k: 10, capacity: 1000}
	storage := &stubStorageFactory{err: expected}

	f, err := New(cf, WithStorage(storage))
	if f != nil {
		t.Errorf("expect filter is nil but got %v", f)
	}
	if err == nil {
		t.Errorf("expect error but got nil")
	}
	if !errors.Is(err, expected) {
		t.Errorf("expect %v but got %v", expected, err)
	}
}

func TestNew_CallsStorageFactoryAndHashFactory(t *testing.T) {
	cf := &dummyConfig{k: 10, capacity: 2000}
	storage := &stubStorageFactory{}
	hash := &stubHashFactory{}

	f, err := New(cf, WithStorage(storage), WithHash(hash))
	if f == nil {
		t.Errorf("expect filter is not nil but got nil")
	}
	if err != nil {
		t.Errorf("expect filter is nil but got %v", err)
	}

	if storage.makeCapacity != cf.capacity {
		t.Errorf("expect %d but got %d", cf.capacity, storage.makeCapacity)
	}
	if hash.makeK != cf.k {
		t.Errorf("expect %d but got %d", cf.k, hash.makeK)
	}
	if hash.makeSize != 11 {
		t.Errorf("expect %d but got %d", 11, hash.makeSize)
	}
}

type stubStorageFactory struct {
	storage      Storage
	err          error
	makeCapacity uint32
}

func (s *stubStorageFactory) Make(capacity uint32) (Storage, error) {
	s.makeCapacity = capacity
	return s.storage, s.err
}

func TestWithStorage(t *testing.T) {
	opt := &Option{}
	ds := &stubStorageFactory{}
	fn := WithStorage(ds)

	fn(opt)
	if opt.storageFactory != ds {
		t.Errorf("Expected storage factory to be %v, got %v", ds, opt.storageFactory)
	}
}

type stubHashFactory struct {
	hash     Hash
	makeK    byte
	makeSize byte
}

func (s *stubHashFactory) Make(numberOfHashFunctions, hashSizeInBits byte) Hash {
	s.makeK = numberOfHashFunctions
	s.makeSize = hashSizeInBits
	return s.hash
}

func TestWithHash(t *testing.T) {
	opt := &Option{}
	dh := &stubHashFactory{}
	fn := WithHash(dh)

	fn(opt)
	if opt.hashFactory != dh {
		t.Errorf("Expected hash factory to be %v, got %v", dh, opt.hashFactory)
	}
}

func TestWithSHA(t *testing.T) {
	opt := &Option{}
	fn := WithSHA()
	fn(opt)

	_, ok := opt.hashFactory.(*shaHashFactory)
	if !ok {
		t.Errorf("Expected hash factory to be shaHashFactory, got %T", opt.hashFactory)
	}
}

func TestWithFNV(t *testing.T) {
	opt := &Option{}
	fn := WithFNV()
	fn(opt)

	_, ok := opt.hashFactory.(*fnvHashFactory)
	if !ok {
		t.Errorf("Expected hash factory to be fnvHashFactory, got %T", opt.hashFactory)
	}
}
