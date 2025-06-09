# dice

A flexible Go package for dice rolling and dice-based mechanics, designed for tabletop role-playing games.

## Features

- Roll various types of dice (d4, d6, d8, d10, d12, d20, d100)
- Create custom dice with any number of sides
- Roll with advantage or disadvantage
- Apply modifiers to dice rolls
- Create dice sets for rolling multiple dice together
- Parse standard dice notation (e.g., "2d6+3")
- Support for critical hits and misses
- Lucky dice that re-roll on a 1
- Difficulty class checks
- Debuff dice (negative values)
- Detailed string representation of dice and rolls

## Installation

```bash
go get github.com/rbrabson/dice
```

## Usage

### Basic Dice Rolling

```go
package main

import (
    "fmt"

    "github.com/rbrabson/dice"
)

func main() {
    // Use predefined dice
    d20 := dice.D20
    roll := d20.Roll()
    fmt.Println("D20 roll:", roll)

    // Create custom dice
    attackDice := dice.NewDice(1, 20, dice.WithModifier(5), dice.WithSource("Longsword"))
    attackRoll := attackDice.Roll()
    fmt.Println("Attack roll:", attackRoll)

    // Parse dice notation
    damageDice := dice.ParseDice("2d6+3")
    damageRoll := damageDice.Roll()
    fmt.Println("Damage roll:", damageRoll)
}
```

### Rolling with Advantage/Disadvantage

```go
// Roll with advantage (roll twice, take the higher value)
advantageRoll := dice.D20.Roll(dice.WithAdvantage())
fmt.Println("Roll with advantage:", advantageRoll)

// Roll with disadvantage (roll twice, take the lower value)
disadvantageRoll := dice.D20.Roll(dice.WithDisadvantage())
fmt.Println("Roll with disadvantage:", disadvantageRoll)
```

### Dice Sets

```go
// Create a dice set for complex damage calculations
damageSet := dice.NewDiceSet(
    dice.ParseDice("1d8", dice.WithSource("Longsword")),
    dice.ParseDice("1d6", dice.WithSource("Fire Damage")),
    dice.NewConstant(3, dice.WithSource("Strength Bonus")),
)
damageRoll := damageSet.Roll()
fmt.Println("Complex damage roll:", damageRoll)
```

### Difficulty Class Checks

```go
// Create a difficulty class
dc := dice.NewDifficultyClass(15)

// Roll against the difficulty class
attackRoll := dice.D20.Roll(dice.WithCriticalHitAllowed())
success := attackRoll.Check(dc)

if success {
    fmt.Println("Attack succeeded!")
} else {
    fmt.Println("Attack failed!")
}
```

### Critical Hits and Misses

```go
// Roll with critical hit/miss detection
attackRoll := dice.D20.Roll(dice.WithCriticalHitAllowed())

if attackRoll.IsCriticalHit() {
    fmt.Println("Critical hit!")
} else if attackRoll.IsCriticalMiss() {
    fmt.Println("Critical miss!")
} else {
    fmt.Println("Normal roll:", attackRoll.Value())
}
```

### Lucky Dice

```go
// Create a lucky dice that re-rolls on a 1
luckyDice := dice.NewDice(1, 6, dice.WithLuck())
luckyRoll := luckyDice.Roll()
fmt.Println("Lucky roll:", luckyRoll)
```

## API Documentation

### Predefined Dice

- `D4`: A standard 4-sided die
- `D6`: A standard 6-sided die
- `D8`: A standard 8-sided die
- `D10`: A standard 10-sided die
- `D12`: A standard 12-sided die
- `D20`: A standard 20-sided die
- `D100`: A standard 100-sided die

### Creating Dice

- `NewDice(numDice, numSides int, opts ...DiceOption)`: Create a new dice with the specified number of dice and sides
- `NewConstant(value int, opts ...DiceOption)`: Create a dice that always returns the same value
- `ParseDice(str string, opts ...DiceOption)`: Parse a string representation of a dice (e.g., "2d6+3")
- `NewDiceSet(dice ...Dice)`: Create a set of dice that can be rolled together

### Dice Options

- `WithModifier(value int)`: Add a constant modifier to the dice roll
- `WithSource(source string)`: Set the source of the dice (for display purposes)
- `AsDebuff()`: Set the dice as a debuff (negates the value)
- `WithLuck()`: Make the dice lucky (re-rolls on a 1)

### Roll Options

- `WithAdvantage()`: Roll with advantage (roll twice, take the higher value)
- `WithDisadvantage()`: Roll with disadvantage (roll twice, take the lower value)
- `WithCriticalHitAllowed()`: Allow critical hits and misses
- `WithCriticalHit(value int)`: Set the value for a critical hit
- `WithCriticalMiss(value int)`: Set the value for a critical miss

### Difficulty Classes

- `NewDifficultyClass(targetValue int)`: Create a new difficulty class with the specified target value

## License

This project is licensed under the terms found in the LICENSE file.
