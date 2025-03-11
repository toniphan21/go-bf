package bf

import (
	"errors"
	"slices"
	"testing"
)

func TestNewBloomFilter(t *testing.T) {
	cases := []struct {
		name                     string
		config                   Config
		hasher                   Hasher
		storage                  Storage
		expansion                expansion
		expectedMaxBitSetByBlock int
		expectedConfigBlock      ConfigBlock
	}{
		{
			name:                     "maxBitSetByBlock should be calculated",
			config:                   &dummyConfig{k: 5, capacity: 1000, s: 32},
			hasher:                   &dummyHasher{},
			storage:                  &mockStorage{newBlock: map[int]StorageBlock{0: &mockStorageBlock{}}},
			expansion:                expansion{rate: 2, ratio: 0.5},
			expectedMaxBitSetByBlock: 500,
			expectedConfigBlock:      ConfigBlock{Capacity: 1000, NumberOfHashFunctions: 5, KeySizeInBits: 32},
		},
		{
			name:                     "maxBitSetByBlock should be trimmed",
			config:                   &dummyConfig{k: 5, capacity: 1000, s: 16},
			hasher:                   &dummyHasher{},
			storage:                  &mockStorage{newBlock: map[int]StorageBlock{0: &mockStorageBlock{}}},
			expansion:                expansion{rate: 2, ratio: 0.33333333},
			expectedMaxBitSetByBlock: 333,
			expectedConfigBlock:      ConfigBlock{Capacity: 1000, NumberOfHashFunctions: 5, KeySizeInBits: 16},
		},
		{
			name:                     "maxBitSetByBlock should not be rounded, just use cast to int",
			config:                   &dummyConfig{k: 5, capacity: 1000, s: 16},
			hasher:                   &dummyHasher{},
			storage:                  &mockStorage{newBlock: map[int]StorageBlock{0: &mockStorageBlock{}}},
			expansion:                expansion{rate: 2, ratio: 0.6666666666},
			expectedMaxBitSetByBlock: 666,
			expectedConfigBlock:      ConfigBlock{Capacity: 1000, NumberOfHashFunctions: 5, KeySizeInBits: 16},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := newBloomFilter(tc.config, tc.hasher, tc.storage, tc.expansion)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if result.count != 0 {
				t.Errorf("expected 0, got %v", result.count)
			}

			if len(result.bitSetByBlocks) != 1 {
				t.Errorf("expected 1, got %v", result.bitSetByBlocks)
			}
			if result.bitSetByBlocks[0] != 0 {
				t.Errorf("expected 0, got %v", result.bitSetByBlocks[0])
			}

			if result.storage.BlockCount() != 1 {
				t.Errorf("expected 1, got %v", result.storage.BlockCount())
			}

			if len(result.maxBitSetByBlock) != 1 {
				t.Errorf("expected 1, got %v", result.maxBitSetByBlock)
			}
			if result.maxBitSetByBlock[0] != tc.expectedMaxBitSetByBlock {
				t.Errorf("expected %v, got %v", tc.expectedMaxBitSetByBlock, result.maxBitSetByBlock)
			}

			if len(result.configBlocks) != 1 {
				t.Errorf("expected 1, got %v", result.configBlocks)
			}
			if result.configBlocks[0] != tc.expectedConfigBlock {
				t.Errorf("expected %v, got %v", tc.expectedConfigBlock, result.configBlocks[0])
			}
		})
	}
}

