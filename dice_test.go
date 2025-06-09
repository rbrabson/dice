package dice

import (
	"strconv"
	"strings"
	"testing"
)

// TestParseDice tests the Dice.ParseDice() function and the Dice.String() function
func TestParseDice(t *testing.T) {
	var d Dice
	var diceStr string

	diceStr = "8d8+4"
	d = ParseDice(diceStr)
	if strings.Compare(d.String(), diceStr) != 0 {
		t.Errorf("Errror parsing: got `%s`, expected `%s`", diceStr, d.String())
	}

	diceStr = "d16"
	d = ParseDice(diceStr)
	if strings.Compare(d.String(), "1d16") != 0 {
		t.Errorf("Errror parsing: got `%s`, expected `%s`", diceStr, d.String())
	}

	diceStr = "1d8-3"
	d = ParseDice(diceStr)
	if strings.Compare(d.String(), diceStr) != 0 {
		t.Errorf("Errror parsing `%s`, value =`%s", diceStr, d.String())
	}

	diceStr = "5"
	d = ParseDice(diceStr)
	if strings.Compare(d.String(), diceStr) != 0 {
		t.Errorf("Errror parsing: got `%s`, expected `%s`", diceStr, d.String())
	}

	diceStr = "-7"
	d = ParseDice(diceStr)
	if strings.Compare(d.String(), diceStr) != 0 {
		t.Errorf("Errror parsing: got `%s`, expected `%s`", diceStr, d.String())
	}

	diceStr = "+4"
	d = ParseDice(diceStr)
	if strings.Compare(d.String(), "4") != 0 {
		t.Errorf("Errror parsing: got `%s`, expected `%s`", "4", d.String())
	}
}

// TestRollD20 tests rolling a D20 dice
func TestRollD20(t *testing.T) {
	// Test D20 constant
	d := D20

	// Roll multiple times to ensure values are within expected range
	for i := 0; i < 100; i++ {
		r := d.Roll(WithCriticalHitAllowed())

		// Value should be between 1 and 20
		if r.Value() < 1 || r.Value() > 20 {
			t.Errorf("D20.Roll().Value() = %d; expected between 1 and 20", r.Value())
		}

		// Critical hit should be true only when value is 20
		if r.IsCriticalHit() != (r.Value() == 20) {
			t.Errorf("D20.Roll().IsCriticalHit() = %v; expected %v for value %d", 
				r.IsCriticalHit(), (r.Value() == 20), r.Value())
		}

		// Critical miss should be true only when value is 1
		if r.IsCriticalMiss() != (r.Value() == 1) {
			t.Errorf("D20.Roll().IsCriticalMiss() = %v; expected %v for value %d", 
				r.IsCriticalMiss(), (r.Value() == 1), r.Value())
		}
	}
}

// TestRollDiceWithModifier tests rolling a dice with a modifier
func TestRollDiceWithModifier(t *testing.T) {
	// Test dice with modifier
	d := NewDice(1, 20, WithModifier(4))

	// Roll multiple times to ensure values are within expected range
	for i := 0; i < 100; i++ {
		r := d.Roll()

		// Value should be between 5 (1+4) and 24 (20+4)
		if r.Value() < 5 || r.Value() > 24 {
			t.Errorf("NewDice(1, 20, WithModifier(4)).Roll().Value() = %d; expected between 5 and 24", r.Value())
		}
	}
}

// TestNewConstant tests creating a constant dice
func TestNewConstant(t *testing.T) {
	// Test constant values
	testCases := []struct {
		value    int
		expected int
	}{
		{0, 0},
		{5, 5},
		{-3, -3},
		{100, 100},
	}

	for _, tc := range testCases {
		d := NewConstant(tc.value)

		// Roll multiple times to ensure the value is always the constant
		for i := 0; i < 10; i++ {
			r := d.Roll()

			// Value should always be the constant
			if r.Value() != tc.expected {
				t.Errorf("NewConstant(%d).Roll().Value() = %d; expected %d", 
					tc.value, r.Value(), tc.expected)
			}

			// A constant should never be a critical hit or miss
			if r.IsCriticalHit() {
				t.Errorf("NewConstant(%d).Roll().IsCriticalHit() = true; expected false", tc.value)
			}

			if r.IsCriticalMiss() {
				t.Errorf("NewConstant(%d).Roll().IsCriticalMiss() = true; expected false", tc.value)
			}
		}

		// Test string representation
		expectedStr := strconv.Itoa(tc.value)
		if d.String() != expectedStr {
			t.Errorf("NewConstant(%d).String() = %s; expected %s", 
				tc.value, d.String(), expectedStr)
		}
	}
}

