package bf

import (
	"errors"
	"testing"
)

type dummyConfig struct {
	info     string
	k        byte
	s        byte
	capacity uint32
}

func (d *dummyConfig) KeySize() byte {
	return d.s
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

func (d *dummyConfig) Next(expansionRate float64) Config {
	return &dummyConfig{
		info:     d.info,
		k:        d.k,
		capacity: uint32(float64(d.capacity) * expansionRate),
	}
}

func assertNewFailedWithError(t *testing.T, f BloomFilter, err error, expected error) {
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

func TestNew_ShouldCheckNilConfig(t *testing.T) {
	f, err := New(nil)
	assertNewFailedWithError(t, f, err, ErrNilConfig)
}

func TestNew_ShouldCheckNilOptionFunc(t *testing.T) {
	cf := &dummyConfig{k: 10, capacity: 1000}
	f, err := New(cf, nil)
	assertNewFailedWithError(t, f, err, ErrNilOptionFunc)
}

func TestNew_ShouldCheckNilStorage(t *testing.T) {
	cf := &dummyConfig{k: 10, capacity: 1000}
	f, err := New(cf, WithStorage(nil))
	assertNewFailedWithError(t, f, err, ErrNilStorage)
}

func TestNew_ShouldCheckNilHasher(t *testing.T) {
	cf := &dummyConfig{k: 10, capacity: 1000}
	f, err := New(cf, WithHasher(nil))
	assertNewFailedWithError(t, f, err, ErrNilHasher)
}

func TestMustPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Must() to panic, but it did not")
		}
	}()

	cf := &dummyConfig{k: 10, capacity: 1000}

	Must(cf, WithStorage(nil))
}

func TestMustDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Did not expect Must() to panic, but it did")
		}
	}()

	cf := &dummyConfig{k: 10, capacity: 1000}

	Must(cf)
}

func TestWithStorage(t *testing.T) {
	opt := &Option{}
	ds := &mockStorage{}
	fn := WithStorage(ds)

	fn(opt)
	if opt.storage != ds {
		t.Errorf("Expected storage to be %v, got %v", ds, opt.storage)
	}
}

func TestWithHasher(t *testing.T) {
	opt := &Option{}
	dh := &dummyHasher{}
	fn := WithHasher(dh)

	fn(opt)
	if opt.hasher != dh {
		t.Errorf("Expected hasher to be %v, got %v", dh, opt.hasher)
	}
}

func TestWithSHA(t *testing.T) {
	opt := &Option{}
	fn := WithSHA()

	fn(opt)
	_, ok := opt.hasher.(*shaHasher)
	if !ok {
		t.Errorf("Expected hasher to be shaHasher, got %T", opt.hasher)
	}
}

func TestWithFNV(t *testing.T) {
	opt := &Option{}
	fn := WithFNV()

	fn(opt)
	_, ok := opt.hasher.(*fnvHasher)
	if !ok {
		t.Errorf("Expected hasher to be fnvHasher, got %T", opt.hasher)
	}
}

func TestWithExpansion(t *testing.T) {
	opt := &Option{}
	fn := WithExpansion(0.1, 1.5)
	expected := expansion{rate: 1.5, ratio: 0.1}

	fn(opt)
	if opt.expansion != expected {
		t.Errorf("Expected expansion to be %v, got %v", expected, opt.expansion)
	}
}

func TestWithFilledRatio(t *testing.T) {
	opt := &Option{}
	fn := WithFilledRatio(0.123)
	expected := expansion{rate: 2, ratio: 0.123}

	fn(opt)
	if opt.expansion != expected {
		t.Errorf("Expected expansion to be %v, got %v", expected, opt.expansion)
	}
}

func TestWithExpansionRate(t *testing.T) {
	opt := &Option{}
	fn := WithExpansionRate(1.5)
	expected := expansion{rate: 1.5, ratio: 0.5}

	fn(opt)
	if opt.expansion != expected {
		t.Errorf("Expected expansion to be %v, got %v", expected, opt.expansion)
	}
}

func TestWithoutExpansion(t *testing.T) {
	opt := &Option{}
	fn := WithoutExpansion()
	expected := expansion{rate: 0, ratio: 0}

	fn(opt)
	if opt.expansion != expected {
		t.Errorf("Expected expansion to be %v, got %v", expected, opt.expansion)
	}
}

func TestWithNoExpansion(t *testing.T) {
	opt := &Option{}
	fn := WithNoExpansion()
	expected := expansion{rate: 0, ratio: 0}

	fn(opt)
	if opt.expansion != expected {
		t.Errorf("Expected expansion to be %v, got %v", expected, opt.expansion)
	}
}
