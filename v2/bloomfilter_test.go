package bf

import (
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