func TestBloomFilter_Add(t *testing.T) {
	t.Run("Add puts the key into the last block", func(t *testing.T) {
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mStorageBlock := &mockStorageBlock{capacity: 10}
		mStorage := &mockStorage{blocks: []StorageBlock{mStorageBlock}}
		mHasher := &mockHasher{
			hash: [][]Key{
				{11, 22, 33},
				{44, 55, 66, 77},
			},
		}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if bf.Count() != 0 {
			t.Errorf("expected Count() is 0 before adding, got %v", bf.Count())
		}

		if bf.bitSetByBlocks[0] != 0 {
			t.Errorf("expected bit set by blocks at 0 is 0 before adding, got %v", bf.bitSetByBlocks[0])
		}

		bf.Add([]byte{10, 20, 30})
		if !isArrayEquals(mHasher.hashInput, []byte{10, 20, 30}) {
			t.Errorf("expected hashed called with %v, got %v", []byte{10, 20, 30}, mHasher.hashInput)
		}

		mStorageBlock.assertSetCalledWith(t, []uint32{1, 2, 3})

		if bf.bitSetByBlocks[0] != 3 {
			t.Errorf("expected bit set by blocks at 0 is 3 after adding, got %v", bf.bitSetByBlocks[0])
		}
		if bf.Count() != 1 {
			t.Errorf("expected Count() is 1 after adding, got %v", bf.Count())
		}
	})

	t.Run("Add ignores if the key already added", func(t *testing.T) {
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mStorageBlock := &mockStorageBlock{
			capacity: 10,
			getData:  map[uint32]bool{0: false, 1: true, 2: true, 3: true},
		}
		mStorage := &mockStorage{blocks: []StorageBlock{mStorageBlock}}
		mHasher := &mockHasher{
			hash: [][]Key{
				{11, 22, 33},
				{44, 55, 66, 77},
			},
		}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		bf.Add([]byte{10, 20, 30})
		if !isArrayEquals(mHasher.hashInput, []byte{10, 20, 30}) {
			t.Errorf("expected hashed called with %v, got %v", []byte{10, 20, 30}, mHasher.hashInput)
		}

		if bf.Count() != 0 {
			t.Errorf("expected Count() is 0 after adding because keys already exists, got %v", bf.Count())
		}
	})

	t.Run("Add ignores if the key already added into previous block", func(t *testing.T) {
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mStorageBlock := &mockStorageBlock{
			capacity: 10,
			getData:  map[uint32]bool{0: false, 1: true, 2: true, 3: true},
		}
		mStorage := &mockStorage{blocks: []StorageBlock{mStorageBlock}}
		mHasher := &mockHasher{
			hash: [][]Key{
				{11, 22, 33},
				{44, 55, 66, 77},
			},
		}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		bf.configBlocks = []ConfigBlock{toConfigBlock(dConfig), toConfigBlock(dConfig)}
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		bf.Add([]byte{10, 20, 30})
		if !isArrayEquals(mHasher.hashInput, []byte{10, 20, 30}) {
			t.Errorf("expected hashed called with %v, got %v", []byte{10, 20, 30}, mHasher.hashInput)
		}

		if bf.Count() != 0 {
			t.Errorf("expected Count() is 0 after adding because keys already exists, got %v", bf.Count())
		}
	})

	t.Run("Add will not expand if bitSetByBlocks not exceeded precalculated value", func(t *testing.T) {
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mStorageBlock := &mockStorageBlock{capacity: 10}
		mStorage := &mockStorage{
			blocks: []StorageBlock{mStorageBlock},
		}
		mHasher := &mockHasher{
			hash: [][]Key{
				{11, 22, 33},
			},
		}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		mStorage.newBlock = map[int]StorageBlock{1: &mockStorageBlock{}}
		bf.maxBitSetByBlock[0] = 100

		bf.Add([]byte{10, 20, 30})
		if !isArrayEquals(mHasher.hashInput, []byte{10, 20, 30}) {
			t.Errorf("expected hashed called with %v, got %v", []byte{10, 20, 30}, mHasher.hashInput)
		}

		if len(mStorage.newBlockCalls) > 0 {
			t.Errorf("expected no new block to be added, got %v", len(mStorage.newBlockCalls))
		}
	})

	t.Run("Add will not expand if bitSetByBlocks exceeded precalculated value but expansion rate is 0", func(t *testing.T) {
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mStorageBlock := &mockStorageBlock{capacity: 10}
		mStorage := &mockStorage{
			blocks: []StorageBlock{mStorageBlock},
		}
		mHasher := &mockHasher{
			hash: [][]Key{
				{11, 22, 33},
			},
		}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 0, ratio: 0.6666666666})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		mStorage.newBlock = map[int]StorageBlock{1: &mockStorageBlock{}}
		bf.maxBitSetByBlock[0] = 1

		bf.Add([]byte{10, 20, 30})
		if !isArrayEquals(mHasher.hashInput, []byte{10, 20, 30}) {
			t.Errorf("expected hashed called with %v, got %v", []byte{10, 20, 30}, mHasher.hashInput)
		}

		if len(mStorage.newBlockCalls) > 0 {
			t.Errorf("expected no new block to be added, got %v", len(mStorage.newBlockCalls))
		}
	})

	t.Run("Add will not expand if bitSetByBlocks exceeded precalculated value", func(t *testing.T) {
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mStorageBlock := &mockStorageBlock{capacity: 10}
		mStorage := &mockStorage{blocks: []StorageBlock{mStorageBlock}}
		mHasher := &mockHasher{
			hash: [][]Key{
				{11, 22, 33},
			},
		}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		mStorage.newBlock = map[int]StorageBlock{1: &mockStorageBlock{}}
		bf.maxBitSetByBlock[0] = 1

		bf.Add([]byte{10, 20, 30})
		if !isArrayEquals(mHasher.hashInput, []byte{10, 20, 30}) {
			t.Errorf("expected hashed called with %v, got %v", []byte{10, 20, 30}, mHasher.hashInput)
		}

		if len(mStorage.newBlockCalls) == 0 {
			t.Errorf("expected new block to be added, got %v", len(mStorage.newBlockCalls))
		}

		if len(bf.maxBitSetByBlock) == 1 {
			t.Errorf("expect new maxBitSetByBlock added")
		}

		if len(bf.bitSetByBlocks) == 1 {
			t.Errorf("expect new bitSetByBlocks added")
		}

		if len(bf.configBlocks) == 1 {
			t.Errorf("expect new configBlock added")
		}
	})
}

