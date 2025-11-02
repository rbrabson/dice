package dice

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/rbrabson/dice/mathx"
)

// rng is a random number generator used for dice rolls
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Pre-defined dice types
var (
	D4 = &dice{
		numDice:  1,
		numSides: 4,
	}
	D6 = &dice{
		numDice:  1,
		numSides: 6,
	}
	D8 = &dice{
		numDice:  1,
		numSides: 8,
	}
	D10 = &dice{
		numDice:  1,
		numSides: 10,
	}
	D12 = &dice{
		numDice:  1,
		numSides: 12,
	}
	D20 = &dice{
		numDice:  1,
		numSides: 20,
	}
	D100 = &dice{
		numDice:  1,
		numSides: 100,
	}
)

type RollType int

const (
	_                    RollType = iota
	RollOnce                      // Single dice is rolled
	RollWithAdvantage             // Two dice are rolled, the highest value is used
	RollWithDisadvantage          // Two dice are rolled, the lowest value is used
)

const (
	CriticalHit  = 20 // The value for a critical hit; this is the highest value that can be rolled on a D20
	CriticalMiss = 1  // The value for a critical miss; this is the lowest value that can be rolled on a D20
)

// Dice is a multi-sided dice that can be rolled for a value, both for a base value or with
// advantage or disadvantage. The multi-sided dice also allows the source for the dice to be
// included in the output.
type Dice interface {
	GetDice() []Dice         // Returns all modifiers for the dice
	GetRoll() Roll           // Gets the last roll of the dice
	IsConstant() bool        // Returns true if the dice is a constant value
	IsDebuff() bool          // Returns true if the dice is a debuff
	IsLucky() bool           // Returns true if the dice is a lucky dice
	NumDice() int            // The number of dice to roll
	NumSides() int           // The number of sides on the dice
	Modifier() int           // The constant modifier to add to the roll
	Source() string          // Get the source for the dice
	Roll(...RollOption) Roll // Rolls a dice and returns the result
	fmt.Stringer             // String representation of the dice
	Str() string             // Returns a string representation of the dice, but without the leading '-'
}

// Roll is a roll of a Dice
type Roll interface {
	Value
	GetAllRolls() []Roll          // Gets the rolls for the dice; there may be two rolls if rolling with advantage/disadvantage
	RolledWithDisadvantage() bool // Returns true if the roll was made with disadvantage
	RolledWithAdvantage() bool    // Returns true if the roll was made with advantage
	ReRoll(...RollOption) Roll    // Re-rolls the dice with the provided options, returning a new Roll
	GetType() RollType            // Gets the type of roll (ROLL_ONCE, ROLL_WITH_ADVANTAGE, ROLL_WITH_DISADVANTATE)
	GetDice() Dice                // The dice used for the roll
	fmt.Stringer                  // Get a string representation of a roll
	Str() string                  // Returns a string representation of the roll, but without the final value
}

// dice is an implementation of the Dice interface.
type dice struct {
	numDice  int    // The number of dice to roll
	numSides int    // The number of sides on the dice
	modifier int    // A constant value to add to the roll
	roll     Roll   // The roll of the dice
	source   string // Source for the dice; used in creating the descripton output
	isLucky  bool   // If true, the dice is a lucky dice that is re-rolled if it rolls a 1
	isDebuff bool   // The dice roll is negated
}

// DiceOption is a function that modifies the default values of a dice.
type DiceOption func(*dice)

// roll is an implementation of the Roll interface
type roll struct {
	rollType           RollType      // Type of roll (ROLL_ONCE, ROLL_ADVANTAGE, ROLL_DISADVANTATE)
	rolls              []*singleRoll // The values rolled for the dice, if rolled with advantage or disadvantage
	value              int           // The value of the roll
	criticalHitAllowed bool          // If true, the dice allows for a critical hit
	criticalHit        int           // The value for a critical hit; defaults to 20
	criticalMiss       int           // The value for a critical miss; defaults to 1
	dice               *dice         // The dice used for the roll
}

// singleRoll represents a single roll of the dice. Whenn rolling a dice, there may be one roll or,
// if rolling with advantage or disadvantage, two rolls.
type singleRoll struct {
	value              int   // The value of the roll
	criticalHitAllowed bool  // If true, the roll allows for a critical hit
	criticalHit        int   // The value for a critical hit; defaults to 20
	criticalMiss       int   // The value for a critical miss; defaults to 1
	dice               *dice // The dice used for the roll
}

// RollOption is a function that can modify the default values of a roll.
type RollOption func(*roll)

