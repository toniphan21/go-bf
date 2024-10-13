package bf

import "testing"

type mockStorage struct {
	setIndex []uint32
	getData  map[uint32]bool
	capacity uint32
}

func (m *mockStorage) Set(index uint32) {
	m.setIndex = append(m.setIndex, index)
}

func (m *mockStorage) Get(index uint32) bool {
	l := len(m.getData)
	if index >= uint32(l) {
		return false
	}
	return m.getData[index]
}

func (m *mockStorage) Capacity() uint32 {
	return m.capacity
}

func (m *mockStorage) assertSetCalledWith(t *testing.T, indices []uint32) {
	if len(m.setIndex) != len(indices) {
		t.Errorf("Set is not called with %v", indices)
	}
	for i, b := range indices {
		if b != m.setIndex[i] {
			t.Errorf("Set is not called with %v", indices)
		}
	}
}

type mockHash struct {
	hashBytes []byte
	hash      []uint32
}

func (m *mockHash) Hash(bytes []byte) []uint32 {
	m.hashBytes = bytes

	return m.hash
}

func (m *mockHash) assertHashCalledWith(t *testing.T, input []byte) {
	if len(m.hashBytes) != len(input) {
		t.Errorf("Hash is not called with %v", input)
	}
	for i, b := range input {
		if b != m.hashBytes[i] {
			t.Errorf("Hash is not called with %v", input)
		}
	}
}

func TestBloomFilter_Add(t *testing.T) {
	hash := &mockHash{hash: []uint32{11, 3, 55, 77}}
	storage := &mockStorage{capacity: 10}
	f := bloomFilter{hash: hash, storage: storage}
	f.Add([]byte("input"))

	hash.assertHashCalledWith(t, []byte("input"))
	storage.assertSetCalledWith(t, []uint32{1, 3, 5, 7})
	if f.count != 1 {
		t.Errorf("expected count is increased to 1")
	}
}

func TestBloomFilter_Count(t *testing.T) {
	hash := &mockHash{hash: []uint32{11, 3, 55, 77}}
	storage := &mockStorage{capacity: 10}
	f := bloomFilter{hash: hash, storage: storage}
	if f.Count() != 0 {
		t.Errorf("expected count is 0")
	}

	f.Add([]byte("input"))
	if f.Count() != 1 {
		t.Errorf("expected count is increased to 1")
	}

	f.Add([]byte("input"))
	if f.Count() != 2 {
		t.Errorf("expected count is increased to 2")
	}
}

func TestBloomFilter_Exists(t *testing.T) {
	cases := []struct {
		name     string
		hash     []uint32
		data     map[uint32]bool
		expected bool
	}{
		{
			name:     "no cell set",
			hash:     []uint32{1, 2},
			data:     map[uint32]bool{0: false, 1: false, 2: false},
			expected: false,
		},
		{
			name:     "1 cell set",
			hash:     []uint32{1, 2},
			data:     map[uint32]bool{0: false, 1: true, 2: false},
			expected: false,
		},
		{
			name:     "2 cells set",
			hash:     []uint32{1, 2},
			data:     map[uint32]bool{0: false, 1: true, 2: true},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hash := &mockHash{hash: tc.hash}
			storage := &mockStorage{getData: tc.data, capacity: 10}
			f := bloomFilter{hash: hash, storage: storage}

			result := f.Exists([]byte("input"))

			hash.assertHashCalledWith(t, []byte("input"))
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestBloomFilter_Data(t *testing.T) {
	storage := &mockStorage{}
	f := bloomFilter{storage: storage}

	result := f.Data()
	if result != storage {
		t.Errorf("expected %v, got %v", storage, result)
	}
}