// TestDiceOptions tests the various dice options
func TestDiceOptions(t *testing.T) {
	// Test WithSource option
	sourceText := "Magic Weapon"
	d := NewDice(1, 6, WithSource(sourceText))

	// The source should be included in the string representation
	if !strings.Contains(d.String(), sourceText) {
		t.Errorf("Dice with source %q should include source in String(), got: %s", 
			sourceText, d.String())
	}

	// Test WithModifier option
	modifier := 3
	d = NewDice(1, 8, WithModifier(modifier))

	// Roll multiple times to ensure the modifier is applied
	for i := 0; i < 20; i++ {
		r := d.Roll()

		// Value should be between 1+modifier and 8+modifier
		if r.Value() < 1+modifier || r.Value() > 8+modifier {
			t.Errorf("Dice with modifier %d rolled %d; expected between %d and %d", 
				modifier, r.Value(), 1+modifier, 8+modifier)
		}
	}

	// Test AsDebuff option
	d = NewDice(1, 4, AsDebuff())

	// Roll multiple times to ensure the value is negated
	for i := 0; i < 20; i++ {
		r := d.Roll()

		// Value should be between -4 and -1
		if r.Value() > -1 || r.Value() < -4 {
			t.Errorf("Debuff dice rolled %d; expected between -4 and -1", r.Value())
		}
	}

	// Test combining options
	d = NewDice(2, 10, WithModifier(5), WithSource("Flaming Sword"), AsDebuff())

	// The string should contain the source
	if !strings.Contains(d.String(), "Flaming Sword") {
		t.Errorf("Dice with source should include source in String(), got: %s", d.String())
	}

	// Roll multiple times to ensure all options are applied
	for i := 0; i < 20; i++ {
		r := d.Roll()

		// Value should be between -(2+5) and -(20+5)
		if r.Value() > -7 || r.Value() < -25 {
			t.Errorf("Dice with combined options rolled %d; expected between -25 and -7", r.Value())
		}
	}
}

// TestRollOptions tests the various roll options
func TestRollOptions(t *testing.T) {
	d := NewDice(1, 20)

	// Test rolling with advantage
	r := d.Roll(WithAdvantage())

	// Verify the roll was with advantage
	if !r.RolledWithAdvantage() {
		t.Errorf("Expected roll to be with advantage, got: %v", r.RolledWithAdvantage())
	}

	// Verify the roll was not with disadvantage
	if r.RolledWithDisadvantage() {
		t.Errorf("Expected roll to not be with disadvantage, got: %v", r.RolledWithDisadvantage())
	}

	// Verify we can get all rolls (for advantage)
	allRolls := r.GetAllRolls()
	if len(allRolls) != 2 {
		t.Errorf("Expected 2 rolls for advantage, got: %d", len(allRolls))
	}

	// Test rolling with disadvantage
	r = d.Roll(WithDisadvantage())

	// Verify the roll was with disadvantage
	if !r.RolledWithDisadvantage() {
		t.Errorf("Expected roll to be with disadvantage, got: %v", r.RolledWithDisadvantage())
	}

	// Verify the roll was not with advantage
	if r.RolledWithAdvantage() {
		t.Errorf("Expected roll to not be with advantage, got: %v", r.RolledWithAdvantage())
	}

	// Verify we can get all rolls (for disadvantage)
	allRolls = r.GetAllRolls()
	if len(allRolls) != 2 {
		t.Errorf("Expected 2 rolls for disadvantage, got: %d", len(allRolls))
	}

	// Test rolling with both advantage and disadvantage (they should cancel out)
	r = d.Roll(WithAdvantage(), WithDisadvantage())

	// Verify the roll was not with advantage
	if r.RolledWithAdvantage() {
		t.Errorf("Expected roll to not be with advantage when both options are used, got: %v", r.RolledWithAdvantage())
	}

	// Verify the roll was not with disadvantage
	if r.RolledWithDisadvantage() {
		t.Errorf("Expected roll to not be with disadvantage when both options are used, got: %v", r.RolledWithDisadvantage())
	}

	// Verify we can get all rolls (for normal roll)
	allRolls = r.GetAllRolls()
	if len(allRolls) != 1 {
		t.Errorf("Expected 1 roll for normal roll, got: %d", len(allRolls))
	}

	// Test rolling with no options
	r = d.Roll()

	// Verify the roll was not with advantage
	if r.RolledWithAdvantage() {
		t.Errorf("Expected roll to not be with advantage, got: %v", r.RolledWithAdvantage())
	}

	// Verify the roll was not with disadvantage
	if r.RolledWithDisadvantage() {
		t.Errorf("Expected roll to not be with disadvantage, got: %v", r.RolledWithDisadvantage())
	}

	// Verify we can get all rolls (for normal roll)
	allRolls = r.GetAllRolls()
	if len(allRolls) != 1 {
		t.Errorf("Expected 1 roll for normal roll, got: %d", len(allRolls))
	}
}

