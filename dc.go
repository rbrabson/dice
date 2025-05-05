package dice

import (
	"fmt"
	"strconv"
)

// DifficultyClass identifies the value a roll of a dice must meet or exceed
// in order to be successful.
type DifficultyClass interface {
	Value
	fmt.Stringer // String returns a string value for the difficulty class
}

// difficultyClass is an implementation of the DifficultyClass interface.
type difficultyClass int

// NewDifficultyClass creates a new difficulty class with the specified target value.
func NewDifficultyClass(targetValue int) DifficultyClass {
	dc := difficultyClass(targetValue)
	return dc
}

// Value returns the value of the variable to be compared in SkillCheck
func (dc difficultyClass) Value() int {
	return int(dc)
}

// IsCriticalHit is always false, as a difficulty class cannot be a critical hit.
func (dc difficultyClass) IsCriticalHit() bool {
	return false
}

// IsCriticalMiss is always false, as a difficulty class cannot be a critical miss.
func (dc difficultyClass) IsCriticalMiss() bool {
	return false
}

// Check determins if the value provided passes the skill check required by the DifficultyClass
func (dc difficultyClass) Check(v Value) bool {
	if v.IsCriticalHit() {
		return true
	}
	if v.IsCriticalMiss() {
		return false
	}
	return v.Value() >= int(dc)
}

// String returns a string value for the difficulty class
func (dc difficultyClass) String() string {
	return strconv.Itoa(int(dc))
}
