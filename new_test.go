package bf

import (
	"github.com/toniphan21/go-bf/internal"
	"testing"
)

func TestNew(t *testing.T) {

}

func TestBloomFilter_NoFalseNegatives(t *testing.T) {
	n, m := 1_000_000, 5_000_000
	cf := WithCapacity(uint32(m), 10)
	bf, _ := New(cf)
	for i := 0; i < n; i++ {
		item := []byte(internal.RandString(10))
		bf.Add(item)
		after := bf.Exists(item)
		if !after {
			t.Fatalf("Bloom Filter has false negative")
		}
	}
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