func TestBloomFilter_Exists(t *testing.T) {
	cases := []struct {
		name     string
		blocks   []StorageBlock
		hash     [][]Key
		expected bool
	}{
		{
			name: "not exists at all - 1 block",
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 10},
			},
			hash:     [][]Key{{11, 22, 33}},
			expected: false,
		},

		{
			name: "not exists at all - 2 blocks",
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 10},
				&mockStorageBlock{capacity: 10},
			},
			hash:     [][]Key{{11, 22, 33}, {44, 55, 66, 77}},
			expected: false,
		},

		{
			name: "not exists - match 2 out of 3 - 1 block",
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 10, getData: map[uint32]bool{0: false, 1: true, 2: true, 3: false}},
			},
			hash:     [][]Key{{11, 22, 33}},
			expected: false,
		},

		{
			name: "exists - 1 block",
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 10, getData: map[uint32]bool{0: false, 1: true, 2: true, 3: true}},
			},
			hash:     [][]Key{{11, 22, 33}},
			expected: true,
		},

		{
			name: "exists - 2 blocks - match at first",
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 10, getData: map[uint32]bool{0: false, 1: true, 2: true, 3: true}},
				&mockStorageBlock{capacity: 10},
			},
			hash:     [][]Key{{11, 22, 33}, {14, 22, 33}},
			expected: true,
		},

		{
			name: "exists - 2 blocks - match at second",
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 10},
				&mockStorageBlock{capacity: 10, getData: map[uint32]bool{0: false, 1: true, 2: true, 3: true}},
			},
			hash:     [][]Key{{44, 55, 66, 77}, {11, 22, 33}},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
			mStorage := &mockStorage{blocks: tc.blocks}
			mHasher := &mockHasher{hash: tc.hash}

			bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			result := bf.Exists([]byte{10, 20, 30})
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestBloomFilter_Storage(t *testing.T) {
	t.Run("it simply returns storage", func(t *testing.T) {
		mStorage := &mockStorage{}
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mHasher := &mockHasher{}
		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		r := bf.Storage()
		if r != mStorage {
			t.Errorf("expected storage to be the same %v, got %v", mStorage, r)
		}
	})
}

func TestBloomFilter_Hasher(t *testing.T) {
	t.Run("it simply returns hasher", func(t *testing.T) {
		mStorage := &mockStorage{}
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mHasher := &mockHasher{}
		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		r := bf.Hasher()
		if r != mHasher {
			t.Errorf("expected hasher to be the same %v, got %v", mStorage, r)
		}
	})
}

