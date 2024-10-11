package bf

import (
	"testing"
)

func TestNew(t *testing.T) {

}

type dummyStorageFactory struct{}

func (d *dummyStorageFactory) Make(capacity uint32) (Storage, error) {
	return nil, nil
}

func TestWithStorage(t *testing.T) {
	opt := &Option{}
	ds := &dummyStorageFactory{}
	fn := WithStorage(ds)

	fn(opt)
	if opt.storageFactory != ds {
		t.Errorf("Expected storage factory to be %v, got %v", ds, opt.storageFactory)
	}
}

type dummyHashFactory struct{}

func (d *dummyHashFactory) Make(numberOfHashFunctions, hashSizeInBits byte) Hash {
	return nil
}

func TestWithHash(t *testing.T) {
	opt := &Option{}
	dh := &dummyHashFactory{}
	fn := WithHash(dh)

	fn(opt)
	if opt.hashFactory != dh {
		t.Errorf("Expected hash factory to be %v, got %v", dh, opt.hashFactory)
	}
}
