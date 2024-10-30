package dice

import "fmt"

type Dice uint8

const (
	D4  Dice = 4
	D6  Dice = 6
	D8  Dice = 8
	D10 Dice = 10
	D12 Dice = 12
	D20 Dice = 20
)

// String returns the string representation of the Dice.
func (d Dice) String() string {
	switch d {
	case D4:
		return "D4"
	case D6:
		return "D6"
	case D8:
		return "D8"
	case D10:
		return "D10"
	case D12:
		return "D12"
	case D20:
		return "D20"
	default:
		panic(fmt.Errorf("unknown Dice: %d", d))
	}
}

type RNG interface {
	// UintN returns, as a uint, a non-negative pseudo-random number in the
	// half-open interval [0,n). It panics if n == 0.
	UintN(n uint) uint
}

func (d Dice) Roll(rng RNG) uint8 {
	// Suppressing "G115: integer overflow conversion uint -> uint8" because we're
	// getting a number from a dice roll, which is always in the [0,255] interval.
	return uint8(rng.UintN(uint(d)) + 1) // nolint:gosec
}
