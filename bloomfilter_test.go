package bf

import (
	"errors"
	"testing"
)

type mockStorage struct {
	setIndex   []uint32
	clearIndex []uint32
	getData    map[uint32]bool
	capacity   uint32
}

func (m *mockStorage) Set(index uint32) {
	m.setIndex = append(m.setIndex, index)
}

func (m *mockStorage) Clear(index uint32) {
	m.clearIndex = append(m.clearIndex, index)
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

func (m *mockStorage) Equals(other Storage) bool {
	o, ok := other.(*mockStorage)
	if !ok {
		return false
	}
	return o.capacity == m.capacity
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

func (m *mockStorage) assertClearCalledWith(t *testing.T, indices []uint32) {
	if len(m.clearIndex) != len(indices) {
		t.Errorf("Clear is not called with %v", indices)
	}
	for i, b := range indices {
		if b != m.clearIndex[i] {
			t.Errorf("Clear is not called with %v", indices)
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

func (m *mockHash) Equals(other Hash) bool {
	o, ok := other.(*mockHash)
	if !ok {
		return false
	}
	return isArrayEquals(o.hashBytes, m.hashBytes) && isArrayEquals(o.hash, m.hash)
}

func isArrayEquals[T byte | uint32](a, b []T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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

func TestBloomFilter_Storage(t *testing.T) {
	storage := &mockStorage{}
	f := bloomFilter{storage: storage}

	result := f.Storage()
	if result != storage {
		t.Errorf("expected %v, got %v", storage, result)
	}
}

func TestBloomFilter_Hash(t *testing.T) {
	hash := &mockHash{hash: []uint32{11, 3, 55, 77}}
	f := bloomFilter{hash: hash}

	result := f.Hash()
	if result != hash {
		t.Errorf("expected %v, got %v", hash, result)
	}
}

func TestIntersect_ReturnsErrStorageAreNotEquals(t *testing.T) {
	a := bloomFilter{storage: &mockStorage{capacity: 1}}
	b := bloomFilter{storage: &mockStorage{capacity: 2}}
	err := a.Intersect(&b)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if !errors.Is(err, ErrStorageDifference) {
		t.Errorf("expected ErrStorageDifference, got %v", err)
	}
}

func TestIntersect_ReturnsErrIfNotHaveTheSameHash(t *testing.T) {
	a := bloomFilter{storage: &mockStorage{capacity: 1}, hash: &mockHash{hash: []uint32{1, 2}}}
	b := bloomFilter{storage: &mockStorage{capacity: 1}, hash: &mockHash{hash: []uint32{2, 1}}}
	err := a.Intersect(&b)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if !errors.Is(err, ErrHashDifference) {
		t.Errorf("expected ErrHashDifference, got %v", err)
	}
}

func TestIntersect_UseClearToChangeDataOfCurrentInstance(t *testing.T) {
	ad := map[uint32]bool{0: false, 1: false, 2: true, 3: true, 4: true}
	bd := map[uint32]bool{0: false, 1: true, 2: false, 3: true, 4: true}
	storage := &mockStorage{capacity: 5, getData: ad}

	a := bloomFilter{storage: storage, hash: &mockHash{hash: []uint32{1, 2}}}
	b := bloomFilter{storage: &mockStorage{capacity: 5, getData: bd}, hash: &mockHash{hash: []uint32{1, 2}}}
	err := a.Intersect(&b)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	storage.assertSetCalledWith(t, []uint32{})
	storage.assertClearCalledWith(t, []uint32{0, 1, 2})
}

func TestIntersect_ShouldUseIntersectIfTheStorageIsBatchIntersect(t *testing.T) {
	as := &bitset{data: []byte{0, 2, 0b00110011}}
	bs := &bitset{data: []byte{1, 0, 0b01010101}}

	a := bloomFilter{storage: as, hash: &mockHash{hash: []uint32{1, 2}}}
	b := bloomFilter{storage: bs, hash: &mockHash{hash: []uint32{1, 2}}}
	err := a.Intersect(&b)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if as.data[0] != 0 || as.data[1] != 0 && as.data[2] != 0b00010001 {
		t.Errorf("Intersect should apply AND operator to all bytes")
	}
	if bs.data[0] != 1 || bs.data[1] != 0 && bs.data[2] != 0b01010101 {
		t.Errorf("Intersect should not changed the given Storage data")
	}
}
