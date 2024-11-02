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

func (d *dummyConfig) KeySize() byte {
	return calcKeyMinSizeFromCapacity(d.capacity)
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

func TestNew_CallsStorageFactoryAndHasherFactory(t *testing.T) {
	cf := &dummyConfig{k: 10, capacity: 2000}
	storage := &stubStorageFactory{}
	hasher := &stubHasherFactory{}

	f, err := New(cf, WithStorage(storage), WithHasher(hasher))
	if f == nil {
		t.Errorf("expect filter is not nil but got nil")
	}
	if err != nil {
		t.Errorf("expect filter is nil but got %v", err)
	}

	if storage.makeCapacity != cf.capacity {
		t.Errorf("expect %d but got %d", cf.capacity, storage.makeCapacity)
	}
	if hasher.makeK != cf.k {
		t.Errorf("expect %d but got %d", cf.k, hasher.makeK)
	}
	if hasher.makeSize != 11 {
		t.Errorf("expect %d but got %d", 11, hasher.makeSize)
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

type stubHasherFactory struct {
	hasher   Hasher
	makeK    byte
	makeSize byte
}

func (s *stubHasherFactory) Make(numberOfHashFunctions, hashSizeInBits byte) Hasher {
	s.makeK = numberOfHashFunctions
	s.makeSize = hashSizeInBits
	return s.hasher
}

func TestWithHasher(t *testing.T) {
	opt := &Option{}
	dh := &stubHasherFactory{}
	fn := WithHasher(dh)

	fn(opt)
	if opt.hasherFactory != dh {
		t.Errorf("Expected hash factory to be %v, got %v", dh, opt.hasherFactory)
	}
}

func TestWithSHA(t *testing.T) {
	opt := &Option{}
	fn := WithSHA()
	fn(opt)

	_, ok := opt.hasherFactory.(shaHasherFactory)
	if !ok {
		t.Errorf("Expected hash factory to be shaHasherFactory, got %T", opt.hasherFactory)
	}
}

func TestWithFNV(t *testing.T) {
	opt := &Option{}
	fn := WithFNV()
	fn(opt)

	_, ok := opt.hasherFactory.(fnvHasherFactory)
	if !ok {
		t.Errorf("Expected hash factory to be fnvHasherFactory, got %T", opt.hasherFactory)
	}
}

func TestMustPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Must() to panic, but it did not")
		}
	}()

	expected := errors.New("whatever")
	cf := &dummyConfig{k: 10, capacity: 1000}
	storage := &stubStorageFactory{err: expected}

	Must(cf, WithStorage(storage))
}

func TestMustDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Did not expect Must() to panic, but it did")
		}
	}()

	cf := &dummyConfig{k: 10, capacity: 1000}
	storage := &stubStorageFactory{}

	Must(cf, WithStorage(storage))
}