// NewDice returns a multi-sided dice with an optional base damage that may be rolled for a value.
func NewDice(numDice int, numSides int, opts ...DiceOption) Dice {
	d := &dice{
		numDice:  numDice,
		numSides: numSides,
		roll:     nil,
	}
	for _, opt := range opts {
		opt(d)
	}

	return d
}

// NewConstant creates a new dice that always returns the provided value.
func NewConstant(value int, opts ...DiceOption) Dice {
	newOpts := make([]DiceOption, 0, len(opts)+1)
	newOpts = append(newOpts, opts...)
	newOpts = append(newOpts, WithModifier(value)) // Add the constant value as a modifier
	return NewDice(0, 0, newOpts...)
}

// ParseDice parses a string representation of a dice into a Dice. Some supported formats are:
// `1d20` `1d20+5`, `1d8-2`, and `d4`.
func ParseDice(str string, opts ...DiceOption) Dice {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	//  Modifiers to apply to the dice
	modifiers := make([]DiceOption, 0, len(opts)+2)
	modifiers = append(modifiers, opts...)

	// Check to see if the dice is to be rolled as a debuff
	if str[0] == '-' {
		modifiers = append(modifiers, AsDebuff())
		str = strings.TrimSpace(str[1:])
	} else if str[0] == '+' {
		str = strings.TrimSpace(str[1:])
	}

	// If there isn't a multi-sided dice included, then it only contains a constant value
	if !strings.Contains(str, "d") {
		c, _ := strconv.Atoi(str)
		return NewConstant(c, modifiers...)
	}

	// Look for any constant damage that is added to a roll (+ or - a number at the end of the string)
	split := strings.Split(str, "+")
	if len(split) > 1 {
		constantValue, _ := strconv.Atoi(split[1])
		modifiers = append(modifiers, WithModifier(constantValue))
	} else {
		split = strings.Split(str, "-")
		if len(split) > 1 {
			constantValue, _ := strconv.Atoi(split[1])
			constantValue *= -1
			modifiers = append(modifiers, WithModifier(constantValue))
		}
	}

	// Get the number of dice and the number of sides per dice, defaulting to one
	// dice if the number of dice isn't specified
	var numDice int
	var numSides int
	split = strings.Split(split[0], "d")
	if len(split) > 1 {
		if split[0] != "" {
			numDice, _ = strconv.Atoi(split[0])
		} else {
			numDice = 1
		}
		numSides, _ = strconv.Atoi(split[1])
	}

	return NewDice(numDice, numSides, modifiers...)
}

// Customize creates a new dice from an existing one, applying the provided options. The dice passed in
// is not modified, but a new dice is returned with the options applied. This allows for
// creating a new dice based on an existing one, but with different options applied.
func (d *dice) Customize(opts ...DiceOption) Dice {
	newDice := &dice{
		numDice:  d.numDice,
		numSides: d.numSides,
		modifier: d.modifier,
		isLucky:  d.isLucky,
		isDebuff: d.isDebuff,
	}

	for _, opt := range opts {
		opt(newDice)
	}

	return newDice
}

// IsConstant returns true if the dice is a constant value.
func (d *dice) IsConstant() bool {
	return d.numDice == 0 && d.numSides == 0
}

// IsDebuff returns true if the dice is a debuff.
func (d *dice) IsDebuff() bool {
	return d.isDebuff
}

// IsLucky returns true if the dice is a lucky dice.
func (d *dice) IsLucky() bool {
	return d.isLucky
}

// NumSides returns the number of sides on the dice.
func (d *dice) NumSides() int {
	return d.numSides
}

// NumDice returns the number of dice to roll.
func (d *dice) NumDice() int {
	return d.numDice
}

// Modifier returns the constant value to add to the roll.
func (d *dice) Modifier() int {
	return d.modifier
}

// GetDice returns the set of dice being used. For a single dice, this will return a slice
// with a single element.
func (d *dice) GetDice() []Dice {
	return []Dice{d}
}

