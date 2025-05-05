package dice

import (
	"fmt"
	"testing"
)

// TestDC tests whether the roll of one dice passes the check for another dice
func TestDc(t *testing.T) {
	dc := NewDifficultyClass(17)
	d := NewDice(1, 20)

	numRolls := 20
	var totalRollValues, numCriticalHits, numCriticalMisses int
	for range numRolls {
		roll := d.Roll()
		switch {
		case roll.IsCriticalHit():
			fmt.Printf("Rolled: %d (Critical!), Pass: %v\n", roll.Value(), dc.Check(roll))
			numCriticalHits++
		case roll.IsCriticalMiss():
			fmt.Printf("Rolled: %d (Miss!), Pass: %v\n", roll.Value(), dc.Check(roll))
			numCriticalMisses++
		default:
			fmt.Printf("Rolled: %d, Pass: %v\n", roll.Value(), dc.Check(roll))
		}
		totalRollValues += roll.Value()
	}
	fmt.Printf("DC=%v\n", dc)
	fmt.Printf("Dice=%v\n", d)
	fmt.Printf("Average Roll: %d, Critical Hit: %d, Critical Miss: %d\n", totalRollValues/numRolls, numCriticalHits, numCriticalMisses)
}