func TestBloomFilter_Intersect(t *testing.T) {
	t.Run("returns ErrNilBloomFilter if the given BloomFilter is nil", func(t *testing.T) {
		a := bloomFilter{}
		err := a.Intersect(nil)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !errors.Is(err, ErrNilBloomFilter) {
			t.Errorf("expected ErrNilBloomFilter, got %v", err)
		}
	})

	t.Run("returns ErrStorageDifference if the given BloomFilter has different storage", func(t *testing.T) {
		a := &bloomFilter{storage: &memoryStorage{blocks: make([]StorageBlock, 1)}}
		b := &bloomFilter{storage: &memoryStorage{blocks: make([]StorageBlock, 2)}}
		err := a.Intersect(b)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !errors.Is(err, ErrStorageDifference) {
			t.Errorf("expected ErrStorageDifference, got %v", err)
		}
	})

	t.Run("returns ErrHasherDifference if the given BloomFilter has different hasher", func(t *testing.T) {
		a := &bloomFilter{storage: &memoryStorage{}, hasher: &shaHasher{}}
		b := &bloomFilter{storage: &memoryStorage{}, hasher: &fnvHasher{}}
		err := a.Intersect(b)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !errors.Is(err, ErrHasherDifference) {
			t.Errorf("expected ErrHasherDifference, got %v", err)
		}
	})

	t.Run("uses Clear to change data of current instance for each storage block", func(t *testing.T) {
		storage := &mockStorage{
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: false, 1: false, 2: true, 3: true, 4: true}},
				&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: true, 1: false, 2: false, 3: true, 4: true}},
			},
		}

		a := bloomFilter{
			storage: storage,
			hasher:  &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
		}

		b := bloomFilter{
			storage: &mockStorage{
				blocks: []StorageBlock{
					&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: false, 1: false, 2: true, 3: false, 4: true}},
					&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: false, 1: false, 2: true, 3: true, 4: true}},
				},
			},
			hasher: &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
		}
		err := a.Intersect(&b)

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if a.count != -1 {
			t.Errorf("expected -1, got %v", a.count)
		}

		msb0, ok := storage.blocks[0].(*mockStorageBlock)
		if !ok {
			t.Errorf("expected mockStorageBlock, got %v", storage.blocks[0])
		}
		msb0.assertSetCalledWith(t, []uint32{})
		msb0.assertClearCalledWith(t, []uint32{0, 1, 3})

		msb1, ok := storage.blocks[1].(*mockStorageBlock)
		if !ok {
			t.Errorf("expected mockStorageBlock, got %v", storage.blocks[0])
		}
		msb1.assertSetCalledWith(t, []uint32{})
		msb1.assertClearCalledWith(t, []uint32{0, 1, 2})
	})

	t.Run("uses Block Intersect if the block implement BatchIntersect", func(t *testing.T) {
		storage := &mockStorage{
			blocks: []StorageBlock{
				&memoryStorageBlock{data: []uint{0, 2, 0b001100101110}},
				&memoryStorageBlock{data: []uint{0, 3, 0b001100111101}},
			},
		}
		otherStorage := &mockStorage{
			blocks: []StorageBlock{
				&memoryStorageBlock{data: []uint{7, 0, 0b101000101110}},
				&memoryStorageBlock{data: []uint{1, 0, 0b010100111101}},
			},
		}

		a := bloomFilter{
			storage: storage,
			hasher:  &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
		}

		b := bloomFilter{
			storage: otherStorage,
			hasher:  &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
			count:   1000,
		}
		err := a.Intersect(&b)

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if a.count != -1 {
			t.Errorf("expected -1, got %v", a.count)
		}
		if b.count != 1000 {
			t.Errorf("Intersect should not change count of given BloomFilter")
		}

		sb0, ok := storage.blocks[0].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb0.data[0] != 0 || sb0.data[1] != 0 || sb0.data[2] != 0b001000101110 {
			t.Errorf("Intersect should apply AND operator to all bytes")
		}

		sb1, ok := storage.blocks[1].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb1.data[0] != 0 || sb1.data[1] != 0 || sb1.data[2] != 0b000100111101 {
			t.Errorf("Intersect should apply AND operator to all bytes")
		}

		sb0, ok = otherStorage.blocks[0].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb0.data[0] != 7 || sb0.data[1] != 0 || sb0.data[2] != 0b101000101110 {
			t.Errorf("Intersect should not change the given Storage data")
		}

		sb1, ok = otherStorage.blocks[1].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb1.data[0] != 1 || sb1.data[1] != 0 || sb1.data[2] != 0b010100111101 {
			t.Errorf("Intersect should not change the given Storage data")
		}
	})
}