// Roll returns the results of rolling a number of multi-sided dice a single time, with or
// without a base damage.
func (d *dice) Roll(opts ...RollOption) Roll {
	r := &roll{
		dice:               d,
		rollType:           RollOnce,
		rolls:              make([]*singleRoll, 0, 2),
		criticalHit:        CriticalHit,
		criticalMiss:       CriticalMiss,
		criticalHitAllowed: false,
	}
	d.roll = r

	for _, opt := range opts {
		opt(r)
	}

	switch r.rollType {
	case RollWithAdvantage:
		r.rolls = []*singleRoll{d.rollDice(r.criticalHitAllowed, r.criticalHit, r.criticalMiss), d.rollDice(r.criticalHitAllowed, r.criticalHit, r.criticalMiss)}
		r.value = max(r.rolls[0].Value(), r.rolls[1].Value())
	case RollWithDisadvantage:
		r.rolls = []*singleRoll{d.rollDice(r.criticalHitAllowed, r.criticalHit, r.criticalMiss), d.rollDice(r.criticalHitAllowed, r.criticalHit, r.criticalMiss)}
		r.value = min(r.rolls[0].Value(), r.rolls[1].Value())
	default:
		r.rolls = []*singleRoll{d.rollDice(r.criticalHitAllowed, r.criticalHit, r.criticalMiss)}
		r.value = r.rolls[0].Value()
	}

	return r
}

// ReRoll re-rolls the dice with the provided options, returning the new Roll.
func (r *roll) ReRoll(opts ...RollOption) Roll {
	return r.dice.Roll(opts...)
}

// AsDebuff sets the dice as a debuff, which negates the value of the roll.
func AsDebuff() DiceOption {
	return func(d *dice) {
		d.isDebuff = true
	}
}

// WithModifier sets a constant value to be added to the roll of the dice.
func WithModifier(value int) DiceOption {
	return func(d *dice) {
		d.modifier = value
	}
}

// WithSource sets the source of the dice, which is used in the string representation of the dice.
func WithSource(source string) DiceOption {
	return func(d *dice) {
		d.source = source
	}
}

// GetRoll returns the last roll of the dice
func (d *dice) GetRoll() Roll {
	return d.roll
}

// Source returns the source of the roll.
func (d *dice) Source() string {
	return d.source
}

// rollDice rolls the dice and returns the value. If the dice is lucky, it will re-roll if it rolls a 1.
func (d *dice) rollDice(criticalHitAllowed bool, criticalHit int, criticalMiss int) *singleRoll {
	value := d.modifier
	for range d.numDice {
		rollValue := rng.Intn(d.numSides) + 1 // rng.Intn returns a value in the range [0, n), so we add 1 to get [1, n]
		if rollValue == 1 && d.isLucky {
			// If the dice is lucky, re-roll if it rolls a 1
			rollValue = rng.Intn(d.numSides) + 1
		}
		value += rollValue
	}

	// If this is a debuff dice, negate the value
	if d.isDebuff {
		value = -value
	}

	roll := &singleRoll{
		value:              value,
		dice:               d,
		criticalHitAllowed: criticalHitAllowed,
		criticalHit:        criticalHit,
		criticalMiss:       criticalMiss,
	}

	return roll
}

// isD20 checks if the dice is a D20. This is used to determine if the roll was a critical hit or miss.
func (d *dice) isD20() bool {
	return d.numSides == 20 && d.numDice == 1 && d.modifier == 0
}

// WithAdvantage rolls the dice with advantage. The dice will be rolled twice, and the highest value will be used.
func WithAdvantage() RollOption {
	return func(r *roll) {
		if r.rollType == RollWithDisadvantage {
			r.rollType = RollOnce
		} else {
			r.rollType = RollWithAdvantage
		}
	}
}

// WithDisadvantage rolls the dice with disadvantage. The dice will be rolled twice, and the lowest value will be used.
func WithDisadvantage() RollOption {
	return func(r *roll) {
		if r.rollType == RollWithAdvantage {
			r.rollType = RollOnce
		} else {
			r.rollType = RollWithDisadvantage
		}
	}
}

// WithLuck sets the dice to be lucky. If the dice rolls a 1, it will be re-rolled one more time and the new value will be used.
func WithLuck() DiceOption {
	return func(d *dice) {
		d.isLucky = true
	}
}

// WithCriticalHit sets the value for a critical hit. If the roll is greater than or equal to this value,
// then it is considered a critical hit.
// Defaults to 20.
func WithCriticalHit(value int) RollOption {
	return func(r *roll) {
		r.criticalHitAllowed = true
		r.criticalHit = value
	}
}

// WithCriticalMiss sets the value for a critical miss. If the roll is less than or equal to this value,
// then it is considered a critical miss.
// Defaults to 1.
func WithCriticalMiss(value int) RollOption {
	return func(r *roll) {
		r.criticalHitAllowed = true
		r.criticalMiss = value
	}
}

