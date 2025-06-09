package dice

import (
	"strconv"
	"testing"
)

// TestDCValue tests the Value method of DifficultyClass
func TestDCValue(t *testing.T) {
	tests := []struct {
		dcValue  int
		expected int
	}{
		{0, 0},
		{10, 10},
		{20, 20},
		{-5, -5},
	}

	for _, test := range tests {
		dc := NewDifficultyClass(test.dcValue)
		if dc.Value() != test.expected {
			t.Errorf("DifficultyClass(%d).Value() = %d; expected %d", test.dcValue, dc.Value(), test.expected)
		}
	}
}

// TestDCCritical tests the IsCriticalHit and IsCriticalMiss methods of DifficultyClass
func TestDCCritical(t *testing.T) {
	dc := NewDifficultyClass(15)

	// DifficultyClass should never be a critical hit or miss
	if dc.IsCriticalHit() {
		t.Errorf("DifficultyClass.IsCriticalHit() = true; expected false")
	}

	if dc.IsCriticalMiss() {
		t.Errorf("DifficultyClass.IsCriticalMiss() = true; expected false")
	}
}

// mockRoll is a test implementation of the Value interface
type mockRoll struct {
	value         int
	isCriticalHit bool
	isCriticalMiss bool
}

// Value returns the value of the roll
func (r mockRoll) Value() int { 
	return r.value 
}

// IsCriticalHit returns whether the roll is a critical hit
func (r mockRoll) IsCriticalHit() bool { 
	return r.isCriticalHit 
}

// IsCriticalMiss returns whether the roll is a critical miss
func (r mockRoll) IsCriticalMiss() bool { 
	return r.isCriticalMiss 
}

// Check determines if the roll has a value equal-to-or-greater than the value passed in
func (r mockRoll) Check(v Value) bool { 
	return r.Value() >= v.Value() 
}

// TestDCCheck tests the Check method of DifficultyClass
func TestDCCheck(t *testing.T) {
	dc := NewDifficultyClass(15)

	tests := []struct {
		roll     mockRoll
		expected bool
	}{
		{mockRoll{10, false, false}, false},  // Below DC
		{mockRoll{15, false, false}, true},   // Equal to DC
		{mockRoll{20, false, false}, true},   // Above DC
		{mockRoll{5, true, false}, true},     // Critical hit always passes
		{mockRoll{25, false, true}, false},   // Critical miss always fails
	}

	for i, test := range tests {
		result := dc.Check(test.roll)
		if result != test.expected {
			t.Errorf("Test %d: DC.Check() = %v; expected %v for roll %+v", i, result, test.expected, test.roll)
		}
	}
}

// TestDCString tests the String method of DifficultyClass
func TestDCString(t *testing.T) {
	tests := []struct {
		dcValue  int
		expected string
	}{
		{0, "0"},
		{10, "10"},
		{20, "20"},
		{-5, "-5"},
	}

	for _, test := range tests {
		dc := NewDifficultyClass(test.dcValue)
		if dc.String() != test.expected {
			t.Errorf("DifficultyClass(%d).String() = %s; expected %s", test.dcValue, dc.String(), test.expected)
		}

		// Also verify that it matches the string conversion of the value
		if dc.String() != strconv.Itoa(test.dcValue) {
			t.Errorf("DifficultyClass(%d).String() = %s; expected %s", test.dcValue, dc.String(), strconv.Itoa(test.dcValue))
		}
	}
}