func TestBloomFilter_Union(t *testing.T) {
	t.Run("returns ErrNilBloomFilter if the given BloomFilter is nil", func(t *testing.T) {
		a := bloomFilter{}
		err := a.Union(nil)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !errors.Is(err, ErrNilBloomFilter) {
			t.Errorf("expected ErrNilBloomFilter, got %v", err)
		}
	})

	t.Run("returns ErrStorageDifference if the given BloomFilter has different storage", func(t *testing.T) {
		a := &bloomFilter{storage: &memoryStorage{blocks: make([]StorageBlock, 1)}}
		b := &bloomFilter{storage: &memoryStorage{blocks: make([]StorageBlock, 2)}}
		err := a.Union(b)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !errors.Is(err, ErrStorageDifference) {
			t.Errorf("expected ErrStorageDifference, got %v", err)
		}
	})

	t.Run("returns ErrHasherDifference if the given BloomFilter has different hasher", func(t *testing.T) {
		a := &bloomFilter{storage: &memoryStorage{}, hasher: &shaHasher{}}
		b := &bloomFilter{storage: &memoryStorage{}, hasher: &fnvHasher{}}
		err := a.Union(b)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !errors.Is(err, ErrHasherDifference) {
			t.Errorf("expected ErrHasherDifference, got %v", err)
		}
	})

	t.Run("uses Set to change data of current instance for each storage block", func(t *testing.T) {
		storage := &mockStorage{
			blocks: []StorageBlock{
				&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: false, 1: false, 2: true, 3: true, 4: true}},
				&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: true, 1: false, 2: false, 3: true, 4: true}},
			},
		}

		a := bloomFilter{
			storage: storage,
			hasher:  &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
		}

		b := bloomFilter{
			storage: &mockStorage{
				blocks: []StorageBlock{
					&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: false, 1: false, 2: true, 3: false, 4: true}},
					&mockStorageBlock{capacity: 5, getData: map[uint32]bool{0: false, 1: false, 2: true, 3: true, 4: true}},
				},
			},
			hasher: &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
		}
		err := a.Union(&b)

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if a.count != -1 {
			t.Errorf("expected -1, got %v", a.count)
		}

		msb0, ok := storage.blocks[0].(*mockStorageBlock)
		if !ok {
			t.Errorf("expected mockStorageBlock, got %v", storage.blocks[0])
		}
		msb0.assertSetCalledWith(t, []uint32{2, 3, 4})
		msb0.assertClearCalledWith(t, []uint32{})

		msb1, ok := storage.blocks[1].(*mockStorageBlock)
		if !ok {
			t.Errorf("expected mockStorageBlock, got %v", storage.blocks[0])
		}
		msb1.assertSetCalledWith(t, []uint32{0, 2, 3, 4})
		msb1.assertClearCalledWith(t, []uint32{})
	})

	t.Run("uses Block Intersect if the block implement BatchUnion", func(t *testing.T) {
		storage := &mockStorage{
			blocks: []StorageBlock{
				&memoryStorageBlock{data: []uint{0, 2, 0b001100101110}},
				&memoryStorageBlock{data: []uint{0, 3, 0b001100111101}},
			},
		}
		otherStorage := &mockStorage{
			blocks: []StorageBlock{
				&memoryStorageBlock{data: []uint{7, 0, 0b101000101110}},
				&memoryStorageBlock{data: []uint{1, 0, 0b010100111101}},
			},
		}

		a := bloomFilter{
			storage: storage,
			hasher:  &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
		}

		b := bloomFilter{
			storage: otherStorage,
			hasher:  &mockHasher{hash: [][]Key{{1, 2}, {3, 4, 5}}},
			count:   1000,
		}
		err := a.Union(&b)

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if a.count != -1 {
			t.Errorf("expected -1, got %v", a.count)
		}
		if b.count != 1000 {
			t.Errorf("Union should not change count of given BloomFilter")
		}

		sb0, ok := storage.blocks[0].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb0.data[0] != 7 || sb0.data[1] != 2 || sb0.data[2] != 0b101100101110 {
			t.Errorf("Union should apply OR operator to all bytes")
		}

		sb1, ok := storage.blocks[1].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb1.data[0] != 1 || sb1.data[1] != 3 || sb1.data[2] != 0b011100111101 {
			t.Errorf("Union should apply OR operator to all bytes")
		}

		sb0, ok = otherStorage.blocks[0].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb0.data[0] != 7 || sb0.data[1] != 0 || sb0.data[2] != 0b101000101110 {
			t.Errorf("Union should not change the given Storage data")
		}

		sb1, ok = otherStorage.blocks[1].(*memoryStorageBlock)
		if !ok {
			t.Errorf("expected memoryStorageBlock, got %v", storage.blocks[0])
		}
		if sb1.data[0] != 1 || sb1.data[1] != 0 || sb1.data[2] != 0b010100111101 {
			t.Errorf("Union should not change the given Storage data")
		}
	})
}

