package bf

import (
	"errors"
	"testing"
)

func TestMemoryStorageFactory_Make(t *testing.T) {
	cases := []struct {
		name             string
		capacity         uint32
		expectedErr      error
		expectedSize     uint32
		expectedCapacity uint32
	}{
		{
			name:        "invalid capacity",
			capacity:    0,
			expectedErr: ErrInvalidStorageCapacity,
		},
		{
			name:             "rounded capacity",
			capacity:         1024,
			expectedSize:     1024 / bitsetDataSize,
			expectedCapacity: 1024,
		},
		{
			name:             "need to be rounded capacity",
			capacity:         1020,
			expectedSize:     1024 / bitsetDataSize,
			expectedCapacity: 1020,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := memoryStorageFactory{}
			r, err := f.Make(tc.capacity)

			if tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("got error %v, want %v", err, tc.expectedErr)
				}
				return
			}

			s, ok := r.(*bitset)
			if !ok {
				t.Errorf("got type %T, want bitset", r)
			}

			if len(s.data) != int(tc.expectedSize) {
				t.Errorf("got size %d, want %d", len(s.data), int(tc.expectedSize))
			}
			if s.capacity != tc.expectedCapacity {
				t.Errorf("got capacity %d, want %d", s.capacity, tc.expectedCapacity)
			}
		})
	}
}
