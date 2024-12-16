package bf

type Storage interface {
	NewBlock(capacity uint32)

	BlockCount() int

	Block(index int) StorageBlock

	IsCompatible(other Storage) bool
}

type memoryStorage struct {
	blocks []StorageBlock
}

func (s *memoryStorage) NewBlock(capacity uint32) {
	if capacity == 0 {
		capacity = DefaultSizeInBits
	}
	n, m := capacity/wordSize, capacity%wordSize
	if m > 0 {
		n += 1
	}
	s.blocks = append(s.blocks, newMemoryStorageBlock(n, capacity))
}

func (s *memoryStorage) BlockCount() int {
	return len(s.blocks)
}

func (s *memoryStorage) Block(index int) StorageBlock {
	return s.blocks[index]
}

func (s *memoryStorage) IsCompatible(other Storage) bool {
	o, ok := other.(*memoryStorage)
	if !ok {
		return false
	}

	if len(s.blocks) != len(o.blocks) {
		return false
	}

	for i := 0; i < len(s.blocks); i++ {
		if !s.blocks[i].IsCompatible(o.blocks[i]) {
			return false
		}
	}

	return true
}

var _ Storage = (*memoryStorage)(nil)
