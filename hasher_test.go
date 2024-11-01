package bf

import (
	"bytes"
	"testing"
)

type mockHashFn struct {
	hashCalledCount int
	hashCalledWith  map[int][]byte
	hashReturn      map[int][]byte
}

func (m *mockHashFn) Hash(input *[]byte) []byte {
	m.hashCalledCount++
	if m.hashCalledWith == nil {
		m.hashCalledWith = make(map[int][]byte)
	}
	clone := make([]byte, len(*input))
	for i := 0; i < len(clone); i++ {
		clone[i] = (*input)[i]
	}
	m.hashCalledWith[m.hashCalledCount-1] = clone

	if m.hashReturn != nil {
		r, ok := m.hashReturn[m.hashCalledCount-1]
		if !ok {
			return []byte{}
		}
		return r
	}
	return []byte{}
}

func (m *mockHashFn) assertCalledCount(t *testing.T, expected int) {
	if m.hashCalledCount != expected {
		t.Errorf("expected: %v, got: %v", expected, m.hashCalledCount)
	}
}

func (m *mockHashFn) assertCalledWith(t *testing.T, expected map[int][]byte) {
	for i, _ := range m.hashCalledWith {
		if !bytes.Equal(m.hashCalledWith[i], expected[i]) {
			t.Errorf("expected: %v, got: %v", expected[i], m.hashCalledWith[i])
		}
	}
}

func TestHasher_MakeKeySplitter(t *testing.T) {
	cases := []struct {
		name               string
		hashSizeInBytes    int
		keyCount           byte
		keySize            int
		count              int
		mockedReturn       map[int][]byte
		expectedHashCalled int
		expected           KeySplitter
	}{
		{
			name:            "it calls hashNTimes with n = 1 if keyCount*keySize*count=1 less than hashSizeInBytes",
			hashSizeInBytes: 10,
			keyCount:        5,
			keySize:         7,
			count:           1,
			mockedReturn: map[int][]byte{
				0: {100, 101, 102, 103, 104, 105, 106, 107, 108, 109},
			},
			expectedHashCalled: 1,
			expected: KeySplitter{
				Source:   []byte{100, 101, 102, 103, 104, 105, 106, 107, 108, 109},
				Count:    1,
				KeyCount: 5,
				KeySize:  7,
			},
		},
		{
			name:            "it calls hashNTimes with n = 1 if keyCount*keySize*count=2 less than hashSizeInBytes",
			hashSizeInBytes: 10,
			keyCount:        5,
			keySize:         7,
			count:           2,
			mockedReturn: map[int][]byte{
				0: {100, 101, 102, 103, 104, 105, 106, 107, 108, 109},
			},
			expectedHashCalled: 1,
			expected: KeySplitter{
				Source:   []byte{100, 101, 102, 103, 104, 105, 106, 107, 108, 109},
				Count:    2,
				KeyCount: 5,
				KeySize:  7,
			},
		},
		{
			name:            "it calls hashNTimes with n = 2 if keyCount*keySize*count=1 greater than hashSizeInBytes",
			hashSizeInBytes: 8,
			keyCount:        10,
			keySize:         7,
			count:           1,
			mockedReturn: map[int][]byte{
				0: {100, 101, 102, 103, 104, 105, 106, 107},
				1: {200, 201, 202, 203, 204, 205, 206, 207},
			},
			expectedHashCalled: 2,
			expected: KeySplitter{
				Source:   []byte{100, 101, 102, 103, 104, 105, 106, 107, 200, 201, 202, 203, 204, 205, 206, 207},
				Count:    1,
				KeyCount: 10,
				KeySize:  7,
			},
		},
		{
			name:            "it calls hashNTimes with n = 2 if keyCount*keySize*count=2 greater than hashSizeInBytes",
			hashSizeInBytes: 8,
			keyCount:        5,
			keySize:         7,
			count:           2,
			mockedReturn: map[int][]byte{
				0: {100, 101, 102, 103, 104, 105, 106, 107},
				1: {200, 201, 202, 203, 204, 205, 206, 207},
			},
			expectedHashCalled: 2,
			expected: KeySplitter{
				Source:   []byte{100, 101, 102, 103, 104, 105, 106, 107, 200, 201, 202, 203, 204, 205, 206, 207},
				Count:    2,
				KeyCount: 5,
				KeySize:  7,
			},
		},
		{
			name:            "it calls hashNTimes with n = 3 if keyCount*keySize*count=2 greater than hashSizeInBytes",
			hashSizeInBytes: 8,
			keyCount:        5,
			keySize:         9,
			count:           3,
			mockedReturn: map[int][]byte{
				0: {0, 1, 2, 3, 4, 5, 6, 7},
				1: {100, 101, 102, 103, 104, 105, 106, 107},
				2: {200, 201, 202, 203, 204, 205, 206, 207},
			},
			expectedHashCalled: 3,
			expected: KeySplitter{
				Source:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 100, 101, 102, 103, 104, 105, 106, 107, 200, 201, 202, 203, 204, 205, 206, 207},
				Count:    3,
				KeyCount: 5,
				KeySize:  9,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &hasher{hashSizeInBytes: tc.hashSizeInBytes, keyCount: tc.keyCount, keySize: tc.keySize}
			m := &mockHashFn{
				hashReturn: tc.mockedReturn,
			}
			result := h.makeKeySplitter(tc.count, []byte{1, 2, 3}, m.Hash)

			m.assertCalledCount(t, tc.expectedHashCalled)
			if !bytes.Equal(result.Source, tc.expected.Source) {
				t.Errorf("expected: %v, got: %v", tc.expected.Source, result.Source)
			}
			if result.Count != tc.expected.Count {
				t.Errorf("expected: %v, got: %v", tc.expected.Count, result.Count)
			}
			if result.KeyCount != tc.expected.KeyCount {
				t.Errorf("expected: %v, got: %v", tc.expected.KeyCount, result.KeyCount)
			}
			if result.KeySize != tc.expected.KeySize {
				t.Errorf("expected: %v, got: %v", tc.expected.KeySize, result.KeySize)
			}
		})
	}
}

