package dice

import (
	"fmt"
	"testing"
)

// TestRollDice tests rolling a multi-sided dice
func TestRollDiceSet(t *testing.T) {
	var d Dice
	var r Roll

	d = NewDiceSet(
		D20,
		ParseDice("1d20+5"),
		NewConstant(-1, WithSource("Strength Modifier")),
		ParseDice("1d4", WithSource("Guidance")),
		ParseDice("1d8", WithSource("Debuff"), AsDebuff()),
	)
	r = d.Roll()
	fmt.Println("Roll:", r, "=> Value:", r.Value())
	fmt.Println("IsCriticalHit:", r.IsCriticalHit())
	fmt.Println("IsCriticalMiss:", r.IsCriticalMiss())

	d = NewDice(1, 20, WithModifier(4))
	r = d.Roll()
	fmt.Println("Roll:", r, "=> Value:", r.Value())
}

// TestSkillCheck tests rolling two multi-sided dice against each other to see if
// the first exceeds the value of the second
func TestDiceSetSkillCheck(t *testing.T) {
	d1 := NewDiceSet(
		NewDice(1, 20),
		NewDice(1, 4, WithSource("Favorable Beginnings")),
	)
	d2 := NewDice(1, 20)
	numChecks := 80

	var numPassed int
	for i := 0; i < numChecks; i++ {
		r1 := d1.Roll()
		r2 := d2.Roll()
		pass := r1.Check(r2)
		fmt.Printf("Pass: %v, %d > %d\n", pass, r1.Value(), r2.Value())
		if pass {
			numPassed++
		}
	}
	fmt.Printf("Passed: %d, Failed: %d\n", numPassed, numChecks-numPassed)
}

// TestRollType tests various ways of rolling with advantage or disadvantage, and
// makes sure the Roll identifies the proper roll type.
func TestDiceSetRollType(t *testing.T) {
	var roll Roll
	d1 := NewDiceSet(
		NewDice(1, 20, WithModifier(30)),
		NewConstant(5, WithSource("Test Constant")),
	)

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
func TestDiceSetDebuff(t *testing.T) {
	d := NewDiceSet(
		NewDice(1, 20, AsDebuff()),
		NewDice(1, 4, WithSource("Guidance"), AsDebuff()),
	)
	fmt.Println(d)
	r := d.Roll(WithAdvantage())
	if r.RolledWithAdvantage() != true {
		t.Errorf("Expected roll to be with advantage, got %s", r.String())
	}
	if r.RolledWithDisadvantage() == true {
		t.Errorf("Expected roll to not be with disadvantage, got %s", r.String())
	}
	if r.Value() >= 0 {
		t.Errorf("Expected roll to be a debuff, got %s", r.String())

	}
	fmt.Println("GetAllRolls", r.GetAllRolls())
}
