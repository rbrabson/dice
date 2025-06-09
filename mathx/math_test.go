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

func TestAbsDiffInt(t *testing.T) {
	tests := []struct {
		x        int
		y        int
		expected int
	}{
		{5, 3, 2},
		{3, 5, 2},
		{-5, -3, 2},
		{-3, -5, 2},
		{5, -3, 8},
		{-3, 5, 8},
		{0, 0, 0},
	}

	for _, test := range tests {
		result := AbsDiff(test.x, test.y)
		if result != test.expected {
			t.Errorf("AbsDiff(%d, %d) = %d; expected %d", test.x, test.y, result, test.expected)
		}
	}
}

func TestAbsDiffFloat(t *testing.T) {
	tests := []struct {
		x        float64
		y        float64
		expected float64
	}{
		{5.5, 3.3, 2.2},
		{3.3, 5.5, 2.2},
		{-5.5, -3.3, 2.2},
		{-3.3, -5.5, 2.2},
		{5.5, -3.3, 8.8},
		{-3.3, 5.5, 8.8},
		{0.0, 0.0, 0.0},
	}

	for _, test := range tests {
		result := AbsDiff(test.x, test.y)
		if result != test.expected {
			t.Errorf("AbsDiff(%f, %f) = %f; expected %f", test.x, test.y, result, test.expected)
		}
	}
}