func TestHasher_HashNTimes(t *testing.T) {
	cases := []struct {
		name               string
		n                  byte
		hashSizeInBytes    int
		input              []byte
		mockedReturn       map[int][]byte
		expectedCount      int
		expectedCalledWith map[int][]byte
		expected           []byte
	}{
		{
			name:            "it should call hashFn() one time with original input if n = 1",
			n:               1,
			hashSizeInBytes: 2,
			input:           []byte{1, 2, 3},
			mockedReturn: map[int][]byte{
				0: {100, 101},
			},
			expectedCount: 1,
			expectedCalledWith: map[int][]byte{
				0: {1, 2, 3},
			},
			expected: []byte{100, 101},
		},
		{
			name:            "it should call hashFn() twice with original input and padded 0 if n = 2",
			n:               2,
			hashSizeInBytes: 2,
			input:           []byte{1, 2, 3},
			mockedReturn: map[int][]byte{
				0: {100, 101},
				1: {102, 103},
			},
			expectedCount: 2,
			expectedCalledWith: map[int][]byte{
				0: {1, 2, 3},
				1: {0, 1, 2, 3},
			},
			expected: []byte{100, 101, 102, 103},
		},
		{
			name:            "it should call hashFn() 2 times with original input and padded 0, 1 if n = 3",
			n:               3,
			hashSizeInBytes: 2,
			input:           []byte{1, 2, 3},
			mockedReturn: map[int][]byte{
				0: {100, 101},
				1: {102, 103},
				2: {104, 105},
			},
			expectedCount: 3,
			expectedCalledWith: map[int][]byte{
				0: {1, 2, 3},
				1: {0, 1, 2, 3},
				2: {1, 1, 2, 3},
			},
			expected: []byte{100, 101, 102, 103, 104, 105},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &hasher{hashSizeInBytes: tc.hashSizeInBytes}
			m := &mockHashFn{
				hashReturn: tc.mockedReturn,
			}
			result := h.hashNTimes(tc.n, &tc.input, m.Hash)

			m.assertCalledCount(t, tc.expectedCount)
			m.assertCalledWith(t, tc.expectedCalledWith)
			if !bytes.Equal(tc.expected, result) {
				t.Errorf("expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}
