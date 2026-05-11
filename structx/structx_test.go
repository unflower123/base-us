package structx

import (
	"errors"
	"reflect"
	"testing"
)

func TestCopyStruct(t *testing.T) {
	tests := []struct {
		name        string
		src         interface{}
		dst         interface{}
		expected    interface{}
		expectedErr error
	}{
		{
			name: "Copy struct to struct pointer",
			src: struct {
				Name string
				Age  int
			}{Name: "Alice", Age: 30},
			dst: &struct {
				Name string
				Age  int
			}{},
			expected: &struct {
				Name string
				Age  int
			}{Name: "Alice", Age: 30},
			expectedErr: nil,
		},
		{
			name: "Copy struct pointer to struct pointer",
			src: &struct {
				Name string
				Age  int
			}{Name: "Bob", Age: 25},
			dst: &struct {
				Name string
				Age  int
			}{},
			expected: &struct {
				Name string
				Age  int
			}{Name: "Bob", Age: 25},
			expectedErr: nil,
		},
		{
			name: "Copy struct with mismatched fields",
			src: struct {
				Name string
				Age  int
			}{Name: "Charlie", Age: 40},
			dst: &struct {
				Name string
				City string
			}{},
			expected: &struct {
				Name string
				City string
			}{Name: "Charlie", City: ""},
			expectedErr: nil,
		},
		{
			name: "Copy struct with unexported fields",
			src: struct {
				Name string
				age  int
			}{Name: "David", age: 50},
			dst: &struct {
				Name string
				age  int
			}{},
			expected: &struct {
				Name string
				age  int
			}{Name: "David", age: 0}, // unexported fields should not be copied
			expectedErr: nil,
		},
		{
			name:        "Copy to nil destination",
			src:         struct{ Name string }{Name: "Eve"},
			dst:         (*struct{ Name string })(nil),
			expected:    nil,
			expectedErr: errors.New("dst must be a non-nil struct pointer"),
		},
		{
			name:        "Copy from non-struct source",
			src:         "not a struct",
			dst:         &struct{ Name string }{},
			expected:    &struct{ Name string }{},
			expectedErr: errors.New("src must be a struct or struct pointer"),
		},
		{
			name:        "Copy to non-struct destination",
			src:         struct{ Name string }{Name: "Frank"},
			dst:         new(int),
			expected:    new(int),
			expectedErr: errors.New("dst must be a non-nil struct pointer"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CopyStruct(tt.src, tt.dst)

			// Check if the error matches the expected error
			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
			}

			// If no error is expected, compare the destination with the expected result
			if err == nil && tt.expectedErr == nil {
				if !reflect.DeepEqual(tt.dst, tt.expected) {
					t.Errorf("Expected: %+v, got: %+v", tt.expected, tt.dst)
				}
			}
		})
	}
}
