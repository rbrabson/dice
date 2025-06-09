package dice

// Value is the value returned by a die roll or difficulty class.
type Value interface {
	Value() int           // Returns the value of the variable
	Check(v Value) bool   // Determines if the roll has a value equal-to-or-greater than the value passed in for the check.
	IsCriticalHit() bool  // Returns true if the value is a critical hit
	IsCriticalMiss() bool // Returns true if the value is a critical miss
}