// WithCriticalHitAllowed sets the roll to allow for critical hits and misses. If set, a roll on a D20
// that is greater than or equal to the critical hit value will be considered a critical hit, and the
// roll that is less than or equal to the critical miss value will be considered a critical miss.
func WithCriticalHitAllowed() RollOption {
	return func(r *roll) {
		r.criticalHitAllowed = true
	}
}

// GetAllRolls returns all rolls for the dice. If a dice is rolled with advantage or disadvantage
// then two rolls will be returned. Otherwise, a single roll is included.
func (r *roll) GetAllRolls() []Roll {
	rolls := make([]Roll, 0, len(r.rolls))
	for _, roll := range r.rolls {
		rolls = append(rolls, roll)
	}
	return rolls
}

// Value gets the value of the roll, including the modifiers.
func (r *roll) Value() int {
	if !r.dice.isDebuff {
		return r.value
	}
	return -1 * r.value
}

// IsCriticalHit returns `true` if the roll was a critical hit; `false` otherwise
func (r *roll) IsCriticalHit() bool {
	return r.criticalHitAllowed && r.value >= r.criticalHit && r.dice.isD20()
}

// IsCriticalMiss returns `true` if the roll was a critical miss; `false` otherwise
func (r *roll) IsCriticalMiss() bool {
	return r.criticalHitAllowed && r.value <= r.criticalMiss && r.dice.isD20()
}

// RolledWithAdvantage returns `true` if the roll was made with advantage; `false` otherwise
func (r *roll) RolledWithAdvantage() bool {
	return r.rollType == RollWithAdvantage
}

// RolledWithDisadvantage returns `true` if the roll was made with disadvantage; `false` otherwise
func (r *roll) RolledWithDisadvantage() bool {
	return r.rollType == RollWithDisadvantage
}

// GetType gets the type of roll (ROLL_ONCE, ROLL_WITH_ADVANTAGE, ROLL_WITH_DISADVANTATE)
func (r *roll) GetType() RollType {
	return r.rollType
}

// GetDice gets the dice that was used for this roll.
func (r *roll) GetDice() Dice {
	return r.dice
}

// Check checks if the roll meets or exceeds the value.
func (r *roll) Check(v Value) bool {
	if r.IsCriticalHit() {
		return true
	}
	if r.IsCriticalMiss() {
		return false
	}
	return r.Value() >= v.Value()
}

// GetAllRolls returns all rolls for the dice. If a dice is rolled with advantage or disadvantage
// then two rolls will be returned. Otherwise, a single roll is included.
func (r *singleRoll) GetAllRolls() []Roll {
	return []Roll{r}
}

// Value gets the value of the roll, including the modifiers.
func (r *singleRoll) Value() int {
	if !r.dice.isDebuff {
		return r.value
	}
	return -1 * r.value
}

// IsCriticalHit returns `true` if the roll was a critical hit; `false` otherwise
func (r *singleRoll) IsCriticalHit() bool {
	return r.criticalHitAllowed && r.value >= r.criticalHit && r.dice.isD20()
}

// IsCriticalMiss returns `true` if the roll was a critical miss; `false` otherwise
func (r *singleRoll) IsCriticalMiss() bool {
	return r.criticalHitAllowed && r.value <= r.criticalMiss && r.dice.isD20()
}

// RolledWithAdvantage returns `true` if the roll was made with advantage; `false` otherwise
func (r *singleRoll) RolledWithAdvantage() bool {
	return false
}

// RolledWithDisadvantage returns `true` if the roll was made with disadvantage; `false` otherwise
func (r *singleRoll) RolledWithDisadvantage() bool {
	return false
}

// GetType gets the type of roll (ROLL_ONCE, ROLL_WITH_ADVANTAGE, ROLL_WITH_DISADVANTATE)
func (r *singleRoll) GetType() RollType {
	return RollOnce
}

// GetDice gets the dice that was used for this roll.
func (r *singleRoll) GetDice() Dice {
	return r.dice
}

// Check checks if the roll meets or exceeds the value.
func (r *singleRoll) Check(v Value) bool {
	if r.IsCriticalHit() {
		return true
	}
	if r.IsCriticalMiss() {
		return false
	}
	return r.Value() >= v.Value()
}

// ReRoll returns a new roll of the dice.
func (r *singleRoll) ReRoll(opts ...RollOption) Roll {
	// Re-roll using a `roll` object. The first roll is always a single roll, so return that.
	oldRoll := &roll{
		rollType:           RollOnce,
		rolls:              make([]*singleRoll, 0, 1),
		criticalHit:        r.criticalHit,
		criticalMiss:       r.criticalMiss,
		criticalHitAllowed: r.criticalHitAllowed,
		dice:               r.dice,
	}
	newRoll := oldRoll.ReRoll(opts...)
	return newRoll.GetAllRolls()[0]
}

