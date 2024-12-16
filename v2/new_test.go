package bf

import "testing"

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

func (d *dummyConfig) Next(expansionRate float64) Config {
	return &dummyConfig{
		info:     d.info,
		k:        d.k,
		capacity: uint32(float64(d.capacity) * expansionRate),
	}
}

func TestWithStorage(t *testing.T) {
	opt := &Option{}
	ds := &dummyStorage{}
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
