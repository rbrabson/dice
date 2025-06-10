package slicex

import "testing"

func TestRotateInt(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		num      int
		expected []int
	}{
		{
			name:     "rotate empty slice",
			slice:    []int{},
			num:      3,
			expected: []int{},
		},
		{
			name:     "rotate by 0",
			slice:    []int{1, 2, 3, 4, 5},
			num:      0,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "rotate by positive number less than length",
			slice:    []int{1, 2, 3, 4, 5},
			num:      2,
			expected: []int{4, 5, 1, 2, 3},
		},
		{
			name:     "rotate by positive number equal to length",
			slice:    []int{1, 2, 3, 4, 5},
			num:      5,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "rotate by positive number greater than length",
			slice:    []int{1, 2, 3, 4, 5},
			num:      7,
			expected: []int{4, 5, 1, 2, 3},
		},
		{
			name:     "rotate by negative number",
			slice:    []int{1, 2, 3, 4, 5},
			num:      -1,
			expected: []int{5, 1, 2, 3, 4},
		},
		{
			name:     "rotate by negative number with absolute value greater than length",
			slice:    []int{1, 2, 3, 4, 5},
			num:      -6,
			expected: []int{5, 1, 2, 3, 4},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Rotate(test.slice, test.num)
			if !slicesEqual(result, test.expected) {
				t.Errorf("Rotate(%v, %d) = %v; expected %v", test.slice, test.num, result, test.expected)
			}
		})
	}
}

func TestRotateString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		num      int
		expected []string
	}{
		{
			name:     "rotate string slice",
			slice:    []string{"a", "b", "c", "d", "e"},
			num:      2,
			expected: []string{"d", "e", "a", "b", "c"},
		},
		{
			name:     "rotate string slice by negative number",
			slice:    []string{"a", "b", "c", "d", "e"},
			num:      -1,
			expected: []string{"e", "a", "b", "c", "d"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Rotate(test.slice, test.num)
			if !slicesEqual(result, test.expected) {
				t.Errorf("Rotate(%v, %d) = %v; expected %v", test.slice, test.num, result, test.expected)
			}
		})
	}
}

// Helper function to compare slices
func slicesEqual[T comparable](a, b []T) bool {
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