// TestSkillCheck tests rolling two multi-sided dice against each other to see if
// the first exceeds the value of the second
func TestSkillCheck(t *testing.T) {
	d1 := NewDiceSet(
		NewDice(1, 20),
		NewDice(1, 4, WithSource("Favorable Beginnings")),
	)
	d2 := NewDice(1, 20)
	numChecks := 80

	for i := 0; i < numChecks; i++ {
		r1 := d1.Roll()
		r2 := d2.Roll()
		pass := r1.Check(r2)

		// Verify that the check result matches the expected comparison
		expectedPass := r1.Value() >= r2.Value()
		if pass != expectedPass {
			t.Errorf("Check result mismatch: got %v, expected %v (r1=%d, r2=%d)", 
				pass, expectedPass, r1.Value(), r2.Value())
		}

		// Also verify that the check is symmetric (if r1 > r2, then r2 should not be > r1)
		if r1.Value() != r2.Value() {
			oppositeCheck := r2.Check(r1)
			if pass == oppositeCheck {
				t.Errorf("Symmetric check failed: r1.Check(r2)=%v, r2.Check(r1)=%v (r1=%d, r2=%d)",
					pass, oppositeCheck, r1.Value(), r2.Value())
			}
		}
	}
}

// TestRollType tests various ways of rolling with advantage or disadvantage, and
// makes sure the Roll identifies the proper roll type.
func TestRollType(t *testing.T) {
	var roll Roll
	d1 := NewDice(1, 20, WithModifier(30))

	roll = d1.Roll()
	if roll.RolledWithAdvantage() == true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}
	if roll.RolledWithDisadvantage() == true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}

	d1 = NewDice(1, 20, WithModifier(30))
	roll = d1.Roll(WithAdvantage())
	if roll.RolledWithAdvantage() != true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}
	if roll.RolledWithDisadvantage() == true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}

	roll = d1.Roll(WithAdvantage(), WithDisadvantage())
	if roll.RolledWithAdvantage() == true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}
	if roll.RolledWithDisadvantage() == true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}

	roll = d1.Roll(WithDisadvantage(), WithAdvantage())
	if roll.RolledWithAdvantage() == true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}
	if roll.RolledWithDisadvantage() == true {
		t.Errorf("Errror rolling with advantage, roll=%s", roll.String())
	}
}

// TestDebuff tests using dice for debuffs
func TestDebuff(t *testing.T) {
	d := ParseDice("1d8", AsDebuff())

	// Verify the dice string representation includes the debuff indicator
	if !strings.Contains(d.String(), "-") {
		t.Errorf("Debuff dice string should contain '-', got: %s", d.String())
	}

	// Test rolling with advantage
	r := d.Roll(WithAdvantage())

	// Verify the roll was with advantage
	if !r.RolledWithAdvantage() {
		t.Errorf("Expected roll to be with advantage, got: %v", r.RolledWithAdvantage())
	}

	// Verify the roll value is within the expected range
	if r.Value() > -1 || r.Value() < -8 {
		t.Errorf("Debuff roll value should be between -8 and -1, got: %d", r.Value())
	}

	// Verify we can get all rolls (for advantage)
	allRolls := r.GetAllRolls()
	if len(allRolls) != 2 {
		t.Errorf("Expected 2 rolls for advantage, got: %d", len(allRolls))
	}

	// Test rolling without advantage
	r = d.Roll()

	// Verify the roll was not with advantage
	if r.RolledWithAdvantage() {
		t.Errorf("Expected roll to not be with advantage, got: %v", r.RolledWithAdvantage())
	}

	// Verify the roll value is within the expected range
	if r.Value() > -1 || r.Value() < -8 {
		t.Errorf("Debuff roll value should be between -8 and -1, got: %d", r.Value())
	}

	// Verify we can get all rolls (for normal roll)
	allRolls = r.GetAllRolls()
	if len(allRolls) != 1 {
		t.Errorf("Expected 1 roll for normal roll, got: %d", len(allRolls))
	}
}