func TestBloomFilter_Clone(t *testing.T) {
	t.Run("copy config, currentConfig, hasher, expansion, count to new instance", func(t *testing.T) {
		mStorage := &mockStorage{}
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mHasher := &mockHasher{}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		bf.count = 100000
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		cloned, ok := bf.Clone().(*bloomFilter)
		if !ok {
			t.Errorf("expected bloomFilter, got %v", bf)
		}

		if bf.config != cloned.config {
			t.Errorf("expected %v, got %v", cloned.config, bf.config)
		}

		if bf.currentConfig != cloned.currentConfig {
			t.Errorf("expected %v, got %v", cloned.config, bf.config)
		}

		if bf.hasher != cloned.hasher {
			t.Errorf("expected %v, got %v", cloned.hasher, bf.hasher)
		}

		if bf.count != cloned.count {
			t.Errorf("expected %v, got %v", cloned.count, bf.count)
		}

		if bf.expansion != cloned.expansion {
			t.Errorf("expected %v, got %v", cloned.expansion, bf.expansion)
		}

		if &bf.expansion == &cloned.expansion {
			t.Errorf("expect expansion should be copied it values, got same instances")
		}
	})

	t.Run("copy bitSetByBlocks, maxBitSetByBlock, configBlocks to new instance", func(t *testing.T) {
		mStorage := &mockStorage{}
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mHasher := &mockHasher{}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		bf.bitSetByBlocks = []int{1, 2, 3, 4, 5}
		bf.maxBitSetByBlock = []int{11, 12, 13, 14, 15}
		bf.configBlocks = []ConfigBlock{
			{Capacity: 1, NumberOfHashFunctions: 2, KeySizeInBits: 3},
			{Capacity: 10, NumberOfHashFunctions: 20, KeySizeInBits: 30},
		}

		cloned, ok := bf.Clone().(*bloomFilter)
		if !ok {
			t.Errorf("expected bloomFilter, got %v", bf)
		}

		if &bf.bitSetByBlocks == &cloned.bitSetByBlocks {
			t.Errorf("expect bitSetByBlocks should be copied it values, got same instances")
		}

		if !slices.Equal(cloned.bitSetByBlocks, []int{1, 2, 3, 4, 5}) {
			t.Errorf("expect bitSetByBlocks should be copied it values, got different values")
		}

		if &bf.maxBitSetByBlock == &cloned.maxBitSetByBlock {
			t.Errorf("expect maxBitSetByBlock should be copied it values, got same instances")
		}

		if !slices.Equal(cloned.maxBitSetByBlock, []int{11, 12, 13, 14, 15}) {
			t.Errorf("expect maxBitSetByBlock should be copied it values, got different values")
		}

		if &bf.configBlocks == &cloned.configBlocks {
			t.Errorf("expect configBlocks should be copied it values, got same instances")
		}

		if !slices.Equal(cloned.configBlocks, []ConfigBlock{
			{Capacity: 1, NumberOfHashFunctions: 2, KeySizeInBits: 3},
			{Capacity: 10, NumberOfHashFunctions: 20, KeySizeInBits: 30},
		}) {
			t.Errorf("expect configBlocks should be copied it values, got different values")
		}
	})

	t.Run("use storage.Clone to clone a storage", func(t *testing.T) {
		mStorage := &mockStorage{}
		dConfig := &dummyConfig{k: 5, capacity: 1000, s: 16}
		mHasher := &mockHasher{}

		bf, err := newBloomFilter(dConfig, mHasher, mStorage, expansion{rate: 2, ratio: 0.6666666666})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		blockA := &memoryStorageBlock{data: []uint{0, 1, 2, 3, 4}, capacity: 10}
		blockB := &memoryStorageBlock{data: []uint{3, 4, 5, 6, 7}, capacity: 10}
		s := &memoryStorage{blocks: []StorageBlock{blockA, blockB}}
		bf.storage = s

		cloned, ok := bf.Clone().(*bloomFilter)
		if !ok {
			t.Errorf("expected bloomFilter, got %v", bf)
		}

		r, ok := cloned.storage.(*memoryStorage)
		if !ok {
			t.Errorf("Expected a memoryStorage, got %T", cloned.storage)
		}

		ap := &s.blocks
		rp := &r.blocks
		if ap == rp {
			t.Errorf("Expected a different blocks to be cloned to different slice")
		}

		for i := 0; i < len(s.blocks); i++ {
			b, ok := s.blocks[i].(*memoryStorageBlock)
			if !ok {
				t.Errorf("Expected a memoryStorageBlock, got %T", s.blocks[i])
			}
			if !b.equals(r.blocks[i]) {
				t.Errorf("Cloned Storage is not equal to the original")
			}
		}
	})
}
