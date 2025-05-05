package dice

import (
	"fmt"
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

// TestRollDice tests rolling a multi-sided dice
func TestRollDice(t *testing.T) {
	var d Dice
	var r Roll

	d = D20
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
func TestSkillCheck(t *testing.T) {
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
	fmt.Println(d)
	r := d.Roll(WithAdvantage())
	fmt.Println(r)
	fmt.Println(r.GetAllRolls())
}