// TestLuckyDice tests the WithLuck option
func TestLuckyDice(t *testing.T) {
	// Create a d6 with the WithLuck option
	d := NewDice(1, 6, WithLuck())

	// Verify the dice is lucky
	if !d.IsLucky() {
		t.Errorf("Expected dice to be lucky, got: %v", d.IsLucky())
	}

	// Roll the dice to ensure it works
	r := d.Roll()
	if r.Value() < 1 || r.Value() > 6 {
		t.Errorf("Lucky dice rolled %d, expected between 1 and 6", r.Value())
	}

	// Compare with a regular dice (not lucky)
	regularDice := NewDice(1, 6)

	// Verify the dice is not lucky
	if regularDice.IsLucky() {
		t.Errorf("Expected regular dice to not be lucky, got: %v", regularDice.IsLucky())
	}

	// Roll the regular dice to ensure it works
	r = regularDice.Roll()
	if r.Value() < 1 || r.Value() > 6 {
		t.Errorf("Regular dice rolled %d, expected between 1 and 6", r.Value())
	}
}

// TestReRoll tests the ReRoll method of Roll
func TestReRoll(t *testing.T) {
	// Create a dice
	d := NewDice(1, 20)

	// Roll the dice
	r1 := d.Roll()

	// Re-roll with advantage
	r2 := r1.ReRoll(WithAdvantage())

	// Verify the re-roll returns a valid roll
	if r2.Value() < 1 || r2.Value() > 20 {
		t.Errorf("Re-roll with advantage returned invalid value: %d", r2.Value())
	}

	// Re-roll with disadvantage
	r3 := r1.ReRoll(WithDisadvantage())

	// Verify the re-roll returns a valid roll
	if r3.Value() < 1 || r3.Value() > 20 {
		t.Errorf("Re-roll with disadvantage returned invalid value: %d", r3.Value())
	}

	// Re-roll with both advantage and disadvantage
	r4 := r1.ReRoll(WithAdvantage(), WithDisadvantage())

	// Verify the re-roll returns a valid roll
	if r4.Value() < 1 || r4.Value() > 20 {
		t.Errorf("Re-roll with both advantage and disadvantage returned invalid value: %d", r4.Value())
	}

	// Test re-rolling a single roll
	// First, get a single roll from a normal roll
	singleRoll := r1.GetAllRolls()[0]

	// Re-roll the single roll with advantage
	r5 := singleRoll.ReRoll(WithAdvantage())

	// Verify the re-roll returns a valid roll
	if r5.Value() < 1 || r5.Value() > 20 {
		t.Errorf("Single roll re-roll returned invalid value: %d", r5.Value())
	}
}

// TestStrMethods tests the Str methods of Dice and Roll
func TestStrMethods(t *testing.T) {
	// Test various dice types
	diceTypes := []Dice{
		NewDice(2, 6),
		NewDice(1, 20, WithModifier(5)),
		NewConstant(10),
		NewDice(3, 8, WithSource("Fireball")),
		NewDice(1, 4, AsDebuff()),
	}

	for _, d := range diceTypes {
		// Test String method of dice
		diceStr := d.String()
		if diceStr == "" {
			t.Errorf("dice.String() returned empty string for %v", d)
		}

		// Test Str method of dice
		diceStrMethod := d.Str()
		if diceStrMethod == "" {
			t.Errorf("dice.Str() returned empty string for %v", d)
		}

		// Roll the dice and test roll string methods
		roll := d.Roll()

		// Test Str method of roll
		rollStr := roll.Str()
		if rollStr == "" {
			t.Errorf("roll.Str() returned empty string for %v", roll)
		}

		// Test String method of roll
		rollString := roll.String()
		if rollString == "" {
			t.Errorf("roll.String() returned empty string for %v", roll)
		}
	}

	// Test roll with advantage and disadvantage
	d := NewDice(1, 20)

	// Roll with advantage
	r := d.Roll(WithAdvantage())
	if r.Str() == "" {
		t.Errorf("roll.Str() with advantage returned empty string")
	}

	// Roll with disadvantage
	r = d.Roll(WithDisadvantage())
	if r.Str() == "" {
		t.Errorf("roll.Str() with disadvantage returned empty string")
	}

	// Test roll with critical hit allowed
	r = d.Roll(WithCriticalHitAllowed())
	if r.Str() == "" {
		t.Errorf("roll.Str() with critical hit allowed returned empty string")
	}
}
