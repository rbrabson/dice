package mathx

import "testing"

func TestAbsInt(t *testing.T) {
	intTests := []struct {
		input    int
		expected int
	}{
		{-5, 5},
		{0, 0},
		{5, 5},
	}

	for _, test := range intTests {
		result := Abs(test.input)
		if result != test.expected {
			t.Errorf("abs(%d) = %d; expected %d", test.input, result, test.expected)
		}
	}
}

func TestAbsFloat(t *testing.T) {
	intTests := []struct {
		input    float64
		expected float64
	}{
		{-5.1, 5.1},
		{0.0, 0.0},
		{5.1, 5.1},
	}

	for _, test := range intTests {
		result := Abs(test.input)
		if result != test.expected {
			t.Errorf("abs(%f) = %f; expected %f", test.input, result, test.expected)
		}
	}
}
