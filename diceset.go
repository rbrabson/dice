package dice

import (
	"strconv"
	"strings"
)

// diceSet is a set of dice that can be rolled together.
// This is useful for rolling multiple dice at once, such as when rolling a set of dice for damage.
type diceSet []Dice

// rollSet is the set of rolls from a diceSet.
type rollSet []Roll

// NewDiceSet creates a new set of dice that can be rolled together.
func NewDiceSet(dice ...Dice) Dice {
	ds := make([]Dice, 0, len(dice))
	ds = append(ds, dice...)
	return diceSet(ds)
}

// IsConstant returns `false` for the dice set, as it is not a constant value.
func (ds diceSet) IsConstant() bool {
	return false
}

// IsDebuff returns `true` if all dice in the set are debuffs, `false` otherwise.
func (ds diceSet) IsDebuff() bool {
	if len(ds) == 0 {
		return false
	}
	for _, d := range ds {
		if !d.IsDebuff() {
			return false
		}
	}
	return true
}

// NumDice returns the number of dice in the first dice in the set.
func (ds diceSet) NumDice() int {
	if len(ds) == 0 {
		return 0
	}
	return ds[0].NumDice()
}

// NumSides returns the number of sides on the first die in the set.
func (ds diceSet) NumSides() int {
	if len(ds) == 0 {
		return 0
	}
	return ds[0].NumSides()
}

// Modifier returns the constant value of the first die in the set.
func (ds diceSet) Modifier() int {
	if len(ds) == 0 {
		return 0
	}
	return ds[0].Modifier()
}

// GetDice returns the set of dice to be rolled.
func (ds diceSet) GetDice() []Dice {
	return ds
}

// GetRoll returns the rolls for the set of dice.
func (ds diceSet) GetRoll() Roll {
	return ds[0].Roll()
}

// Source returns the source of the dice set, which is the source of the first die in the set.
func (ds diceSet) Source() string {
	return ds[0].Source()
}

// Roll rolls the dice set and returns the result. The options are applied only to the first dice that
// is rolled, and can be used to roll with advantage or disadvantage.
func (ds diceSet) Roll(opts ...RollOption) Roll {
	rollSet := newRollSet(ds)
	for i, d := range ds {
		var roll Roll
		if i == 0 {
			// If this is the first dice, then we roll it with the options applied
			roll = d.Roll(opts...)
		} else {
			// Otherwise, roll the dice with no options applied
			roll = d.Roll()
		}
		rollSet = append(rollSet, roll)
	}

	return rollSet
}

// newRollSet creates a new set of rolls for the given diceSet.
func newRollSet(diceSet diceSet) rollSet {
	return make([]Roll, 0, len(diceSet))
}

// Value returns the total value of the rolls in the roll set.
func (rs rollSet) Value() int {
	var value int
	for _, r := range rs {
		value += r.Value()
	}

	// If all dice in the set are debuffs, the entire set is a debuff
	// and the value should be negative
	allDebuffs := true
	for _, r := range rs {
		if !r.GetDice().IsDebuff() {
			allDebuffs = false
			break
		}
	}

	if allDebuffs && len(rs) > 0 && value > 0 {
		value = -value
	}

	return value
}

// Check checks if the roll set meets or exceeds the provided value.
func (rs rollSet) Check(v Value) bool {
	if rs.IsCriticalHit() {
		return true
	}
	if rs.IsCriticalMiss() {
		return false
	}
	return rs.Value() >= v.Value()
}

// IsCriticalHit checks if the roll is a critical success.
// This is true iff the first roll in the set is a critical success.
func (rs rollSet) IsCriticalHit() bool {
	return rs[0].IsCriticalHit()
}

// IsCriticalMiss checks if the roll is a critical failure.
// This is true iff the first roll in the set is a critical failure.
func (rs rollSet) IsCriticalMiss() bool {
	return rs[0].IsCriticalMiss()
}

// IsLucky returns `false` for the dice set, as it is not a lucky roll.
func (ds diceSet) IsLucky() bool {
	return false
}

// GetAllRolls returns all the roll values in the roll set.
func (rs rollSet) GetAllRolls() []Roll {
	return rs
}

// RolledWithDisadvantage checks if the roll was made with disadvantage.
func (rs rollSet) RolledWithDisadvantage() bool {
	return rs[0].RolledWithDisadvantage()
}

// RolledWithAdvantage checks if the roll was made with advantage.
func (rs rollSet) RolledWithAdvantage() bool {
	return rs[0].RolledWithAdvantage()
}

// ReRoll re-rolls the dice in the roll set with the provided options.
func (rs rollSet) ReRoll(opts ...RollOption) Roll {
	d := rs.GetDice()
	return d.Roll(opts...)
}

// GetType returns `ROLL_ONCE` for the roll set, as it is a single roll.
func (rs rollSet) GetType() RollType {
	return RollOnce
}

// GetDice returns the set of dice that were rolled in the roll set.
func (rs rollSet) GetDice() Dice {
	dice := make([]Dice, 0, len(rs))
	for _, r := range rs {
		dice = append(dice, r.GetDice())
	}
	return NewDiceSet(dice...)
}

// String returns a string representation of the dice set.
func (ds diceSet) String() string {
	return ds.Str()
}

// Str returns a string representation of the dice set.
func (ds diceSet) Str() string {
	var sb strings.Builder
	for i, d := range ds {
		if i == 0 {
			if d.IsDebuff() {
				sb.WriteString("-")
			}
		} else {
			if d.IsDebuff() {
				sb.WriteString(" - ")
			} else {
				sb.WriteString(" + ")
			}
		}
		sb.WriteString(d.Str())
	}
	return sb.String()
}

// String returns a string representation of the roll set.
func (rs rollSet) String() string {
	var sb strings.Builder
	sb.WriteString(rs.Str())
	sb.WriteString(" = ")
	sb.WriteString(strconv.Itoa(rs.Value()))

	return sb.String()
}

// Str returns a Str representation of the roll set, but without the final value.
func (rs rollSet) Str() string {
	var sb strings.Builder
	for i, r := range rs {
		if i == 0 {
			if r.GetDice().IsDebuff() {
				sb.WriteString("-")
			}
		} else {
			if r.Value() < 0 {
				sb.WriteString(" - ")
			} else {
				sb.WriteString(" + ")
			}
		}
		sb.WriteString(r.Str())
	}

	return sb.String()
}
