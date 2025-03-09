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
