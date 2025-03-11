package bf

import (
	"testing"
)

type mockStorage struct {
	blocks        []StorageBlock
	newBlock      map[int]StorageBlock
	newBlockCalls map[int]uint32
}

func (m *mockStorage) NewBlock(capacity uint32) {
	if m.newBlock == nil {
		return
	}

	idx := len(m.newBlockCalls)
	m.blocks = append(m.blocks, m.newBlock[idx])

	if m.newBlockCalls == nil {
		m.newBlockCalls = make(map[int]uint32)
	}
	m.newBlockCalls[idx] = capacity
}

func (m *mockStorage) BlockCount() int {
	return len(m.blocks)
}

func (m *mockStorage) Block(index int) StorageBlock {
	return m.blocks[index]
}

func (m *mockStorage) IsCompatible(other Storage) bool {
	o, ok := other.(*mockStorage)
	if !ok {
		return false
	}

	if len(m.blocks) != len(o.blocks) {
		return false
	}

	for i := 0; i < len(m.blocks); i++ {
		if !m.blocks[i].IsCompatible(o.blocks[i]) {
			return false
		}
	}

	return true
}

func (m *mockStorage) Clone() Storage {
	cloned := &mockStorage{}
	if m.blocks != nil {
		cloned.blocks = make([]StorageBlock, len(m.blocks))
		for i, block := range m.blocks {
			cloned.blocks[i] = block.Clone()
		}
	}

	if m.newBlock != nil {
		cloned.newBlock = make(map[int]StorageBlock, len(m.newBlock))
		for k, v := range m.newBlock {
			cloned.newBlock[k] = v.Clone()
		}
	}

	if m.newBlockCalls != nil {
		cloned.newBlockCalls = make(map[int]uint32, len(m.newBlockCalls))
		for k, v := range m.newBlockCalls {
			cloned.newBlockCalls[k] = v
		}
	}

	return cloned
}

func TestMemoryStorage_NewBlock(t *testing.T) {
	cases := []struct {
		name                string
		blocks              []StorageBlock
		capacity            uint32
		expectedN           uint32
		expectedBlocksCount int
	}{
		{
			name:                "if capacity is 0, use default capacity",
			capacity:            0,
			expectedN:           DefaultSizeInBits / wordSize,
			expectedBlocksCount: 1,
		},

		{
			name:                "less than wordSize, no blocks",
			capacity:            wordSize - 1,
			blocks:              make([]StorageBlock, 0),
			expectedBlocksCount: 1,
			expectedN:           1,
		},

		{
			name:                "less than wordSize, 1 block",
			capacity:            wordSize - 1,
			blocks:              make([]StorageBlock, 1),
			expectedBlocksCount: 2,
			expectedN:           1,
		},

		{
			name:                "exact wordSize, no blocks",
			capacity:            wordSize,
			blocks:              make([]StorageBlock, 0),
			expectedBlocksCount: 1,
			expectedN:           1,
		},

		{
			name:                "exact wordSize, 1 block",
			capacity:            wordSize,
			blocks:              make([]StorageBlock, 1),
			expectedBlocksCount: 2,
			expectedN:           1,
		},

		{
			name:                "n * wordSize + 1, no blocks",
			capacity:            100*wordSize + 1,
			blocks:              make([]StorageBlock, 0),
			expectedBlocksCount: 1,
			expectedN:           101,
		},

		{
			name:                "n*wordSize + 1, 1 block",
			capacity:            100*wordSize + 1,
			blocks:              make([]StorageBlock, 1),
			expectedBlocksCount: 2,
			expectedN:           101,
		},

		{
			name:                "exact n*wordSize, no blocks",
			capacity:            100 * wordSize,
			blocks:              make([]StorageBlock, 0),
			expectedBlocksCount: 1,
			expectedN:           100,
		},

		{
			name:                "exact n*wordSize, 1 block",
			capacity:            100 * wordSize,
			blocks:              make([]StorageBlock, 1),
			expectedBlocksCount: 2,
			expectedN:           100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ms := &memoryStorage{}
			ms.blocks = tc.blocks
			ms.NewBlock(tc.capacity)

			if ms.BlockCount() != tc.expectedBlocksCount {
				t.Errorf("expected %v, got %v", tc.expectedBlocksCount, ms.BlockCount())
			}
			last := ms.Block(tc.expectedBlocksCount - 1)
			msb, ok := last.(*memoryStorageBlock)
			if !ok {
				t.Errorf("expected a memoryStorageBlock, got %T", last)
			}

			if len(msb.data) != int(tc.expectedN) {
				t.Errorf("expected %v, got %v", tc.expectedN, len(msb.data))
			}
		})
	}
}

func TestMemoryStorage_IsCompatible(t *testing.T) {
	a := &memoryStorageBlock{data: []uint{0, 1, 2, 3, 4}, capacity: 10}
	b := &memoryStorageBlock{data: []uint{0, 1, 2}, capacity: 10}
	c := &memoryStorageBlock{data: []uint{0, 1, 2}, capacity: 20}

	cases := []struct {
		name     string
		instance Storage
		given    Storage
		expected bool
	}{
		{
			name:     "if given storage is not memoryStorage, not compatible",
			instance: &memoryStorage{},
			given:    &mockStorage{},
			expected: false,
		},

		{
			name:     "if given blocks has different length, not compatible",
			instance: &memoryStorage{blocks: make([]StorageBlock, 1)},
			given:    &memoryStorage{blocks: make([]StorageBlock, 2)},
			expected: false,
		},

		{
			name:     "if given blocks has same length, different capacity, not compatible - 1",
			instance: &memoryStorage{blocks: []StorageBlock{a}},
			given:    &memoryStorage{blocks: []StorageBlock{c}},
			expected: false,
		},

		{
			name:     "if given blocks has same length, different capacity, not compatible - 2",
			instance: &memoryStorage{blocks: []StorageBlock{a, b}},
			given:    &memoryStorage{blocks: []StorageBlock{a, c}},
			expected: false,
		},

		{
			name:     "if given blocks has same length, same capacity, compatible",
			instance: &memoryStorage{blocks: []StorageBlock{a, b, c}},
			given:    &memoryStorage{blocks: []StorageBlock{a, b, c}},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.instance.IsCompatible(tc.given)
			if r != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, r)
			}
		})
	}
}

func TestMemoryStorage_Clone(t *testing.T) {
	blockA := &memoryStorageBlock{data: []uint{0, 1, 2, 3, 4}, capacity: 10}
	blockB := &memoryStorageBlock{data: []uint{3, 4, 5, 6, 7}, capacity: 10}
	s := &memoryStorage{blocks: []StorageBlock{blockA, blockB}}
	r, ok := s.Clone().(*memoryStorage)
	if !ok {
		t.Errorf("Expected a memoryStorage, got %T", s)
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
}
