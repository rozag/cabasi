package dice

type Dice uint8

const (
	D4  Dice = 4
	D6  Dice = 6
	D8  Dice = 8
	D10 Dice = 10
	D12 Dice = 12
	D20 Dice = 20
)

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