// getDiceString returns a string representation of the dice. This includes both
// the dice as well as the modifier.
func getDiceString(numDice, numSides, modifier int) string {
	var sb strings.Builder
	if numDice > 0 {
		sb.WriteString(strconv.Itoa(numDice))
		sb.WriteString("d")
		sb.WriteString(strconv.Itoa(numSides))
	}

	if modifier != 0 {
		if modifier < 0 {
			sb.WriteString("-")
			sb.WriteString(strconv.Itoa(-1 * modifier))
		} else {
			if numDice != 0 {
				sb.WriteString("+")
			}
			sb.WriteString(strconv.Itoa(modifier))
		}
	} else if numDice == 0 && numSides == 0 {
		// For constant dice with value 0, return "0"
		sb.WriteString("0")
	}

	return sb.String()
}

// String returns a string representation of the dice. This includes both
// the dice as well as the source of the dice (if provided).
func (d *dice) String() string {
	var sb strings.Builder
	if d.isDebuff {
		sb.WriteString("-")
	}
	sb.WriteString(d.Str())

	return sb.String()
}

// Str returns a string representation of the dice. This includes both
// the dice and the source of the dice (if provided).
func (d *dice) Str() string {
	var sb strings.Builder
	sb.WriteString(getDiceString(d.numDice, d.numSides, d.modifier))

	if d.source != "" {
		sb.WriteString(" (")
		sb.WriteString(d.source)
		sb.WriteString(")")
	}

	return sb.String()
}

// String returns a string representation of the roll. This includes both
// the roll, the source of the dice (if provided), and the value of the roll.
func (r *roll) String() string {
	var sb strings.Builder

	if r.Value() < 0 {
		sb.WriteString("-")
	}
	sb.WriteString(r.Str())
	sb.WriteString(" = ")
	sb.WriteString(strconv.Itoa(r.Value()))

	return sb.String()
}

// Str returns a Str representation of the roll. This includes both
// the roll as well as the source of the dice (if provided), but not the
// value of the roll.
func (r *roll) Str() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(mathx.Abs(r.Value())))
	switch {
	case r.IsCriticalHit():
		sb.WriteString(" (Critical!)")
	case r.IsCriticalMiss():
		sb.WriteString(" (Miss!)")
	case !r.dice.IsConstant():
		sb.WriteString(" (")
		sb.WriteString(getDiceString(r.dice.numDice, r.dice.numSides, r.dice.modifier))
		if r.dice.Source() != "" {
			sb.WriteString(", ")
			sb.WriteString(r.dice.Source())
		}
		switch {
		case r.RolledWithAdvantage():
			sb.WriteString(", Advantage")
		case r.RolledWithDisadvantage():
			sb.WriteString(", Disadvantage")
		}
		sb.WriteString(")")
	case r.dice.Source() != "":
		sb.WriteString(" (")
		sb.WriteString(r.dice.Source())
		sb.WriteString(")")
	}

	return sb.String()
}

// String returns a string representation of the roll. This includes both
// the roll, the source of the dice (if provided), and the value of the roll.
func (r *singleRoll) String() string {
	var sb strings.Builder

	if r.Value() < 0 {
		sb.WriteString("-")
	}
	sb.WriteString(r.Str())
	sb.WriteString(" = ")
	sb.WriteString(strconv.Itoa(r.Value()))

	return sb.String()
}

// Str returns a Str representation of the roll. This includes both
// the roll as well as the source of the dice (if provided), but not the
// value of the roll.
func (r *singleRoll) Str() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(mathx.Abs(r.Value())))
	switch {
	case r.IsCriticalHit():
		sb.WriteString(" (Critical!)")
	case r.IsCriticalMiss():
		sb.WriteString(" (Miss!)")
	case !r.dice.IsConstant():
		sb.WriteString(" (")
		sb.WriteString(getDiceString(r.dice.numDice, r.dice.numSides, r.dice.modifier))
		if r.dice.Source() != "" {
			sb.WriteString(", ")
			sb.WriteString(r.dice.Source())
		}
		switch {
		case r.RolledWithAdvantage():
			sb.WriteString(", Advantage")
		case r.RolledWithDisadvantage():
			sb.WriteString(", Disadvantage")
		}
		sb.WriteString(")")
	case r.dice.Source() != "":
		sb.WriteString(" (")
		sb.WriteString(r.dice.Source())
		sb.WriteString(")")
	}

	return sb.String()
}
